// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/pools"
)

var (
	errPoolImmutable = errors.New("pool is immutable")
)

type PoolController struct {
	CommonController
}

// GetPools GET /pools
func (c PoolController) GetPools(params pools.GetPoolsParams) middleware.Responder {
	filter := make(map[string]any)
	pagination := db.Pagination{
		HTTPRequest: params.HTTPRequest,
		Limit:       params.Limit,
		Marker:      params.Marker,
		PageReverse: params.PageReverse,
		Sort:        params.Sort,
	}
	sql := `SELECT * FROM pool`
	if params.DomainID != nil {
		filter["domain_id"] = *params.DomainID
		sql = `SELECT 
    		pool.id AS id, 
    		pool.name AS name, 
    		pool.provisioning_status AS provisioning_status,
    		pool.status AS status,
    		pool.admin_state_up AS admin_state_up,
    		pool.project_id AS project_id,
    		pool.created_at AS created_at,
    		pool.updated_at AS updated_at
		FROM pool JOIN domain_pool_relation ON pool.id = domain_pool_relation.pool_id`
	}
	if params.PoolID != nil {
		filter["id"] = *params.PoolID
	}
	rows, err := pagination.Query(c.db, sql, filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return pools.NewGetPoolsBadRequest().WithPayload(utils.InvalidMarker)
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

		// Filter result based on policy
		requestVars := map[string]string{"project_id": *pool.ProjectID}
		if _, err = auth.Authenticate(params.HTTPRequest, requestVars); err == nil {
			if err := PopulatePoolDomains(c.db, &pool); err != nil {
				panic(err)
			}
			if err := PopulatePoolMembers(c.db, &pool); err != nil {
				panic(err)
			}
			if err := PopulatePoolMonitors(c.db, &pool); err != nil {
				panic(err)
			}
			_pools = append(_pools, &pool)
		}
	}
	_links := pagination.GetLinks(_pools)
	payload := pools.GetPoolsOKBody{Pools: _pools, Links: _links}
	return pools.NewGetPoolsOK().WithPayload(&payload)
}

