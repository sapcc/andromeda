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
	"github.com/lib/pq"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
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

//GetMembers GET /pools/:id/members
func (c MemberController) GetMembers(params members.GetPoolsPoolIDMembersParams) middleware.Responder {
	pagination := db.NewPagination("member", params.Limit, params.Marker, params.Sort, params.PageReverse)
	//filter for pool_id, pool_id is safe and type validated
	filter := []string{fmt.Sprintf("pool_id = '%s'", params.PoolID)}
	rows, err := pagination.Query(c.db, params.HTTPRequest, filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return members.NewGetPoolsPoolIDMembersBadRequest().WithPayload(utils.InvalidMarker)
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
	payload := members.GetPoolsPoolIDMembersOKBody{Members: _members, Links: _links}
	return members.NewGetPoolsPoolIDMembersOK().WithPayload(&payload)
}

//PostMembers POST /pools/:id/members
func (c MemberController) PostMembers(params members.PostPoolsPoolIDMembersParams) middleware.Responder {
	member := params.Member.Member
	projectID, err := auth.ProjectScopeForRequest(params.HTTPRequest)
	if err != nil {
		panic(err)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, projectID) {
		return utils.GetPolicyForbiddenResponse()
	}

	pool := models.Pool{ID: params.PoolID}
	if err := PopulatePool(c.db, &pool, []string{"project_id"}, false); err != nil || *pool.ProjectID != projectID {
		return members.NewPostPoolsPoolIDMembersNotFound().WithPayload(utils.GetErrorPoolNotFound(&params.PoolID))
	}
	member.PoolID = params.PoolID
	member.ProjectID = &projectID

	// Set default values
	if err := utils.SetModelDefaults(member); err != nil {
		panic(err)
	}

	sql := `
		INSERT INTO member
		    (name, admin_state_up, project_id, address, port, pool_id, datacenter_id)
		VALUES
		    (:name, :admin_state_up, :project_id, :address, :port, :pool_id, :datacenter_id)
		RETURNING *
	`
	stmt, err := c.db.PrepareNamed(sql)
	if err != nil {
		panic(err)
	}
	if err := stmt.Get(member, member); err != nil {
		var pe *pq.Error
		if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
			return members.NewPostPoolsPoolIDMembersDefault(409).WithPayload(utils.DuplicateMember)
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return members.NewPostPoolsPoolIDMembersDefault(409).WithPayload(utils.DuplicateMember)
		}
		panic(err)
	}
	return members.NewPostPoolsPoolIDMembersCreated().
		WithPayload(&members.PostPoolsPoolIDMembersCreatedBody{Member: member})
}

//GetMembersMemberID GET /pools/:id/members/:id
func (c MemberController) GetMembersMemberID(params members.GetPoolsPoolIDMembersMemberIDParams) middleware.Responder {
	member := models.Member{ID: params.MemberID, PoolID: params.PoolID}
	if err := PopulateMember(c.db, &member, []string{"*"}); err != nil {
		return members.NewGetPoolsPoolIDMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}

	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *member.ProjectID) {
		return utils.GetPolicyForbiddenResponse()
	}
	return members.NewGetPoolsPoolIDMembersMemberIDOK().
		WithPayload(&members.GetPoolsPoolIDMembersMemberIDOKBody{Member: &member})
}

//PutMembersMemberID PUT /pools/:id/members/:id
func (c MemberController) PutMembersMemberID(params members.PutPoolsPoolIDMembersMemberIDParams) middleware.Responder {
	member := models.Member{ID: params.MemberID, PoolID: params.PoolID}
	if err := PopulateMember(c.db, &member, []string{"project_id"}); err != nil {
		return members.NewPutPoolsPoolIDMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *member.ProjectID) {
		return utils.GetPolicyForbiddenResponse()
	}

	params.Member.Member.ID = params.MemberID
	params.Member.Member.PoolID = params.PoolID
	sql := `
		UPDATE member SET
			name = COALESCE(:name, name),
			admin_state_up = COALESCE(:admin_state_up, admin_state_up),
			address = COALESCE(:address, address),
			port = COALESCE(:port, port),
		    updated_at = NOW()
		WHERE id = :id
	`
	if _, err := c.db.NamedExec(sql, params.Member.Member); err != nil {
		panic(err)
	}

	// Update member response
	if err := PopulateMember(c.db, &member, []string{"*"}); err != nil {
		panic(err)
	}

	return members.NewPutPoolsPoolIDMembersMemberIDAccepted().
		WithPayload(&members.PutPoolsPoolIDMembersMemberIDAcceptedBody{Member: &member})
}

//DeleteMembersMemberID DELETE /pools/:id/members/:id
func (c MemberController) DeleteMembersMemberID(params members.DeletePoolsPoolIDMembersMemberIDParams) middleware.Responder {
	member := models.Member{ID: params.MemberID, PoolID: params.PoolID}
	if err := PopulateMember(c.db, &member, []string{"project_id"}); err != nil {
		return members.NewDeletePoolsPoolIDMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *member.ProjectID) {
		return utils.GetPolicyForbiddenResponse()
	}

	sql := c.db.Rebind(`DELETE FROM member WHERE id = ?`)
	res := c.db.MustExec(sql, params.MemberID)
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		members.NewDeletePoolsPoolIDMembersMemberIDNotFound().WithPayload(utils.NotFound)
	}

	return members.NewDeletePoolsPoolIDMembersMemberIDNoContent()
}

func PopulateMember(db *sqlx.DB, member *models.Member, fields []string) error {
	sql := db.Rebind(fmt.Sprintf(`SELECT %s FROM member WHERE pool_id = ? AND id = ?`, strings.Join(fields, ", ")))
	if err := db.Get(member, sql, member.PoolID, member.ID); err != nil {
		return err
	}
	return nil
}
