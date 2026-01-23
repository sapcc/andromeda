// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	dbsql "database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/geographic_maps"
)

type GeoMapController struct {
	CommonController
}

// GetGeomaps GET /geoMaps
func (c GeoMapController) GetGeomaps(params geographic_maps.GetGeomapsParams) middleware.Responder {
	filter := make(map[string]any, 0)
	pagination := db.Pagination{
		HTTPRequest: params.HTTPRequest,
		Limit:       params.Limit,
		Marker:      params.Marker,
		PageReverse: params.PageReverse,
		Sort:        params.Sort,
	}
	sql := `SELECT * FROM geographic_map`
	if params.DatacenterID != nil {
		filter["datacenter"] = *params.DatacenterID
		sql = `SELECT 
    		geographic_map.id AS id, 
    		geographic_map.name AS name, 
    		geographic_map.provisioning_status AS provisioning_status,
    		geographic_map.scope AS scope,
    		geographic_map.project_id AS project_id,
    		geographic_map.provider AS provider,
    		geographic_map.created_at AS created_at,
    		geographic_map.updated_at AS updated_at
			FROM geographic_map
		    JOIN geographic_map_assignment ON 
		        geographic_map.id = geographic_map_assignment.geographic_map_id OR 
		        geographic_map.default_datacenter = geographic_map_assignment.datacenter`
	}
	if params.DefaultDatacenterID != nil {
		filter["default_datacenter"] = *params.DefaultDatacenterID
	}
	rows, err := pagination.Query(c.db, sql, filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return geographic_maps.NewGetGeomapsBadRequest().WithPayload(utils.InvalidMarker)
		}
		panic(err)
	}

	//goland:noinspection GoPreferNilSlice
	var _geoMaps = []*models.Geomap{}
	for rows.Next() {
		geoMap := models.Geomap{}
		if err := rows.StructScan(&geoMap); err != nil {
			panic(err)
		}

		// Filter result based on policy
		requestVars := map[string]string{"project_id": *geoMap.ProjectID, "scope": *geoMap.Scope}
		if _, err = auth.Authenticate(params.HTTPRequest, requestVars); err == nil {
			if err := PopulateGeoMapAssignments(c.db, &geoMap); err != nil {
				panic(err)
			}
			_geoMaps = append(_geoMaps, &geoMap)
		}
	}

	_links := pagination.GetLinks(_geoMaps)
	payload := geographic_maps.GetGeomapsOKBody{Geomaps: _geoMaps, Links: _links}
	return geographic_maps.NewGetGeomapsOK().WithPayload(&payload)
}

