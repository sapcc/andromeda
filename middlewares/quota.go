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
	"net/http"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
)

type quotaController struct {
	db *sqlx.DB
}

func NewQuotaController(db *sqlx.DB) *quotaController {
	return &quotaController{db: db}
}

// QuotaHandler provides the quota enforcement.
func (qc *quotaController) QuotaHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := strings.Split(policy.RuleFromHTTPRequest(r), ":")

		if len(target) != 3 {
			next.ServeHTTP(w, r)
			return
		}

		// resource and action are sql-safe resource names from swagger spec
		resource := target[1]
		action := target[2]

		// Handle quota only for resource creation/deletion
		if action != "post" {
			next.ServeHTTP(w, r)
			return
		}

		// Get project scope
		project, err := auth.ProjectScopeForRequest(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Check for quota
		var quotaAvailable, quotaUsed int

		insert := sq.Insert("quota").
			Columns("project_id", "domain", "pool", "member", "monitor", "datacenter").
			Values(
				project,
				config.Global.Quota.DefaultQuotaDomain,
				config.Global.Quota.DefaultQuotaPool,
				config.Global.Quota.DefaultQuotaMember,
				config.Global.Quota.DefaultQuotaMonitor,
				config.Global.Quota.DefaultQuotaDatacenter,
			)
		if qc.db.DriverName() == "mysql" {
			insert = insert.Options("IGNORE").
				PlaceholderFormat(sq.Question)
		} else {
			insert = insert.Suffix("ON CONFLICT (project_id) DO NOTHING").
				PlaceholderFormat(sq.Dollar)
		}

		sql, args, err := insert.ToSql()
		if err != nil {
			panic(err)
		}
		if _, err := qc.db.Exec(sql, args...); err != nil {
			panic(err)
		}

		subquery := sq.
			Select("COUNT(id)").
			From(resource).
			Where(sq.Eq{"project_id": project}).
			Where(sq.NotEq{"provisioning_status": "DELETED"})

		query := sq.Select(resource).
			Column(sq.Alias(subquery, "quota_usage")).
			From("quota").
			Where(sq.Eq{"project_id": project})

		sql, args, err = query.ToSql()
		if err != nil {
			panic(err)
		}

		if err := qc.db.QueryRowx(qc.db.Rebind(sql), args...).Scan(&quotaAvailable, &quotaUsed); err != nil {
			panic(err)
		}

		logger.Debugf("Quota %s of project %s is %d of %d", resource, project, quotaUsed, quotaAvailable)
		if quotaAvailable-quotaUsed < 1 {
			ret, _ := utils.GetQuotaMetResponse(resource).MarshalBinary()
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(403)
			if _, err := w.Write(ret); err != nil {
				panic(err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
