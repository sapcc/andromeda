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
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go-micro.dev/v4"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/domains"
)

type DomainController struct {
	db *sqlx.DB
	sv micro.Service
}

// GetDomains GET /domains
func (c DomainController) GetDomains(params domains.GetDomainsParams) middleware.Responder {
	filter := make(map[string]any, 0)
	if projectId, err := auth.Authenticate(params.HTTPRequest); err != nil {
		return domains.NewGetDomainsDefault(403).WithPayload(utils.PolicyForbidden)
	} else if projectId != "" {
		filter["project_id"] = projectId
	}

	pagination := db.Pagination(params)
	rows, err := pagination.Query(c.db, "SELECT * FROM domain", filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return domains.NewGetDomainsDefault(400).WithPayload(utils.InvalidMarker)
		}
		panic(err)
	}

	//goland:noinspection GoPreferNilSlice
	var _domains = []*models.Domain{}
	for rows.Next() {
		domain := models.Domain{Pools: []strfmt.UUID{}}
		if err := rows.StructScan(&domain); err != nil {
			panic(err)
		}
		if *domain.Provider == "akamai" {
			cname := strfmt.Hostname(fmt.Sprintf("%s.%s", domain.Fqdn.String(), config.Global.AkamaiConfig.Domain))
			domain.CnameTarget = &cname
		}
		if err := PopulateDomainPools(c.db, &domain); err != nil {
			panic(err)
		}
		_domains = append(_domains, &domain)
	}
	_links := pagination.GetLinks(_domains)
	payload := domains.GetDomainsOKBody{Domains: _domains, Links: _links}
	return domains.NewGetDomainsOK().WithPayload(&payload)
}