// PostGeomaps POST /geoMaps
func (c GeoMapController) PostGeomaps(params geographic_maps.PostGeomapsParams) middleware.Responder {
	geomap := params.Geomap.Geomap

	projectID, err := auth.Authenticate(params.HTTPRequest, nil)
	if err != nil {
		return geographic_maps.NewPostGeomapsDefault(403).WithPayload(utils.PolicyForbidden)
	} else {
		geomap.ProjectID = &projectID
	}

	if geomap.Provider == "" {
		return geographic_maps.NewPostGeomapsBadRequest().WithPayload(utils.MissingProvider)
	}

	// Set default values
	if err := utils.SetModelDefaults(geomap); err != nil {
		panic(err)
	}

	// Wrap insert and relations into transaction
	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO geographic_map 
				(name, default_datacenter, scope, provider, project_id)
			VALUES
				(:name, :default_datacenter, :scope, :provider, :project_id)
			RETURNING *
		`
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			return err
		}
		if err := stmt.Get(geomap, geomap); err != nil {
			return err
		}

		sql = `
			INSERT INTO geographic_map_assignment 
				(geographic_map_id, datacenter, country)
			VALUES
				(?, ?, ?)
			RETURNING *
		`
		for _, assignment := range geomap.Assignments {
			if assignment == nil {
				continue
			}
			if _, err := tx.Exec(tx.Rebind(sql), geomap.ID, assignment.Datacenter, assignment.Country); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
			return geographic_maps.NewPostGeomapsBadRequest().WithPayload(utils.NotFound)
		}

		panic(err)
	}

	_ = c.PendingSync()
	return geographic_maps.NewPostGeomapsCreated().WithPayload(&geographic_maps.PostGeomapsCreatedBody{Geomap: geomap})
}

// GetGeomapsGeoMapID GET /geoMaps/:id
func (c GeoMapController) GetGeomapsGeoMapID(params geographic_maps.GetGeomapsGeomapIDParams) middleware.Responder {
	q := sq.Select("*").
		From("geographic_map").
		Where("id = ?", params.GeomapID)

	sql, args := q.MustSql()
	var geomap models.Geomap
	if err := c.db.Get(&geomap, c.db.Rebind(sql), args...); err != nil {
		if errors.Is(err, dbsql.ErrNoRows) {
			return geographic_maps.NewGetGeomapsGeomapIDNotFound().WithPayload(utils.NotFound)
		}
		panic(err)
	}

	requestVars := map[string]string{"project_id": *geomap.ProjectID, "scope": *geomap.Scope}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return geographic_maps.NewGetGeomapsGeomapIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	sql = c.db.Rebind(`SELECT datacenter, country FROM geographic_map_assignment WHERE geographic_map_id = ?`)
	if err := c.db.Select(&geomap.Assignments, sql, geomap.ID); err != nil {
		panic(err)
	}

	return geographic_maps.NewGetGeomapsGeomapIDOK().WithPayload(&geographic_maps.GetGeomapsGeomapIDOKBody{Geomap: &geomap})
}

// PutGeomapsGeoMapID PUT /geoMaps/:id
func (c GeoMapController) PutGeomapsGeoMapID(params geographic_maps.PutGeomapsGeomapIDParams) middleware.Responder {
	return middleware.NotImplemented("operation geographic_maps.PutGeomapsGeomapID has not yet been implemented")
}

// DeleteGeomapsGeoMapID DELETE /geoMaps/:id
func (c GeoMapController) DeleteGeomapsGeoMapID(params geographic_maps.DeleteGeomapsGeomapIDParams) middleware.Responder {
	geomap := models.Geomap{ID: params.GeomapID}
	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql, args, err := sq.Select("project_id", "scope").
			From("geographic_map").
			Where("id = ?", geomap.ID).
			Suffix("FOR UPDATE").
			ToSql()
		if err != nil {
			return err
		}
		if err = tx.Get(&geomap, tx.Rebind(sql), args...); err != nil {
			return err
		}

		requestVars := map[string]string{"project_id": *geomap.ProjectID, "scope": *geomap.Scope}
		if _, err = auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
			return err
		}

		sql, args, err = sq.Update("geographic_map").
			Where("id = ?", geomap.ID).
			Set("provisioning_status", models.GeomapProvisioningStatusPENDINGDELETE).
			Set("updated_at", sq.Expr("NOW()")).
			ToSql()
		if err != nil {
			return err
		}
		if _, err = tx.Exec(c.db.Rebind(sql), args...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, dbsql.ErrNoRows) {
			return geographic_maps.NewDeleteGeomapsGeomapIDNotFound().WithPayload(utils.NotFound)
		} else if errors.Is(err, auth.ErrForbidden) {
			return geographic_maps.NewDeleteGeomapsGeomapIDDefault(403).WithPayload(utils.PolicyForbidden)
		}
		panic(err)
	}

	_ = c.PendingSync()
	return geographic_maps.NewDeleteGeomapsGeomapIDNoContent()
}

// PopulateGeoMapAssignments populates a domain instance with associated pools
func PopulateGeoMapAssignments(db *sqlx.DB, geomap *models.Geomap) error {
	// Get pool_ids associated
	sql := db.Rebind(`SELECT datacenter, country FROM geographic_map_assignment WHERE geographic_map_id = ?`)
	if err := db.Select(&geomap.Assignments, sql, geomap.ID); err != nil {
		return err
	}
	return nil
}
