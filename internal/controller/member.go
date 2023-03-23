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
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go-micro.dev/v4"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/members"
)

type MemberController struct {
	db *sqlx.DB
	sv micro.Service
}

// GetMembers GET /members
func (c MemberController) GetMembers(params members.GetMembersParams) middleware.Responder {
	pagination := db.NewPagination("member", params.Limit, params.Marker, params.Sort, params.PageReverse)
	// filter for pool_id, pool_id is safe and type validated
	filter := []string{fmt.Sprintf("pool_id = '%s'", params.PoolID)}
	rows, err := pagination.Query(c.db, params.HTTPRequest, filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return members.NewGetMembersBadRequest().WithPayload(utils.InvalidMarker)
		}
		if errors.Is(err, db.ErrPolicyForbidden) {
			return utils.GetPolicyForbiddenResponse()
		}
		panic(err)
	}

	//goland:noinspection GoPreferNilSlice
	var _members = []*models.Member{}
	for rows.Next() {
		member := models.Member{}
		if err := rows.StructScan(&member); err != nil {
			panic(err)
		}
		_members = append(_members, &member)
	}
	_links := pagination.GetLinks(_members, params.HTTPRequest)
	payload := members.GetMembersOKBody{Members: _members, Links: _links}
	return members.NewGetMembersOK().WithPayload(&payload)
}

// PostMembers POST /members
func (c MemberController) PostMembers(params members.PostMembersParams) middleware.Responder {
	if params.Member.Member.PoolID == nil {
		return members.NewPostMembersBadRequest().WithPayload(utils.PoolIDRequired)
	}

	member := params.Member.Member
	projectID, err := auth.ProjectScopeForRequest(params.HTTPRequest)
	if err != nil {
		panic(err)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, projectID) {
		return utils.GetPolicyForbiddenResponse()
	}

	pool := models.Pool{ID: *params.Member.Member.PoolID}
	if err := PopulatePool(c.db, &pool, []string{"project_id"}, false); err != nil || *pool.ProjectID != projectID {
		return members.NewPostMembersNotFound().WithPayload(utils.GetErrorPoolNotFound(&pool.ID))
	}
	member.PoolID = &pool.ID
	member.ProjectID = &projectID

	// Set default values
	if err := utils.SetModelDefaults(member); err != nil {
		panic(err)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		// Run insert transaction
		sql := `
			INSERT INTO member
				(name, admin_state_up, project_id, address, port, pool_id, datacenter_id)
			VALUES
				(:name, :admin_state_up, :project_id, :address, :port, :pool_id, :datacenter_id)
			RETURNING *
		`
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			panic(err)
		}
		if err := stmt.Get(member, member); err != nil {
			return err
		}

		return UpdateCascadePool(tx, pool.ID, "PENDING_UPDATE")
	}); err != nil {
		var pe *pq.Error
		if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
			return members.NewPostMembersDefault(409).WithPayload(utils.DuplicateMember)
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return members.NewPostMembersDefault(409).WithPayload(utils.DuplicateMember)
		}
		panic(err)
	}

	return members.NewPostMembersCreated().
		WithPayload(&members.PostMembersCreatedBody{Member: member})
}

// GetMembersMemberID GET /members/:id
func (c MemberController) GetMembersMemberID(params members.GetMembersMemberIDParams) middleware.Responder {
	member := models.Member{ID: params.MemberID}
	if err := PopulateMember(c.db, &member, []string{"*"}); err != nil {
		return members.NewGetMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}

	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *member.ProjectID) {
		return utils.GetPolicyForbiddenResponse()
	}
	return members.NewGetMembersMemberIDOK().
		WithPayload(&members.GetMembersMemberIDOKBody{Member: &member})
}

// PutMembersMemberID PUT /members/:id
func (c MemberController) PutMembersMemberID(params members.PutMembersMemberIDParams) middleware.Responder {
	member := models.Member{ID: params.MemberID}
	if err := PopulateMember(c.db, &member, []string{"project_id", "pool_id"}); err != nil {
		return members.NewPutMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *member.ProjectID) {
		return utils.GetPolicyForbiddenResponse()
	}

	if params.Member.Member.PoolID != nil && *params.Member.Member.PoolID != *member.PoolID {
		return members.NewPutMembersMemberIDBadRequest().WithPayload(utils.PoolIDImmutable)
	}

	params.Member.Member.ID = params.MemberID
	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
			UPDATE member SET
				name = COALESCE(:name, name),
				admin_state_up = COALESCE(:admin_state_up, admin_state_up),
				address = COALESCE(:address, address),
				port = COALESCE(:port, port),
				updated_at = NOW(),
				datacenter_id = COALESCE(:datacenter_id, datacenter_id),
				provisioning_status = 'PENDING_UPDATE'
			WHERE id = :id
		`
		if _, err := tx.NamedExec(sql, params.Member.Member); err != nil {
			panic(err)
		}
		return UpdateCascadePool(tx, *member.PoolID, "PENDING_UPDATE")
	}); err != nil {
		panic(err)
	}

	// Update member response
	if err := PopulateMember(c.db, &member, []string{"*"}); err != nil {
		panic(err)
	}

	return members.NewPutMembersMemberIDAccepted().
		WithPayload(&members.PutMembersMemberIDAcceptedBody{Member: &member})
}

// DeleteMembersMemberID DELETE /pools/:id/members/:id
func (c MemberController) DeleteMembersMemberID(params members.DeleteMembersMemberIDParams) middleware.Responder {
	member := models.Member{ID: params.MemberID}
	if err := PopulateMember(c.db, &member, []string{"project_id", "pool_id"}); err != nil {
		return members.NewDeleteMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *member.ProjectID) {
		return utils.GetPolicyForbiddenResponse()
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := tx.Rebind(`DELETE FROM member WHERE id = ?`)
		res := tx.MustExec(sql, params.MemberID)
		if deleted, _ := res.RowsAffected(); deleted != 1 {
			return EmptyResultError
		}
		return UpdateCascadePool(tx, *member.PoolID, "PENDING_UPDATE")
	}); err != nil {
		if errors.Is(err, EmptyResultError) {
			return members.NewDeleteMembersMemberIDNotFound().WithPayload(utils.NotFound)
		}
		panic(err)
	}

	return members.NewDeleteMembersMemberIDNoContent()
}

func PopulateMember(db *sqlx.DB, member *models.Member, fields []string) error {
	sql := db.Rebind(
		fmt.Sprintf(`SELECT %s FROM member WHERE id = ?`,
			strings.Join(fields, ", ")))
	if err := db.Get(member, sql, member.ID); err != nil {
		return err
	}
	return nil
}
