// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuildMemberStatusUpdateRequest(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if it cannot fetch datacenters over RPC", func(t *testing.T) {
		session := new(mockedBigIPSession)
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{}, errors.New("RPC failed for datacenters"))
		_, err := buildMemberStatusUpdateRequest(session, store)
		assert.ErrorContains(err, "RPC failed for datacenters")
		store.AssertCalled(t, "GetDatacenters")
		store.AssertNotCalled(t, "GetDomains")
		session.AssertNotCalled(t, "APIRequest", mock.Anything)
	})

	t.Run("Fails if it cannot fetch domains over RPC", func(t *testing.T) {
		session := new(mockedBigIPSession)
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{{Id: "dc1-uuid", Name: "dc1-name"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{}, errors.New("RPC failed for domains"))
		_, err := buildMemberStatusUpdateRequest(session, store)
		assert.ErrorContains(err, "RPC failed for domains")
		store.AssertCalled(t, "GetDatacenters")
		store.AssertCalled(t, "GetDomains")
		session.AssertNotCalled(t, "APIRequest", mock.Anything)
	})

	t.Run("Succeeds if the domain/pool/member entries returned by RPC are correct and F5 API works", func(t *testing.T) {
		expectedURLPaths := []string{
			"gtm/pool/a/~domain_tenant_dom1-uuid~application~pool_pool1-uuid/members/~Common~cc_andromeda_srv_10.10.0.11_dc1-name:10.10.0.11:80/stats",
			"gtm/pool/a/~domain_tenant_dom1-uuid~application~pool_pool1-uuid/members/~Common~cc_andromeda_srv_10.10.0.12_dc2-name:10.10.0.12:80/stats",
			"gtm/pool/a/~domain_tenant_dom1-uuid~application~pool_pool1-uuid/members/~Common~cc_andromeda_srv_10.10.0.13_dc2-name:10.10.0.13:80/stats",
			"gtm/pool/a/~domain_tenant_dom1-uuid~application~pool_pool1-uuid/members/~Common~cc_andromeda_srv_10.10.0.14_dc2-name:10.10.0.14:80/stats"}
		session := new(mockedBigIPSession)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[0]}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "offline"}}}}}}`), nil)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[1]}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "available"}}}}}}`), nil)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[2]}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "unknown"}}}}}}`), nil)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[3]}).
			Return([]byte(`{"code": 404}`), &errBigIPEntityNotFound{})
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{
			{Id: "dc1-uuid", Name: "dc1-name"},
			{Id: "dc2-uuid", Name: "dc2-name"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{{Id: "dom1-uuid", Pools: []*rpcmodels.Pool{
			{Id: "pool1-uuid", Members: []*rpcmodels.Member{
				{Id: "member1-uuid", Address: "10.10.0.11", Port: 80, DatacenterId: "dc1-uuid"},
				{Id: "member2-uuid", Address: "10.10.0.12", Port: 80, DatacenterId: "dc2-uuid"},
				{Id: "member3-uuid", Address: "10.10.0.13", Port: 80, DatacenterId: "dc2-uuid"},
				{Id: "member4-uuid", Address: "10.10.0.14", Port: 80, DatacenterId: "dc2-uuid"}}}}}}, nil)
		req, err := buildMemberStatusUpdateRequest(session, store)
		assert.Nil(err)
		expectedReq := &server.MemberStatusRequest{MemberStatus: []*server.MemberStatusRequest_MemberStatus{
			{Id: "member1-uuid", Status: server.MemberStatusRequest_MemberStatus_OFFLINE},
			{Id: "member2-uuid", Status: server.MemberStatusRequest_MemberStatus_ONLINE},
			{Id: "member3-uuid", Status: server.MemberStatusRequest_MemberStatus_UNKNOWN},
			{Id: "member4-uuid", Status: server.MemberStatusRequest_MemberStatus_UNKNOWN}}}
		assert.Equal(expectedReq, req)
	})
}

func TestFetchPoolTypeAMemberAvailability(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if F5 API request fails", func(t *testing.T) {
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(""), errors.New("API Request failed"))
		_, err := fetchPoolTypeAMemberAvailability(session, "/some/path")
		assert.ErrorContains(err, "API Request failed")
	})

	t.Run("Fails if JSON response payload contains more than one key in `.entries`", func(t *testing.T) {
		jsonDoc := `{
			"kind": "tm:gtm:pool:a:members:membersstats",
			"entries": {
			  "theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "offline"}}}},
			  "badKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "offline"}}}}
			}
		}`
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(jsonDoc), nil)
		_, err := fetchPoolTypeAMemberAvailability(session, "/some/path")
		assert.ErrorContains(err, "expected exactly 1 key")
	})

	t.Run("Fails if JSON response payload contains an invalid value at `.entries.theKey.nestedStats`", func(t *testing.T) {
		jsonDoc := `{
			"kind": "tm:gtm:pool:a:members:membersstats",
			"entries": {"theKey": "badNestedStats"}
		}`
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(jsonDoc), nil)
		_, err := fetchPoolTypeAMemberAvailability(session, "/some/path")
		assert.ErrorContains(err, "could not decode nested member stats")
	})

	t.Run("Succeds if JSON payload is decoded correctly all the way to the availability state description", func(t *testing.T) {
		jsonDoc := `{
			"kind": "tm:gtm:pool:a:members:membersstats",
			"entries": {"theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "offline"}}}}}
		}`
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(jsonDoc), nil)
		status, err := fetchPoolTypeAMemberAvailability(session, "/some/path")
		session.AssertCalled(t, "APICall", &bigip.APIRequest{Method: "get", URL: "/some/path", ContentType: "application/json"})
		assert.Nil(err)
		assert.Equal("offline", status)
	})
}

func TestFetchPoolTypeAMemberStats(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if F5 API request fails", func(t *testing.T) {
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(""), errors.New("API Request failed"))
		_, err := fetchPoolTypeAMemberStats(session, "/some/path")
		assert.ErrorContains(err, "API Request failed")
	})

	t.Run("Fails with special error if F5 API request soft fails with 404", func(t *testing.T) {
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(`{"code": 404}`), errors.New("404"))
		_, err := fetchPoolTypeAMemberStats(session, "/some/path")
		assert.ErrorContains(err, "entity not found")
	})

	t.Run("Fails if F5 API request succeeds but response JSON payload cannot be decoded", func(t *testing.T) {
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(`{`), nil)
		_, err := fetchPoolTypeAMemberStats(session, "/some/path")
		assert.ErrorContains(err, "unexpected end of JSON input")
	})

	t.Run("Succeeds if F5 API request succeeds and JSON payload can be decoded", func(t *testing.T) {
		jsonDoc := `{
			"kind": "tm:gtm:pool:a:members:membersstats",
			"entries": {"theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "offline"}}}}}
		}`
		session := new(mockedBigIPSession)
		session.On("APICall", mock.Anything).Return([]byte(jsonDoc), nil)
		mcs, err := fetchPoolTypeAMemberStats(session, "/some/path")
		session.AssertCalled(t, "APICall", &bigip.APIRequest{Method: "get", URL: "/some/path", ContentType: "application/json"})
		assert.Nil(err)
		assert.Equal("tm:gtm:pool:a:members:membersstats", mcs.Kind)
	})
}
