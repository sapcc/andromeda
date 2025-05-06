// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/apex/log"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
)

var (
	// resource types that do not map directly to a `quota` table column
	// due to their quota being bound to a provider
	providerBoundResourceQuotas = []string{"domain"}
)

type (
	errMissingProvider struct{ Resource string }
	errUnknownResource struct{ Resource string }

	providerBoundResourceRequest interface {
		GetResource() resource
	}

	newDomainRequest struct {
		Resource resource `json:"domain"`
	}

	resource struct {
		Provider string `json:"provider"`
	}

	quotaController struct {
		db *sqlx.DB
	}
)

func (r *newDomainRequest) GetResource() resource {
	return r.Resource
}

func NewQuotaController(db *sqlx.DB) *quotaController {
	return &quotaController{db: db}
}

func (e *errMissingProvider) Error() string {
	return fmt.Sprintf("request body for '%s' contains empty 'provider' field when non-empty 'provider' field required", e.Resource)
}

func (e *errUnknownResource) Error() string {
	return fmt.Sprintf("unknown resource '%s'", e.Resource)
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

		// Handle quota only for resource creation
		if action != "post" {
			next.ServeHTTP(w, r)
			return
		}

		var (
			provider string
			err      error
		)

		if slices.Contains(providerBoundResourceQuotas, resource) {
			log.Debugf("Resource quota for '%s' is bound to provider", resource)
			provider, err = providerFromHTTPRequest(resource, w, r)
			if err != nil {
				middleware.
					Error(401, utils.GetInvalidProviderBoundResourceResponse(resource), utils.JSONHeader).
					WriteResponse(w, runtime.JSONProducer())
				return
			}
			log.Debugf("Provider read from request body: %s", provider)
		}

		// Get project scope
		project, err := auth.ProjectScopeForRequest(r)
		if err != nil {
			middleware.
				Error(401, utils.Unauthorized(err), utils.JSONHeader).
				WriteResponse(w, runtime.JSONProducer())
			return
		}

		// Skip quota check disabled keystonemiddleware
		if project == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Check for quota
		var quotaAvailable, quotaUsed int

		insert := sq.Insert("quota").
			Columns("project_id", "domain_akamai", "domain_f5", "pool", "member", "monitor", "datacenter").
			Values(
				project,
				config.Global.Quota.DefaultQuotaDomainAkamai,
				config.Global.Quota.DefaultQuotaDomainF5,
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

		if provider != "" {
			subquery = subquery.Where(sq.Eq{"provider": provider})
		}

		query := sq.Select(quotaTableColumnForResource(resource, provider)).
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

		log.Debugf("Quota %s of project %s is %d of %d", quotaTableColumnForResource(resource, provider), project, quotaUsed, quotaAvailable)
		if quotaAvailable-quotaUsed < 1 {
			middleware.
				Error(403, utils.GetQuotaMetResponse(resource), utils.JSONHeader).
				WriteResponse(w, runtime.JSONProducer())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func providerFromHTTPRequest(resource string, w http.ResponseWriter, r *http.Request) (string, error) {
	resourceReq, err := newResourceRequest(resource)
	if err != nil {
		return "", err
	}
	bodyBytes, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20)) // 1MB
	if err != nil {
		return "", err
	}
	err = r.Body.Close()
	if err != nil {
		return "", err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&resourceReq); err != nil {
		return "", err
	}
	provider := resourceReq.GetResource().Provider
	if provider == "" {
		return "", &errMissingProvider{resource}
	}
	return provider, nil
}

func newResourceRequest(resource string) (providerBoundResourceRequest, error) {
	if resource == "domain" {
		return &newDomainRequest{}, nil
	}
	return nil, &errUnknownResource{resource}
}

func quotaTableColumnForResource(resource string, provider string) string {
	if slices.Contains(providerBoundResourceQuotas, resource) {
		return fmt.Sprintf("%s_%s", resource, provider)
	}
	return resource
}
