/*
 *   Copyright 2022 SAP SE
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
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/restapi/operations/members"
)

func (t *SuiteTest) TestMembers() {
	mc := t.c.Members
	rr := httptest.NewRecorder()

	poolID := t.createPool()
	res := mc.GetMembersMemberID(members.GetMembersMemberIDParams{
		MemberID: "test123",
	})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound, rr.Body)

	member := members.PostMembersBody{}
	_ = member.UnmarshalBinary([]byte(`{ "member": { "name": "test", "address": "1.2.3.4", "port": 1234 } }`))
	member.Member.PoolID = &poolID

	// Write new member
	res = mc.PostMembers(members.PostMembersParams{Member: member})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusCreated, rr.Body)

	// Get all members
	res = mc.GetMembers(members.GetMembersParams{PoolID: &poolID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	membersResponse := members.GetMembersOKBody{}
	_ = membersResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(membersResponse.Members), 1, rr.Body)
	assert.Equal(t.T(), membersResponse.Members[0].ID, member.Member.ID, rr.Body)
	assert.Equal(t.T(), *membersResponse.Members[0].Name, "test", rr.Body)

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
	t.deletePool(poolID)
}
