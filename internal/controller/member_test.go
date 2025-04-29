// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/restapi/operations/members"
)

func (t *SuiteTest) TestMembers() {
	mc := t.c.Members
	rr := httptest.NewRecorder()

	poolID := t.createPool(nil)
	res := mc.GetMembersMemberID(members.GetMembersMemberIDParams{
		MemberID: "test123",
	})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusNotFound, rr.Code, rr.Body)

	member := members.PostMembersBody{}
	_ = member.UnmarshalBinary([]byte(`{ "member": { "name": "test", "address": "1.2.3.4", "port": 1234 } }`))
	member.Member.PoolID = &poolID

	// Write new member
	res = mc.PostMembers(members.PostMembersParams{Member: member})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	// Get all members
	res = mc.GetMembers(members.GetMembersParams{PoolID: &poolID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	membersResponse := members.GetMembersOKBody{}
	_ = membersResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(membersResponse.Members), 1, rr.Body)
	assert.Equal(t.T(), membersResponse.Members[0].ID, member.Member.ID, rr.Body)
	assert.Equal(t.T(), *membersResponse.Members[0].Name, "test", rr.Body)

	// Fetch all members without pool filter
	res = mc.GetMembers(members.GetMembersParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	membersResponse = members.GetMembersOKBody{}
	_ = membersResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(membersResponse.Members), 1, rr.Body)

	// Get specific member
	res = mc.GetMembersMemberID(members.GetMembersMemberIDParams{
		MemberID: member.Member.ID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	// Delete specific member
	res = mc.DeleteMembersMemberID(members.DeleteMembersMemberIDParams{
		MemberID: member.Member.ID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusNoContent, rr.Code, rr.Body)

	// Cleanup
	if _, err := t.db.Exec("DELETE FROM pool"); err != nil {
		t.FailNow(err.Error())
	}
}
