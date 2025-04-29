// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/datacenters"
)

type DatacenterController struct {
	CommonController
}

// GetDatacenters GET /datacenters
func (c DatacenterController) GetDatacenters(params datacenters.GetDatacentersParams) middleware.Responder {
	filter := make(map[string]interface{}, 0)
	pagination := db.Pagination{
		HTTPRequest: params.HTTPRequest,
		Limit:       params.Limit,
		Marker:      params.Marker,
		PageReverse: params.PageReverse,
		Sort:        params.Sort,
	}
	if params.DatacenterID != nil {
		filter["id"] = *params.DatacenterID
	}
	rows, err := pagination.Query(c.db, "SELECT * FROM datacenter", filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return datacenters.NewGetDatacentersBadRequest().WithPayload(utils.InvalidMarker)
		}
		panic(err)
	}

	//goland:noinspection GoPreferNilSlice
	var _datacenters = []*models.Datacenter{}
	for rows.Next() {
		datacenter := models.Datacenter{}
		if err := rows.StructScan(&datacenter); err != nil {
			panic(err)
		}

		// Filter result based on policy
		requestVars := map[string]string{"project_id": *datacenter.ProjectID, "scope": *datacenter.Scope}
		if _, err = auth.Authenticate(params.HTTPRequest, requestVars); err == nil {
			_datacenters = append(_datacenters, &datacenter)
		}
	}

	_links := pagination.GetLinks(_datacenters)
	payload := datacenters.GetDatacentersOKBody{Datacenters: _datacenters, Links: _links}
	return datacenters.NewGetDatacentersOK().WithPayload(&payload)
}

// PostDatacenters POST /datacenters
func (c DatacenterController) PostDatacenters(params datacenters.PostDatacentersParams) middleware.Responder {
	datacenter := params.Datacenter.Datacenter

	if projectId, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return datacenters.NewPostDatacentersDefault(403).WithPayload(utils.PolicyForbidden)
	} else {
		datacenter.ProjectID = &projectId
	}

	// Set default values
	if err := utils.SetModelDefaults(datacenter); err != nil {
		panic(err)
	}

	sql := `
		INSERT INTO datacenter 
    		(name, admin_state_up, continent, country, state_or_province, 
    		 city, latitude, longitude, scope, project_id, provider)
		VALUES 
		    (:name, :admin_state_up, :continent, :country, :state_or_province, 
		     :city, :latitude, :longitude, :scope, :project_id, :provider)
		RETURNING *
	`
	stmt, _ := c.db.PrepareNamed(sql)
	if err := stmt.Get(datacenter, datacenter); err != nil {
		panic(err)
	}

	_ = PendingSync(c.rpc)
	return datacenters.NewPostDatacentersCreated().WithPayload(&datacenters.PostDatacentersCreatedBody{Datacenter: datacenter})
}

// GetDatacentersDatacenterID GET /datacenters/:id
func (c DatacenterController) GetDatacentersDatacenterID(params datacenters.GetDatacentersDatacenterIDParams) middleware.Responder {
	datacenter := models.Datacenter{ID: params.DatacenterID}
	err := PopulateDatacenter(c.db, &datacenter, []string{"*"})
	if err != nil {
		return datacenters.NewGetDatacentersDatacenterIDNotFound().WithPayload(utils.NotFound)
	}

	requestVars := map[string]string{"project_id": *datacenter.ProjectID, "scope": *datacenter.Scope}
	if _, err = auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return datacenters.NewGetDatacentersDatacenterIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	return datacenters.NewGetDatacentersDatacenterIDOK().WithPayload(&datacenters.GetDatacentersDatacenterIDOKBody{Datacenter: &datacenter})
}

// PutDatacentersDatacenterID PUT /datacenters/:id
func (c DatacenterController) PutDatacentersDatacenterID(params datacenters.PutDatacentersDatacenterIDParams) middleware.Responder {
	datacenter := models.Datacenter{ID: params.DatacenterID}
	err := PopulateDatacenter(c.db, &datacenter, []string{"project_id"})
	if err != nil {
		return datacenters.NewPutDatacentersDatacenterIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *datacenter.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return datacenters.NewPutDatacentersDatacenterIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	params.Datacenter.Datacenter.ID = params.DatacenterID
	sql := `
		UPDATE datacenter SET
			name = COALESCE(:name, name),
			admin_state_up = COALESCE(:admin_state_up, admin_state_up),
			continent = COALESCE(:continent, continent), 
			country = COALESCE(:country, country), 
			state_or_province = COALESCE(:state_or_province, state_or_province), 
			city = COALESCE(:city, city), 
			latitude = COALESCE(:latitude, latitude), 
			longitude = COALESCE(:longitude, longitude), 
			scope = COALESCE(:scope, scope),
		    updated_at = NOW()
		WHERE id = :id
	`
	if _, err := c.db.NamedExec(sql, params.Datacenter.Datacenter); err != nil {
		panic(err)
	}

	// Update datacenter response
	if err := PopulateDatacenter(c.db, &datacenter, []string{"*"}); err != nil {
		panic(err)
	}

	_ = PendingSync(c.rpc)
	return datacenters.NewPutDatacentersDatacenterIDAccepted().WithPayload(
		&datacenters.PutDatacentersDatacenterIDAcceptedBody{Datacenter: &datacenter})
}

// DeleteDatacentersDatacenterID DELETE /datacenters/:id
func (c DatacenterController) DeleteDatacentersDatacenterID(params datacenters.DeleteDatacentersDatacenterIDParams) middleware.Responder {
	datacenter := models.Datacenter{ID: params.DatacenterID}
	if err := PopulateDatacenter(c.db, &datacenter, []string{"project_id"}); err != nil {
		return datacenters.NewDeleteDatacentersDatacenterIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *datacenter.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return datacenters.NewDeleteDatacentersDatacenterIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	var count int
	sql := c.db.Rebind(`SELECT count(*) FROM member WHERE datacenter_id = ?`)
	if err := c.db.Get(&count, sql, datacenter.ID); err != nil {
		panic(err)
	} else if count > 0 {
		return datacenters.NewDeleteDatacentersDatacenterIDDefault(409).WithPayload(utils.DatacenterInUse)
	}

	sql = c.db.Rebind(`UPDATE datacenter SET provisioning_status = 'PENDING_DELETE' WHERE id = ?`)
	res, err := c.db.Exec(sql, params.DatacenterID)
	if err != nil {
		var pe *pq.Error
		if errors.As(err, &pe) && pgerrcode.IsIntegrityConstraintViolation(string(pe.Code)) {
			return datacenters.NewDeleteDatacentersDatacenterIDDefault(409).WithPayload(utils.DatacenterInUse)
		}
		if utils.MySQLForeignKeyViolation.Is(err) {
			return datacenters.NewDeleteDatacentersDatacenterIDDefault(409).WithPayload(utils.DatacenterInUse)
		}
		panic(err)
	}
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		return datacenters.NewDeleteDatacentersDatacenterIDNotFound().WithPayload(utils.NotFound)
	}

	_ = PendingSync(c.rpc)
	return datacenters.NewDeleteDatacentersDatacenterIDNoContent()
}

// PopulateDatacenter populates attributes of a datacenter instance based on it's ID
func PopulateDatacenter(db *sqlx.DB, datacenter *models.Datacenter, fields []string) error {
	sql := db.Rebind(fmt.Sprintf(`SELECT %s FROM datacenter WHERE id = ?`, strings.Join(fields, ", ")))
	if err := db.Get(datacenter, sql, datacenter.ID); err != nil {
		return err
	}
	return nil
}
