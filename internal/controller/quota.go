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

package controller

import (
	dbsql "database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
	"github.com/sapcc/andromeda/restapi/operations/domains"
)

type QuotaController struct {
	db *sqlx.DB
}

func getQuotas(db *sqlx.DB, projectID *string) ([]*administrative.GetQuotasOKBodyQuotasItems0, error) {
	var rows *sqlx.Rows
	var err error
	if projectID != nil {
		rows, err = db.Queryx(`SELECT * FROM quota WHERE project_id = ?`, projectID)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = db.Queryx(`SELECT * FROM quota`)
		if err != nil {
			return nil, err
		}
	}

	//goland:noinspection GoPreferNilSlice
	var quotas = []*administrative.GetQuotasOKBodyQuotasItems0{}
	for rows.Next() {
		var p administrative.GetQuotasOKBodyQuotasItems0
		if err = rows.StructScan(&p); err != nil {
			return nil, err
		}
		quotas = append(quotas, &p)
	}
	return quotas, nil
}

// GetQuotas GET /quotas
func (c QuotaController) GetQuotas(params administrative.GetQuotasParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest); err != nil {
		return administrative.NewGetQuotasDefault(403).WithPayload(utils.PolicyForbidden)
	}

	responseQuotas, err := getQuotas(c.db, params.ProjectID)
	if err != nil {
		return domains.NewGetDomainsDefault(404).WithPayload(utils.NotFound)
	}
	return administrative.NewGetQuotasOK().WithPayload(&administrative.GetQuotasOKBody{Quotas: responseQuotas})
}

func (c QuotaController) GetQuotasProjectID(params administrative.GetQuotasProjectIDParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest); err != nil {
		return administrative.NewGetQuotasProjectIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	body := administrative.GetQuotasProjectIDOKBody{}

	sql, args, err := sq.Select("domain, pool, member, monitor, datacenter").
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("domain").Where(sq.Eq{"project_id": params.ProjectID}).
				Where(sq.NotEq{"provisioning_status": "DELETED"}),
			"in_use_domain")).
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("pool").Where(sq.Eq{"project_id": params.ProjectID}),
			"in_use_pool")).
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("member").Where(sq.Eq{"project_id": params.ProjectID}),
			"in_use_member")).
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("monitor").Where(sq.Eq{"project_id": params.ProjectID}),
			"in_use_monitor")).
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("datacenter").Where(sq.Eq{"project_id": params.ProjectID}),
			"in_use_datacenter")).
		From("quota").
		Where(sq.Eq{"project_id": params.ProjectID}).
		ToSql()

	if err != nil {
		panic(err)
	}

	if err := c.db.Get(&body.Quota, c.db.Rebind(sql), args...); err != nil {
		if errors.Is(err, dbsql.ErrNoRows) {
			return administrative.NewGetQuotasProjectIDNotFound().WithPayload(utils.NotFound)
		}
		panic(err)
	}
	return administrative.NewGetQuotasProjectIDOK().WithPayload(&body)

}

func (c QuotaController) GetQuotasDefaults(params administrative.GetQuotasDefaultsParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest); err != nil {
		return administrative.NewGetQuotasDefaultsDefault(403).WithPayload(utils.PolicyForbidden)
	}

	body := administrative.GetQuotasDefaultsOKBody{
		Quota: &models.Quota{
			Datacenter: &config.Global.Quota.DefaultQuotaDatacenter,
			Domain:     &config.Global.Quota.DefaultQuotaDomain,
			Member:     &config.Global.Quota.DefaultQuotaMember,
			Monitor:    &config.Global.Quota.DefaultQuotaMonitor,
			Pool:       &config.Global.Quota.DefaultQuotaPool,
		},
	}
	return administrative.NewGetQuotasDefaultsOK().WithPayload(&body)
}

func (c QuotaController) PutQuotasProjectID(params administrative.PutQuotasProjectIDParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest); err != nil {
		return administrative.NewPutQuotasProjectIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	// Set defaults
	if params.Quota.Quota.Datacenter == nil {
		params.Quota.Quota.Datacenter = &config.Global.Quota.DefaultQuotaDatacenter
	}
	if params.Quota.Quota.Domain == nil {
		params.Quota.Quota.Domain = &config.Global.Quota.DefaultQuotaDomain
	}
	if params.Quota.Quota.Member == nil {
		params.Quota.Quota.Member = &config.Global.Quota.DefaultQuotaMember
	}
	if params.Quota.Quota.Monitor == nil {
		params.Quota.Quota.Monitor = &config.Global.Quota.DefaultQuotaMonitor
	}
	if params.Quota.Quota.Pool == nil {
		params.Quota.Quota.Pool = &config.Global.Quota.DefaultQuotaPool
	}

	quota := struct {
		*models.Quota
		ProjectID *string
	}{params.Quota.Quota, &params.ProjectID}

	var sql string
	if c.db.DriverName() == "mysql" {
		sql = `
			INSERT INTO quota SET
				domain = :domain,
				pool = :pool,
				member = :member,
				monitor = :monitor,
			    datacenter = :datacenter,
				project_id = :project_id
			ON DUPLICATE KEY UPDATE
				domain = COALESCE(:domain, domain),
				pool = COALESCE(:pool, pool),
				member = COALESCE(:member, member), 
				monitor = COALESCE(:monitor, monitor),
				datacenter = COALESCE(:datacenter, datacenter)
		`
	} else {
		sql = `
			INSERT INTO quota
				(domain, pool, member, monitor, datacenter, project_id)
			VALUES 
			    (:domain, :pool, :member, :monitor, :datacenter, :project_id)
			ON CONFLICT (project_id) DO UPDATE SET 
				domain = COALESCE(:domain, quota.domain),
				pool = COALESCE(:pool, quota.pool),
				member = COALESCE(:member, quota.member), 
				monitor = COALESCE(:monitor, quota.monitor),
				datacenter = COALESCE(:datacenter, quota.datacenter)
		`
	}
	if _, err := c.db.NamedExec(sql, quota); err != nil {
		panic(err)
	}

	return administrative.NewPutQuotasProjectIDAccepted().WithPayload(
		&administrative.PutQuotasProjectIDAcceptedBody{Quota: quota.Quota})
}

func (c QuotaController) DeleteQuotasProjectID(params administrative.DeleteQuotasProjectIDParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest); err != nil {
		return administrative.NewDeleteQuotasProjectIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	sql := c.db.Rebind(`DELETE FROM quota WHERE project_id = ?`)
	res := c.db.MustExec(sql, params.ProjectID)
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		return administrative.NewDeleteQuotasProjectIDNotFound().WithPayload(utils.NotFound)
	}
	return administrative.NewDeleteQuotasProjectIDNoContent()
}