// PostPools POST /pools
func (c PoolController) PostPools(params pools.PostPoolsParams) middleware.Responder {
	pool := params.Pool.Pool

	projectID, err := auth.Authenticate(params.HTTPRequest, nil)
	if err != nil {
		return pools.NewPostPoolsDefault(403).WithPayload(utils.PolicyForbidden)
	} else {
		pool.ProjectID = &projectID
	}

	// Set default values
	if err := utils.SetModelDefaults(pool); err != nil {
		panic(err)
	}
	pool.Members = []strfmt.UUID{}

	// Wrap insert and relations into transaction
	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO pool
				(name, admin_state_up, project_id)
			VALUES
				(:name, :admin_state_up, :project_id)
			RETURNING *
		`

		stmt, _ := tx.PrepareNamed(sql)
		if err := stmt.Get(pool, pool); err != nil {
			return err
		}

		for _, domainId := range params.Pool.Pool.Domains {
			if _, err := insertDomainPoolRelations(tx, domainId, projectID, []strfmt.UUID{pool.ID}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	_ = PendingSync(c.rpc)
	return pools.NewPostPoolsCreated().WithPayload(&pools.PostPoolsCreatedBody{Pool: pool})
}

// GetPoolsPoolID GET /pools/:id
func (c PoolController) GetPoolsPoolID(params pools.GetPoolsPoolIDParams) middleware.Responder {
	//zero-length slice used because we want [] via json encoder, nil encodes null
	pool := models.Pool{ID: params.PoolID, Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
	if err := PopulatePool(c.db, &pool, []string{"*"}, true); err != nil {
		return pools.NewGetPoolsPoolIDNotFound().WithPayload(utils.NotFound)
	}

	requestVars := map[string]string{"project_id": *pool.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return pools.NewGetPoolsPoolIDDefault(403).WithPayload(utils.PolicyForbidden)
	}
	return pools.NewGetPoolsPoolIDOK().WithPayload(&pools.GetPoolsPoolIDOKBody{Pool: &pool})
}

// PutPoolsPoolID PUT /pools/:id
func (c PoolController) PutPoolsPoolID(params pools.PutPoolsPoolIDParams) middleware.Responder {
	//zero-length slice used because we want [] via json encoder, nil encodes null
	pool := models.Pool{ID: params.PoolID, Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
	if err := PopulatePool(c.db, &pool, []string{"id", "project_id"}, true); err != nil {
		return pools.NewPutPoolsPoolIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *pool.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return pools.NewPutPoolsPoolIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		if pool.Domains != nil {
			// Check if domains are in progress of being updated
			if isPoolImmutable(&pool, tx) {
				return errPoolImmutable
			}
		}

		// Populate args
		if params.Pool.Pool.Domains != nil {
			var existingDomainRefs []strfmt.UUID
			sql := tx.Rebind(`SELECT domain_id FROM domain_pool_relation WHERE pool_id = ? FOR UPDATE`)
			if err := tx.Select(&existingDomainRefs, sql, params.PoolID); err != nil {
				return err
			}

			domainsRemoved := utils.UUIDsDifference(existingDomainRefs, params.Pool.Pool.Domains)
			domainsAdded := utils.UUIDsDifference(params.Pool.Pool.Domains, existingDomainRefs)
			for _, domain := range domainsRemoved {
				if _, err := removeDomainPoolRelations(tx, domain, []strfmt.UUID{params.PoolID}); err != nil {
					return err
				}
			}
			for _, domain := range domainsAdded {
				if _, err := insertDomainPoolRelations(tx, domain, *pool.ProjectID, []strfmt.UUID{params.PoolID}); err != nil {
					return err
				}
			}
		}

		params.Pool.Pool.ID = params.PoolID
		sql := `
			UPDATE pool SET
				name = COALESCE(:name, name),
				admin_state_up = COALESCE(:admin_state_up, admin_state_up),
				updated_at = NOW()
			WHERE id = :id
		`
		if _, err := tx.NamedExec(sql, params.Pool.Pool); err != nil {
			return err
		}

		return UpdateCascadePool(tx, params.PoolID, "PENDING_UPDATE")
	}); err != nil {
		if errors.Is(err, errPoolImmutable) {
			return pools.NewPutPoolsPoolIDConflict().WithPayload(utils.GetErrorImmutable("pool", "domains"))
		}
		panic(err)
	}

	// Update pool response
	if err := PopulatePool(c.db, &pool, []string{"*"}, true); err != nil {
		panic(err)
	}

	_ = PendingSync(c.rpc)
	return pools.NewPutPoolsPoolIDAccepted().WithPayload(&pools.PutPoolsPoolIDAcceptedBody{Pool: &pool})
}

// DeletePoolsPoolID DELETE /pools/:id
func (c PoolController) DeletePoolsPoolID(params pools.DeletePoolsPoolIDParams) middleware.Responder {
	// zero-length slice used because we want [] via json encoder, nil encodes null
	pool := models.Pool{ID: params.PoolID, Members: []strfmt.UUID{}, Domains: []strfmt.UUID{}}
	if err := PopulatePool(c.db, &pool, []string{"id", "project_id"}, true); err != nil {
		return pools.NewDeletePoolsPoolIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *pool.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return pools.NewDeletePoolsPoolIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		if isPoolImmutable(&pool, tx) {
			return errPoolImmutable
		}

		if len(pool.Domains) > 0 {
			return UpdateCascadePool(tx, params.PoolID, "PENDING_DELETE")
		}

		_, err := tx.Exec(tx.Rebind(`DELETE FROM pool WHERE id = ?`), params.PoolID)
		return err
	}); err != nil {
		if errors.Is(err, errPoolImmutable) {
			return pools.NewDeletePoolsPoolIDConflict().WithPayload(
				utils.GetErrorImmutable("pool", "domains"))
		}
		panic(err)
	}

	_ = PendingSync(c.rpc)
	return pools.NewDeletePoolsPoolIDNoContent()
}

// isPoolImmutable checks if a pool is immutable because related domains are being updated/deleted.
// Should be called with a fully populated pool model instance.
func isPoolImmutable(pool *models.Pool, tx *sqlx.Tx) bool {
	// sqlx.In does not support []strfmt.UUID, so we convert to []string
	var domains = make([]string, 0, len(pool.Domains))
	for _, domain := range pool.Domains {
		domains = append(domains, domain.String())
	}

	if len(domains) == 0 {
		return false
	}

	// Check if related domains are in progress of being updated/deleted
	query := `SELECT count(id) FROM domain WHERE provisioning_status != 'ACTIVE' AND id IN (?)`
	sql, args, err := sqlx.In(query, domains)
	if err != nil {
		panic(err)
	}

	var countNonActive int
	rows, err := tx.Query(tx.Rebind(sql), args...)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		if err = rows.Scan(&countNonActive); err != nil {
			panic(err)
		}
	}

	return countNonActive > 0
}

// PopulatePoolMembers populates a pool instance with associated members
func PopulatePoolMembers(db *sqlx.DB, pool *models.Pool) error {
	// Get datacenter_id's associated
	sql := db.Rebind(`SELECT id FROM member WHERE pool_id = ?`)
	if err := db.Select(&pool.Members, sql, pool.ID); err != nil {
		return err
	}
	return nil
}

// PopulatePoolDomains populates a pool instance with associated domains
func PopulatePoolDomains(db *sqlx.DB, pool *models.Pool) error {
	//Get domain_id's associated
	sql := db.Rebind(`SELECT domain_id FROM domain_pool_relation WHERE pool_id = ?`)
	err := db.Select(&pool.Domains, sql, pool.ID)
	if err != nil {
		return err
	}
	return nil
}

// PopulatePoolMonitors populates a pool instance with associated monitors
func PopulatePoolMonitors(db *sqlx.DB, pool *models.Pool) error {
	//Get domain_id's associated
	sql := db.Rebind(`SELECT id FROM monitor WHERE pool_id = ?`)
	err := db.Select(&pool.Monitors, sql, pool.ID)
	if err != nil {
		return err
	}
	return nil
}

// PopulatePool populates attributes of a pool instance based on it's ID
func PopulatePool(db *sqlx.DB, pool *models.Pool, fields []string, fullyPopulate bool) error {
	//Get pool
	sql := db.Rebind(fmt.Sprintf(`SELECT %s FROM pool WHERE id = ?`, strings.Join(fields, ", ")))
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

func UpdateCascadePool(tx *sqlx.Tx, poolID strfmt.UUID, provisioningStatus string) error {
	var sql string
	// Pending Domain
	if tx.DriverName() == "mysql" {
		sql = tx.Rebind(`
			UPDATE domain
			JOIN domain_pool_relation dpr on domain.id = dpr.domain_id
			SET provisioning_status = 'PENDING_UPDATE', updated_at = NOW()
			WHERE dpr.pool_id = ?`)
	} else {
		sql = tx.Rebind(`
			UPDATE domain
			SET provisioning_status = 'PENDING_UPDATE', updated_at = NOW()
			FROM domain_pool_relation
			WHERE domain.id = domain_pool_relation.domain_id AND domain_pool_relation.pool_id = ?`)
	}
	if _, err := tx.Exec(sql, poolID); err != nil {
		return err
	}

	sql = fmt.Sprintf(`UPDATE pool SET provisioning_status = '%s' WHERE id = ?`, provisioningStatus)
	if _, err := tx.Exec(tx.Rebind(sql), poolID); err != nil {
		return err
	}

	return nil
}
