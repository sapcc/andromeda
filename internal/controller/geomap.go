/*
 *   Copyright 2020 SAP SE
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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/geographic_maps"
)

type GeoMapController struct {
	db *sqlx.DB
	sv micro.Service
}

// GetGeomaps GET /geoMaps
func (c GeoMapController) GetGeomaps(params geographic_maps.GetGeomapsParams) middleware.Responder {
	pagination := db.Pagination(params)
	rows, err := pagination.Query(c.db, "SELECT * FROM geographic_map", nil)
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

	_ = PendingSync(c.sv)
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
			Set("provisioning_status", "PENDING_DELETE").
			Set("updated_at", sq.Expr("NOW()")).
			ToSql()
		if err != nil {
			return err
		}
		if _, err = c.db.Exec(c.db.Rebind(sql), args...); err != nil {
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

	_ = PendingSync(c.sv)
	return geographic_maps.NewDeleteGeomapsGeomapIDNoContent()
}

// PopulateDomainPools populates a domain instance with associated pools
func PopulateGeoMapAssignments(db *sqlx.DB, geomap *models.Geomap) error {
	// Get pool_ids associated
	sql := db.Rebind(`SELECT datacenter, country FROM geographic_map_assignment WHERE geographic_map_id = ?`)
	if err := db.Select(&geomap.Assignments, sql, geomap.ID); err != nil {
		return err
	}
	return nil
}
