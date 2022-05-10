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
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/pools"
)

type PoolController struct {
	db *sqlx.DB
	sv micro.Service
}

//GetPools GET /pools
func (c PoolController) GetPools(params pools.GetPoolsParams) middleware.Responder {
	pagination := db.NewPagination("pool", params.Limit, params.Marker, params.Sort, params.PageReverse)
	rows, err := pagination.Query(c.db, params.HTTPRequest, nil)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return pools.NewGetPoolsBadRequest().WithPayload(InvalidMarker)
		}
		if errors.Is(err, db.ErrPolicyForbidden) {
			return GetPolicyForbiddenResponse()
		}
		panic(err)
	}

	//zero-length slice used because we want [] via json encoder, nil encodes null
	//goland:noinspection GoPreferNilSlice
	var _pools = []*models.Pool{}
	for rows.Next() {
		pool := models.Pool{Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
		if err := rows.StructScan(&pool); err != nil {
			panic(err)
		}
		if err := PopulatePoolDomains(c.db, &pool); err != nil {
			panic(err)
		}
		_pools = append(_pools, &pool)
	}
	_links := pagination.GetLinks(_pools, params.HTTPRequest)
	payload := pools.GetPoolsOKBody{Pools: _pools, Links: _links}
	return pools.NewGetPoolsOK().WithPayload(&payload)
}

//PostPools POST /pools
func (c PoolController) PostPools(params pools.PostPoolsParams) middleware.Responder {
	pool := params.Pool.Pool
	projectID, err := auth.ProjectScopeForRequest(params.HTTPRequest)
	if err != nil {
		panic(err)
	}

	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, projectID) {
		return GetPolicyForbiddenResponse()
	}
	pool.ProjectID = &projectID

	// Set default values
	if err := utils.SetModelDefaults(pool); err != nil {
		panic(err)
	}
	pool.Members = []strfmt.UUID{}
	pool.Domains = []strfmt.UUID{}

	sql := `
		INSERT INTO pool
		    (name, admin_state_up, project_id)
		VALUES
		    (:name, :admin_state_up, :project_id)
		RETURNING *
	`
	stmt, _ := c.db.PrepareNamed(sql)
	if err := stmt.Get(pool, pool); err != nil {
		panic(err)
	}
	return pools.NewPostPoolsCreated().WithPayload(&pools.PostPoolsCreatedBody{Pool: pool})
}

//GetPoolsPoolID GET /pools/:id
func (c PoolController) GetPoolsPoolID(params pools.GetPoolsPoolIDParams) middleware.Responder {
	//zero-length slice used because we want [] via json encoder, nil encodes null
	pool := models.Pool{ID: params.PoolID, Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
	if err := PopulatePool(c.db, &pool, []string{"*"}, true); err != nil {
		return pools.NewGetPoolsPoolIDNotFound().WithPayload(NotFound)
	}

	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *pool.ProjectID) {
		return GetPolicyForbiddenResponse()
	}
	return pools.NewGetPoolsPoolIDOK().WithPayload(&pools.GetPoolsPoolIDOKBody{Pool: &pool})
}

//PutPoolsPoolID PUT /pools/:id
func (c PoolController) PutPoolsPoolID(params pools.PutPoolsPoolIDParams) middleware.Responder {
	//zero-length slice used because we want [] via json encoder, nil encodes null
	pool := models.Pool{ID: params.PoolID, Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
	if err := PopulatePool(c.db, &pool, []string{"id", "project_id"}, false); err != nil {
		return pools.NewPutPoolsPoolIDNotFound().WithPayload(NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *pool.ProjectID) {
		return GetPolicyForbiddenResponse()
	}

	params.Pool.Pool.ID = params.PoolID
	sql := `
		UPDATE pool SET
			name = COALESCE(:name, name),
			admin_state_up = COALESCE(:admin_state_up, admin_state_up),
		    updated_at = NOW()
		WHERE id = :id
	`
	if _, err := c.db.NamedExec(sql, params.Pool.Pool); err != nil {
		panic(err)
	}

	// Update pool response
	if err := PopulatePool(c.db, &pool, []string{"*"}, true); err != nil {
		panic(err)
	}

	return pools.NewPutPoolsPoolIDAccepted().WithPayload(&pools.PutPoolsPoolIDAcceptedBody{Pool: &pool})
}

//DeletePoolsPoolID DELETE /pools/:id
func (c PoolController) DeletePoolsPoolID(params pools.DeletePoolsPoolIDParams) middleware.Responder {
	//zero-length slice used because we want [] via json encoder, nil encodes null
	pool := models.Pool{ID: params.PoolID, Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
	if err := PopulatePool(c.db, &pool, []string{"id", "project_id"}, false); err != nil {
		return pools.NewDeletePoolsPoolIDNotFound().WithPayload(NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *pool.ProjectID) {
		return GetPolicyForbiddenResponse()
	}

	sql := `DELETE FROM pool WHERE id = ?`
	res := c.db.MustExec(sql, params.PoolID)
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		return pools.NewDeletePoolsPoolIDNotFound().WithPayload(NotFound)
	}
	return pools.NewDeletePoolsPoolIDNoContent()
}

//PopulatePoolMembers populates a pool instance with associated members
func PopulatePoolMembers(db *sqlx.DB, pool *models.Pool) error {
	// Get datacenter_id's associated
	sql := `SELECT id FROM member WHERE pool_id = ?`
	if err := db.Select(&pool.Members, sql, pool.ID); err != nil {
		return err
	}
	return nil
}

//PopulatePoolDomains populates a pool instance with associated domains
func PopulatePoolDomains(db *sqlx.DB, pool *models.Pool) error {
	//Get domain_id's associated
	sql := `SELECT domain_id FROM domain_pool_relation WHERE pool_id = ?`
	err := db.Select(&pool.Domains, sql, pool.ID)
	if err != nil {
		return err
	}
	return nil
}

//PopulatePoolMonitors populates a pool instance with associated monitors
func PopulatePoolMonitors(db *sqlx.DB, pool *models.Pool) error {
	//Get domain_id's associated
	sql := `SELECT id FROM monitor WHERE pool_id = ?`
	err := db.Select(&pool.Monitors, sql, pool.ID)
	if err != nil {
		return err
	}
	return nil
}

//PopulatePool populates attributes of a pool instance based on it's ID
func PopulatePool(db *sqlx.DB, pool *models.Pool, fields []string, fullyPopulate bool) error {
	//Get pool
	sql := fmt.Sprintf(`SELECT %s FROM pool WHERE id = ?`, strings.Join(fields, ", "))
	if err := db.Get(pool, sql, pool.ID); err != nil {
		return err
	}
	if fullyPopulate {
		if err := PopulatePoolMembers(db, pool); err != nil {
			return err
		}
		if err := PopulatePoolDomains(db, pool); err != nil {
			return err
		}
		if err := PopulatePoolMonitors(db, pool); err != nil {
			return err
		}
	}
	return nil
}
