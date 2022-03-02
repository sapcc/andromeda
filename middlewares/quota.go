/*
 *   Copyright 2021 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asim/go-micro/v3/logger"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
)

type quotaController struct {
	db *sqlx.DB
}

func NewQuotaController(db *sqlx.DB) *quotaController {
	return &quotaController{db: db}
}

//quotaResponseWriter is a wrapper of regular ResponseWriter
type quotaResponseWriter struct {
	http.ResponseWriter
	resource  string
	action    string
	projectID string
	db        *sqlx.DB
}

func (qc *quotaController) NewQuotaResponseWriter(w http.ResponseWriter, r string, a string, p string) *quotaResponseWriter {
	return &quotaResponseWriter{w, r, a, p, qc.db}
}

//WriteHeader
func (qrw *quotaResponseWriter) WriteHeader(code int) {
	qrw.ResponseWriter.WriteHeader(code)

	var operator rune
	if qrw.action == "post" {
		operator = '+'
	} else {
		operator = '-'
	}

	sql := fmt.Sprintf(
		`INSERT INTO quota 
                SET in_use_%s = 1,
					domain = ?,
					pool = ?,
					member = ?,
					monitor = ?,
                    project_id = ?
                ON DUPLICATE KEY UPDATE 
					in_use_%s = in_use_%s %c 1`,
		qrw.resource, qrw.resource, qrw.resource, operator)
	if 200 < code && 205 > code {
		if _, err := qrw.db.Exec(sql,
			config.Global.Quota.DefaultQuotaDomain,
			config.Global.Quota.DefaultQuotaPool,
			config.Global.Quota.DefaultQuotaMember,
			config.Global.Quota.DefaultQuotaMonitor,
			qrw.projectID); err != nil {
			logger.Error(err)
		}
	}
}

//QuotaHandler provides the quota enforcement.
func (qc *quotaController) QuotaHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := strings.Split(policy.RuleFromHTTPRequest(r), ":")

		if len(target) != 3 {
			next.ServeHTTP(w, r)
		}

		// resource and action are sql-safe resource names from swagger spec
		resource := target[1]
		action := target[2]

		// Handle quota only for resource creation/deletion
		if action != "delete" && action != "post" {
			next.ServeHTTP(w, r)
		}

		// Get project scope
		project, err := auth.ProjectScopeForRequest(r)
		if err != nil {
			next.ServeHTTP(w, r)
		}

		// Check Quota increase possible before processing request
		if action == "post" {
			var quotaAvailable int
			sql := fmt.Sprintf(`SELECT %s - in_use_%s FROM quota WHERE project_id = ?`, resource, resource)
			if err := qc.db.Get(&quotaAvailable, sql, project); err != nil {
				logger.Debug(err)
			} else {
				logger.Debug("Quota is ", quotaAvailable)
				if quotaAvailable < 1 {
					err := fmt.Sprintf("Quota has been met for resources: %s", resource)
					http.Error(w, err, http.StatusForbidden)
					return
				}
			}
		}

		qrw := qc.NewQuotaResponseWriter(w, resource, action, project)
		next.ServeHTTP(qrw, r)
	})
}