// PostDomains POST /domains
func (c DomainController) PostDomains(params domains.PostDomainsParams) middleware.Responder {
	domain := params.Domain.Domain

	projectID, err := auth.Authenticate(params.HTTPRequest)
	if err != nil {
		return domains.NewPostDomainsDefault(403).WithPayload(utils.PolicyForbidden)
	} else {
		domain.ProjectID = &projectID
	}

	if domain.Pools == nil {
		domain.Pools = []strfmt.UUID{}
	}

	// Set default values
	if err := utils.SetModelDefaults(domain); err != nil {
		panic(err)
	}

	// Wrap insert and relations into transaction
	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO domain 
				(name, fqdn, record_type, mode, admin_state_up, provider, project_id)
			VALUES
				(:name, :fqdn, :record_type, :mode, :admin_state_up, :provider, :project_id)
			RETURNING *
		`
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			return err
		}
		if err := stmt.Get(domain, domain); err != nil {
			return err
		}

		// Add pool/domain relationship
		if len(params.Domain.Domain.Pools) > 0 {
			if _, err := insertDomainPoolRelations(tx, domain.ID, projectID, domain.Pools); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		var rnfError *utils.ResourcesNotFoundError
		if errors.As(err, &rnfError) {
			errMsg := "Invalid value for 'pools': " + rnfError.Error()
			return domains.NewPostDomainsDefault(400).WithPayload(
				&models.Error{Code: 400, Message: errMsg})
		}
		var pe *pq.Error
		if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
			return domains.NewPostDomainsDefault(409).WithPayload(utils.DuplicateDomain)
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return domains.NewPostDomainsDefault(409).WithPayload(utils.DuplicateDomain)
		}
		panic(err)
	}

	return domains.NewPostDomainsCreated().WithPayload(&domains.PostDomainsCreatedBody{Domain: domain})
}

// GetDomainsDomainID GET /domains/:id
func (c DomainController) GetDomainsDomainID(params domains.GetDomainsDomainIDParams) middleware.Responder {
	// Get domain
	domain := models.Domain{ID: params.DomainID, Pools: []strfmt.UUID{}}

	if err := PopulateDomain(c.db, &domain, []string{"*"}); err != nil {
		return domains.NewGetDomainsDomainIDNotFound().WithPayload(utils.NotFound)
	}
	if projectID, err := auth.Authenticate(params.HTTPRequest); err != nil || projectID != *domain.ProjectID {
		return domains.NewGetDomainsDomainIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	return domains.NewGetDomainsDomainIDOK().WithPayload(&domains.GetDomainsDomainIDOKBody{Domain: &domain})
}

// PutDomainsDomainID PUT /domains/:id
func (c DomainController) PutDomainsDomainID(params domains.PutDomainsDomainIDParams) middleware.Responder {
	domain := models.Domain{Pools: []strfmt.UUID{}, ID: params.DomainID}
	if err := PopulateDomain(c.db, &domain, []string{"*"}); err != nil {
		return domains.NewPutDomainsDomainIDNotFound().WithPayload(utils.NotFound)
	}
	if projectID, err := auth.Authenticate(params.HTTPRequest); err != nil || projectID != *domain.ProjectID {
		return domains.NewPutDomainsDomainIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	if params.Domain.Domain.Provider != nil && *params.Domain.Domain.Provider != *domain.Provider {
		// cannot change provider
		return domains.NewPutDomainsDomainIDNotFound().WithPayload(utils.ProviderUnchangeable)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		// Populate args
		if params.Domain.Domain.Pools != nil {
			var existingPoolRefs []strfmt.UUID
			sql := tx.Rebind(`SELECT pool_id FROM domain_pool_relation WHERE domain_id = ? FOR UPDATE`)
			if err := tx.Select(&existingPoolRefs, sql, params.DomainID); err != nil {
				return err
			}

			poolsRemoved := utils.UUIDsDifference(existingPoolRefs, params.Domain.Domain.Pools)
			poolsAdded := utils.UUIDsDifference(params.Domain.Domain.Pools, existingPoolRefs)
			if poolsRemoved != nil {
				if _, err := removeDomainPoolRelations(tx, params.DomainID, poolsRemoved); err != nil {
					return err
				}
			}

			if poolsAdded != nil {
				if _, err := insertDomainPoolRelations(tx, domain.ID, *domain.ProjectID, poolsAdded); err != nil {
					return err
				}
			}
		}

		// Update
		params.Domain.Domain.ID = params.DomainID
		sql := `
			UPDATE domain SET
				name = COALESCE(:name, name),
				admin_state_up = COALESCE(:admin_state_up, admin_state_up),
				fqdn = COALESCE(:fqdn, fqdn), 
				mode = COALESCE(:mode, mode), 
				record_type = COALESCE(:record_type, record_type), 
			    provisioning_status = 'PENDING_UPDATE',
				updated_at = NOW()
			WHERE id = :id
		`
		if _, err := tx.NamedExec(sql, params.Domain.Domain); err != nil {
			return err
		}

		return nil
	}); err != nil {
		var rnfError *utils.ResourcesNotFoundError
		if errors.As(err, &rnfError) {
			errMsg := "Invalid value for 'pools': " + rnfError.Error()
			return domains.NewPutDomainsDomainIDBadRequest().WithPayload(
				&models.Error{Code: 400, Message: errMsg})
		}
		if errors.Is(err, dbsql.ErrNoRows) {
			return domains.NewPutDomainsDomainIDNotFound().WithPayload(utils.NotFound)
		}
		// Unknown Error
		panic(err)
	}

	// Update domain response
	if err := PopulateDomain(c.db, &domain, []string{"*"}); err != nil {
		panic(err)
	}

	return domains.NewPutDomainsDomainIDAccepted().WithPayload(&domains.PutDomainsDomainIDAcceptedBody{Domain: &domain})
}

// DeleteDomainsDomainID DELETE /domains/:id
func (c DomainController) DeleteDomainsDomainID(params domains.DeleteDomainsDomainIDParams) middleware.Responder {
	domain := models.Domain{ID: params.DomainID}
	if err := PopulateDomain(c.db, &domain, []string{"id", "project_id"}); err != nil {
		return domains.NewDeleteDomainsDomainIDNotFound().WithPayload(utils.NotFound)
	}
	if projectID, err := auth.Authenticate(params.HTTPRequest); err != nil || projectID != *domain.ProjectID {
		return domains.NewDeleteDomainsDomainIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	sql := c.db.Rebind(`UPDATE domain SET provisioning_status = 'PENDING_DELETE', updated_at = NOW() WHERE id = ?`)
	res := c.db.MustExec(sql, params.DomainID)
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		return domains.NewDeleteDomainsDomainIDNotFound().WithPayload(utils.NotFound)
	}
	return domains.NewDeleteDomainsDomainIDNoContent()
}

// removeDomainPoolRelations removes pools associated to a domain inside a transaction
func removeDomainPoolRelations(tx *sqlx.Tx, domainID strfmt.UUID, poolIDs []strfmt.UUID) (*strfmt.UUID, error) {
	sql := tx.Rebind(`DELETE FROM domain_pool_relation WHERE domain_id = ? and pool_id = ?`)
	for _, poolID := range poolIDs {
		if _, err := tx.Exec(sql, domainID, poolID); err != nil {
			return &poolID, err
		}
	}
	return nil, UpdateCascadeDomain(tx, domainID, "PENDING_UPDATE")
}

// insertDomainPoolRelations adds pools associated to a domain inside a transaction
func insertDomainPoolRelations(tx *sqlx.Tx, domainID strfmt.UUID, projectID string, poolIDs []strfmt.UUID) (*strfmt.UUID, error) {
	//check pools are belonging to same project
	sql := `SELECT id FROM pool WHERE id IN (?) AND project_id = ?;`
	query, args, err := sqlx.In(sql, poolIDs, projectID)
	if err != nil {
		return nil, err
	}
	query = tx.Rebind(query)
	var validPoolsFound []strfmt.UUID
	if err := tx.Select(&validPoolsFound, query, args...); err != nil {
		return nil, err
	}
	if len(validPoolsFound) != len(poolIDs) {
		missingPools := utils.UUIDsDifference(poolIDs, validPoolsFound)
		return nil, &utils.ResourcesNotFoundError{Ids: missingPools, Resource: "Pool"}
	}

	for _, poolID := range poolIDs {
		if _, err := tx.NamedExec(
			"INSERT INTO domain_pool_relation (domain_id, pool_id) VALUES (:domain_id, :pool_id)",
			map[string]interface{}{
				"domain_id": domainID,
				"pool_id":   poolID,
			},
		); err != nil {
			return &poolID, err
		}
	}
	return nil, UpdateCascadeDomain(tx, domainID, "PENDING_UPDATE")
}

// PopulateDomainPools populates a domain instance with associated pools
func PopulateDomainPools(db *sqlx.DB, domain *models.Domain) error {
	// Get pool_ids associated
	sql := db.Rebind(`SELECT pool_id FROM domain_pool_relation WHERE domain_id = ?`)
	if err := db.Select(&domain.Pools, sql, domain.ID); err != nil {
		return err
	}
	return nil
}

// PopulateDomain populates attributes of a domain instance based on it's ID
func PopulateDomain(db *sqlx.DB, domain *models.Domain, fields []string) error {
	// Get domain
	// zero-length slice used because we want [] via json encoder, nil encodes null
	sql := db.Rebind(fmt.Sprintf(`SELECT %s FROM domain WHERE id = ?`, strings.Join(fields, ", ")))
	if err := db.Get(domain, sql, domain.ID); err != nil {
		return err
	}
	if err := PopulateDomainPools(db, domain); err != nil {
		return err
	}
	return nil
}

func UpdateCascadeDomain(tx *sqlx.Tx, domainID strfmt.UUID, provisioningStatus string) error {
	// Pending Domain
	sql := fmt.Sprintf(`UPDATE domain SET provisioning_status = '%s' WHERE id = ?`, provisioningStatus)
	if _, err := tx.Exec(tx.Rebind(sql), domainID); err != nil {
		return err
	}
	return nil
}
