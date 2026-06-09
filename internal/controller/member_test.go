// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/domains"
	"github.com/sapcc/andromeda/restapi/operations/members"
)

func (t *SuiteTest) createF5Domain() strfmt.UUID {
	fqdn := strfmt.Hostname("f5test.com")
	domain := domains.PostDomainsBody{
		Domain: &models.Domain{
			Fqdn:     &fqdn,
			Name:     conv.Pointer("f5test"),
			Provider: conv.Pointer("f5"),
		},
	}

	res := t.c.Domains.PostDomains(domains.PostDomainsParams{Domain: domain})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	domainResponse := domains.PostDomainsCreatedBody{}
	_ = domainResponse.UnmarshalBinary(rr.Body.Bytes())
	return domainResponse.Domain.ID
}

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

func (t *SuiteTest) TestMembersF5AddressValidation() {
	mc := t.c.Members
	defer t.cleanupDomains()
	defer t.cleanupPools()

	// Pool with F5 domain attached
	domainID := t.createF5Domain()
	f5PoolID := t.createPool([]strfmt.UUID{domainID})

	// Pool with Akamai domain attached
	akamaiDomainID := t.createDomain()
	akamaiPoolID := t.createPool([]strfmt.UUID{akamaiDomainID})

	// Pool with no domain
	plainPoolID := t.createPool(nil)

	t.Run("F5 pool: valid IPv4 address is accepted", func() {
		body := members.PostMembersBody{}
		_ = body.UnmarshalBinary([]byte(`{ "member": { "address": "192.0.2.1", "port": 80 } }`))
		body.Member.PoolID = &f5PoolID

		res := mc.PostMembers(members.PostMembersParams{Member: body})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

		// Cleanup member
		resp := members.PostMembersCreatedBody{}
		_ = resp.UnmarshalBinary(rr.Body.Bytes())
		mc.DeleteMembersMemberID(members.DeleteMembersMemberIDParams{MemberID: resp.Member.ID})
	})

	t.Run("F5 pool: valid IPv6 address is accepted", func() {
		body := members.PostMembersBody{}
		_ = body.UnmarshalBinary([]byte(`{ "member": { "address": "::1", "port": 80 } }`))
		body.Member.PoolID = &f5PoolID

		res := mc.PostMembers(members.PostMembersParams{Member: body})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

		resp := members.PostMembersCreatedBody{}
		_ = resp.UnmarshalBinary(rr.Body.Bytes())
		mc.DeleteMembersMemberID(members.DeleteMembersMemberIDParams{MemberID: resp.Member.ID})
	})

	t.Run("F5 pool: hostname address is rejected", func() {
		body := members.PostMembersBody{}
		_ = body.UnmarshalBinary([]byte(`{ "member": { "address": "my.server.example.com", "port": 80 } }`))
		body.Member.PoolID = &f5PoolID

		res := mc.PostMembers(members.PostMembersParams{Member: body})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusBadRequest, rr.Code, rr.Body)
	})

	t.Run("Akamai pool: hostname address is allowed", func() {
		body := members.PostMembersBody{}
		_ = body.UnmarshalBinary([]byte(`{ "member": { "address": "my.server.example.com", "port": 80 } }`))
		body.Member.PoolID = &akamaiPoolID

		res := mc.PostMembers(members.PostMembersParams{Member: body})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

		resp := members.PostMembersCreatedBody{}
		_ = resp.UnmarshalBinary(rr.Body.Bytes())
		mc.DeleteMembersMemberID(members.DeleteMembersMemberIDParams{MemberID: resp.Member.ID})
	})

	t.Run("Pool with no domain: hostname address is allowed", func() {
		body := members.PostMembersBody{}
		_ = body.UnmarshalBinary([]byte(`{ "member": { "address": "my.server.example.com", "port": 80 } }`))
		body.Member.PoolID = &plainPoolID

		res := mc.PostMembers(members.PostMembersParams{Member: body})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

		resp := members.PostMembersCreatedBody{}
		_ = resp.UnmarshalBinary(rr.Body.Bytes())
		mc.DeleteMembersMemberID(members.DeleteMembersMemberIDParams{MemberID: resp.Member.ID})
	})

	t.Run("F5 pool: PUT with hostname address is rejected", func() {
		// Create member with a valid IP first
		createBody := members.PostMembersBody{}
		_ = createBody.UnmarshalBinary([]byte(`{ "member": { "address": "10.0.0.1", "port": 80 } }`))
		createBody.Member.PoolID = &f5PoolID

		createRes := mc.PostMembers(members.PostMembersParams{Member: createBody})
		createRR := httptest.NewRecorder()
		createRes.WriteResponse(createRR, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusCreated, createRR.Code, createRR.Body)

		created := members.PostMembersCreatedBody{}
		_ = created.UnmarshalBinary(createRR.Body.Bytes())
		memberID := created.Member.ID

		// Attempt to update address to a hostname
		updateBody := members.PutMembersMemberIDBody{}
		hostname := "bad.hostname.example.com"
		updateBody.Member = &models.Member{Address: &hostname}

		putRes := mc.PutMembersMemberID(members.PutMembersMemberIDParams{
			MemberID: memberID,
			Member:   updateBody,
		})
		putRR := httptest.NewRecorder()
		putRes.WriteResponse(putRR, runtime.JSONProducer())
		assert.Equal(t.T(), http.StatusBadRequest, putRR.Code, putRR.Body)

		mc.DeleteMembersMemberID(members.DeleteMembersMemberIDParams{MemberID: memberID})
	})
}
