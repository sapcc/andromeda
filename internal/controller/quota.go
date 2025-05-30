// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	dbsql "database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-openapi/runtime/middleware"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type QuotaController struct {
	CommonController
}

// GetQuotas GET /quotas
func (c QuotaController) GetQuotas(params administrative.GetQuotasParams) middleware.Responder {
	rows, err := c.db.Queryx(`SELECT * FROM quota`)
	if err != nil {
		panic(err)
	}

	//zero-length slice used because we want [] via json encoder, nil encodes null
	//goland:noinspection GoPreferNilSlice
	var _quotas = []*administrative.GetQuotasOKBodyQuotasItems0{}
	for rows.Next() {
		var p administrative.GetQuotasOKBodyQuotasItems0
		if err = rows.StructScan(&p); err != nil {
			panic(err)
		}
		if _, err = auth.Authenticate(params.HTTPRequest, map[string]string{"project_id": *p.ProjectID}); err == nil {
			_quotas = append(_quotas, &p)
		}
	}

	return administrative.NewGetQuotasOK().WithPayload(&administrative.GetQuotasOKBody{Quotas: _quotas})
}

func (c QuotaController) GetQuotasProjectID(params administrative.GetQuotasProjectIDParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, map[string]string{"project_id": params.ProjectID}); err != nil {
		return administrative.NewGetQuotasProjectIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	body := administrative.GetQuotasProjectIDOKBody{}

	sql, args, err := sq.Select("domain_akamai, domain_f5, pool, member, monitor, datacenter").
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("domain").Where(sq.Eq{"project_id": params.ProjectID, "provider": "akamai"}).
				Where(sq.NotEq{"provisioning_status": "DELETED"}),
			"in_use_domain_akamai")).
		Column(sq.Alias(
			sq.Select("COUNT(id)").From("domain").Where(sq.Eq{"project_id": params.ProjectID, "provider": "f5"}).
				Where(sq.NotEq{"provisioning_status": "DELETED"}),
			"in_use_domain_f5")).
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
	if _, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return administrative.NewGetQuotasDefaultsDefault(403).WithPayload(utils.PolicyForbidden)
	}

	body := administrative.GetQuotasDefaultsOKBody{
		Quota: &models.Quota{
			Datacenter:   &config.Global.Quota.DefaultQuotaDatacenter,
			DomainAkamai: &config.Global.Quota.DefaultQuotaDomainAkamai,
			DomainF5:     &config.Global.Quota.DefaultQuotaDomainF5,
			Member:       &config.Global.Quota.DefaultQuotaMember,
			Monitor:      &config.Global.Quota.DefaultQuotaMonitor,
			Pool:         &config.Global.Quota.DefaultQuotaPool,
		},
	}
	return administrative.NewGetQuotasDefaultsOK().WithPayload(&body)
}

func (c QuotaController) PutQuotasProjectID(params administrative.PutQuotasProjectIDParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, map[string]string{"project_id": params.ProjectID}); err != nil {
		return administrative.NewPutQuotasProjectIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	quota := struct {
		*models.Quota
		ProjectID *string
	}{
		params.Quota.Quota,
		&params.ProjectID,
	}

	var sql string
	if c.db.DriverName() == "mysql" {
		sql = `
			INSERT INTO quota SET
				domain_akamai = COALESCE(:domain_akamai, %d),
				domain_f5 = COALESCE(:domain_f5, %d),
				pool = COALESCE(:pool, %d),
				member = COALESCE(:member, %d),
				monitor = COALESCE(:monitor, %d),
			    datacenter = COALESCE(:datacenter, %d),
				project_id = :project_id
			ON DUPLICATE KEY UPDATE
				domain_akamai = COALESCE(:domain_akamai, domain_akamai),
				domain_f5 = COALESCE(:domain_f5, domain_f5),
				pool = COALESCE(:pool, pool),
				member = COALESCE(:member, member), 
				monitor = COALESCE(:monitor, monitor),
				datacenter = COALESCE(:datacenter, datacenter)
			RETURNING *
		`
	} else {
		sql = `
			INSERT INTO quota
				(domain_akamai, domain_f5, pool, member, monitor, datacenter, project_id)
			VALUES 
			    (
					 COALESCE(:domain_akamai, %d),
					 COALESCE(:domain_f5, %d),
					 COALESCE(:pool, %d),
					 COALESCE(:member, %d),
					 COALESCE(:monitor, %d),
					 COALESCE(:datacenter, %d),
					 :project_id
			     )
			ON CONFLICT (project_id) DO UPDATE SET 
				domain_akamai = COALESCE(:domain_akamai, quota.domain_akamai),
				domain_f5 = COALESCE(:domain_f5, quota.domain_f5),
				pool = COALESCE(:pool, quota.pool),
				member = COALESCE(:member, quota.member), 
				monitor = COALESCE(:monitor, quota.monitor),
				datacenter = COALESCE(:datacenter, quota.datacenter)
			RETURNING *
		`
	}
	sql = fmt.Sprintf(sql,
		config.Global.Quota.DefaultQuotaDomainAkamai,
		config.Global.Quota.DefaultQuotaDomainF5,
		config.Global.Quota.DefaultQuotaPool,
		config.Global.Quota.DefaultQuotaMember,
		config.Global.Quota.DefaultQuotaMonitor,
		config.Global.Quota.DefaultQuotaDatacenter)
	nstmt, err := c.db.PrepareNamed(sql)
	if err != nil {
		panic(err)
	}
	if err = nstmt.Get(&quota, quota); err != nil {
		panic(err)
	}
	return administrative.NewPutQuotasProjectIDAccepted().WithPayload(
		&administrative.PutQuotasProjectIDAcceptedBody{Quota: quota.Quota})
}

func (c QuotaController) DeleteQuotasProjectID(params administrative.DeleteQuotasProjectIDParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, map[string]string{"project_id": params.ProjectID}); err != nil {
		return administrative.NewDeleteQuotasProjectIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	sql := c.db.Rebind(`DELETE FROM quota WHERE project_id = ?`)
	res := c.db.MustExec(sql, params.ProjectID)
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		return administrative.NewDeleteQuotasProjectIDNotFound().WithPayload(utils.NotFound)
	}
	return administrative.NewDeleteQuotasProjectIDNoContent()
}
