// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeclarationSync(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if building AS3 declaration fails", func(t *testing.T) {
		session := new(mockedBigIPSession)
		rpc := new(mockedRPCClient)
		rpc.On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).Return(&server.DatacentersResponse{}, errors.New("RPC GetDatacenters() failed"))
		err := declarationSync(config.F5Config{}, session, rpc)
		assert.ErrorContains(err, "RPC GetDatacenters() failed")
	})

	t.Run("Fails if encoding JSON declaration fails", func(t *testing.T) {
		t.Skip("Not possible to simulate without mocking buildAS3Declaration()")
	})

	t.Run("After successfully building the AS3 declaration, it FAILS...", func(t *testing.T) {
		rpc := new(mockedRPCClient)
		rpc.On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).Return(&server.DatacentersResponse{
			Response: []*rpcmodels.Datacenter{
				{Id: "dc1-uuid", Name: "dc1"},
			},
		}, nil)
		rpc.On("GetDomains", mock.Anything, mock.Anything, mock.Anything).Return(&server.DomainsResponse{
			Response: []*rpcmodels.Domain{
				{Id: "domain1-uuid"},
			},
		}, nil)
		rpc.On("GetMembers", mock.Anything, mock.Anything, mock.Anything).Return(&server.MembersResponse{
			Response: []*rpcmodels.Member{
				{Id: "member1-uuid", Address: "200.10.0.1", Port: 80, DatacenterId: "dc1-uuid"},
			},
		}, nil)

		t.Run("... if it cannot post the AS3 declaration request", func(t *testing.T) {
			session := new(mockedBigIPSession)
			session.On("PostAs3Bigip", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("BigIP PostAs3Bigip() failed"), "", "")
			err := declarationSync(config.F5Config{}, session, rpc)
			assert.ErrorContains(err, "BigIP PostAs3Bigip() failed")
		})

		t.Run("... if it cannot post the update request over RPC", func(t *testing.T) {
			session := new(mockedBigIPSession)
			session.On("PostAs3Bigip", mock.Anything, mock.Anything, mock.Anything).Return(nil, "", "")
			rpc.On("UpdateProvisioningStatus", mock.Anything, mock.Anything, mock.Anything).Return(
				&server.ProvisioningStatusResponse{}, errors.New("RPC UpdateProvisioningStatus() failed"))
			err := declarationSync(config.F5Config{}, session, rpc)
			assert.ErrorContains(err, "RPC UpdateProvisioningStatus() failed")
		})
	})

	t.Run("Succeeds otherwise", func(t *testing.T) {
		session := new(mockedBigIPSession)
		expectedJSONReqPayload := `{"Common":{"Shared":{"cc_andromeda_monitor_mon1-uuid":{"class":"GSLB_Monitor","monitorType":"https","interval":60,"probeTimeout":10,"receive":"bar","send":"foo"},"cc_andromeda_srv_200.10.0.1_dc1":{"class":"GSLB_Server","dataCenter":{"bigip":"/Common/dc1"},"devices":[{"address":"200.10.0.1"}],"virtualServers":[{"address":"200.10.0.1","name":"200.10.0.1:80","port":80}],"serverType":"generic-host"},"cc_andromeda_srv_200.10.0.2_dc2":{"class":"GSLB_Server","dataCenter":{"bigip":"/Common/dc2"},"devices":[{"address":"200.10.0.2"}],"virtualServers":[{"address":"200.10.0.2","name":"200.10.0.2:80","port":80}],"serverType":"generic-host"},"class":"Application","label":"","remark":"","template":"shared"},"class":"Tenant","label":"","remark":""},"class":"ADC","domain_dom1-uuid":{"application":{"class":"Application","label":"","pool_pool1-uuid":{"class":"GSLB_Pool","resourceRecordType":"A","members":[{"server":{"use":"/Common/Shared/cc_andromeda_srv_200.10.0.1_dc1"},"virtualServer":"200.10.0.1:80"}],"lbModePreferred":"global-availability","lbModeAlternate":"none","lbModeFallback":"none"},"pool_pool2-uuid":{"class":"GSLB_Pool","resourceRecordType":"A","members":[{"server":{"use":"/Common/Shared/cc_andromeda_srv_200.10.0.2_dc2"},"virtualServer":"200.10.0.2:80"}],"lbModePreferred":"global-availability","lbModeAlternate":"none","lbModeFallback":"none"},"remark":"","template":"","wideip":{"class":"GSLB_Domain","domainName":"hello-world.local","resourceRecordType":"A","poolLbMode":"global-availability","pools":[{"use":"pool_pool1-uuid"},{"use":"pool_pool2-uuid"}]}},"class":"Tenant","label":"","remark":""},"id":"","schemaVersion":"3.36.0","updateMode":"complete"}`
		session.On("PostAs3Bigip", expectedJSONReqPayload, "", "").Return(nil, "", "")
		expectedRequest := &server.ProvisioningStatusRequest{
			ProvisioningStatus: []*server.ProvisioningStatusRequest_ProvisioningStatus{
				{Id: "mon1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MONITOR, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "member1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "member2-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "pool1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "pool2-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "dom1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			},
		}
		expectedDatacentersSearchRequest := &server.SearchRequest{Provider: "f5", ResultPerPage: 1000}
		rpc := new(mockedRPCClient)
		rpc.On("GetDatacenters", mock.Anything, expectedDatacentersSearchRequest, mock.Anything).Return(&server.DatacentersResponse{
			Response: []*rpcmodels.Datacenter{
				{Id: "dc1-uuid", Name: "dc1"},
				{Id: "dc2-uuid", Name: "dc2"},
			},
		}, nil)
		expectedMembersSearchRequests := []*server.SearchRequest{
			&server.SearchRequest{DatacenterId: "dc1-uuid", ResultPerPage: 1000},
			&server.SearchRequest{DatacenterId: "dc2-uuid", ResultPerPage: 1000},
		}
		rpc.On("GetMembers", mock.Anything, expectedMembersSearchRequests[0], mock.Anything).Return(&server.MembersResponse{
			Response: []*rpcmodels.Member{
				{Id: "member1-uuid", Address: "200.10.0.1", Port: 80, DatacenterId: "dc1-uuid"},
			},
		}, nil)
		rpc.On("GetMembers", mock.Anything, expectedMembersSearchRequests[1], mock.Anything).Return(&server.MembersResponse{
			Response: []*rpcmodels.Member{
				{Id: "member2-uuid", Address: "200.10.0.2", Port: 80, DatacenterId: "dc2-uuid"},
			},
		}, nil)
		expectedDomainsSearchRequest := &server.SearchRequest{Provider: "f5", ResultPerPage: 1000, FullyPopulated: true}
		rpc.On("GetDomains", mock.Anything, expectedDomainsSearchRequest, mock.Anything).Return(&server.DomainsResponse{
			Response: []*rpcmodels.Domain{
				{
					Id:         "dom1-uuid",
					Fqdn:       "hello-world",
					Mode:       models.DomainModeAVAILABILITY,
					RecordType: "A",
					Pools: []*rpcmodels.Pool{
						{
							Id: "pool1-uuid",
							Members: []*rpcmodels.Member{
								{Id: "member1-uuid", Address: "200.10.0.1", Port: 80, DatacenterId: "dc1-uuid"},
							},
							Monitors: []*rpcmodels.Monitor{
								{
									Id:       "mon1-uuid",
									PoolId:   "pool1-uuid",
									Type:     rpcmodels.Monitor_HTTPS,
									Send:     "foo",
									Receive:  "bar",
									Interval: 60,
									Timeout:  10,
								},
							},
						},
						{
							Id: "pool2-uuid",
							Members: []*rpcmodels.Member{
								{Id: "member2-uuid", Address: "200.10.0.2", Port: 80, DatacenterId: "dc2-uuid"},
							},
						},
					},
				},
			},
		}, nil)
		rpc.On("UpdateProvisioningStatus", mock.Anything, expectedRequest, mock.Anything).Return(&server.ProvisioningStatusResponse{}, nil)
		err := declarationSync(config.F5Config{DomainSuffix: ".local"}, session, rpc)
		assert.Nil(err)
		session.AssertCalled(t, "PostAs3Bigip", expectedJSONReqPayload, "", "")
		rpc.AssertCalled(t, "GetDatacenters", mock.Anything, expectedDatacentersSearchRequest, mock.Anything)
		rpc.AssertCalled(t, "GetMembers", mock.Anything, expectedMembersSearchRequests[0], mock.Anything)
		rpc.AssertCalled(t, "GetMembers", mock.Anything, expectedMembersSearchRequests[1], mock.Anything)
		rpc.AssertCalled(t, "GetDomains", mock.Anything, expectedDomainsSearchRequest, mock.Anything)
		rpc.AssertCalled(t, "UpdateProvisioningStatus", mock.Anything, expectedRequest, mock.Anything)
	})
}

func TestStatusSync(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if it cannot build the update request", func(t *testing.T) {
		session := new(mockedBigIPSession)
		rpc := new(mockedRPCClient)
		rpc.On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).Return(&server.DatacentersResponse{}, errors.New("RPC failed for datacenters"))
		err := statusSync(config.F5Config{}, session, rpc)
		assert.ErrorContains(err, "RPC failed for datacenters")
		rpc.AssertNumberOfCalls(t, "GetDatacenters", 1)
		rpc.AssertNotCalled(t, "GetDomains")
		rpc.AssertNotCalled(t, "UpdateMemberStatus")
		session.AssertNotCalled(t, "APICall", mock.Anything)
	})

	t.Run("Fails if it cannot post the update request over RPC", func(t *testing.T) {
		session := new(mockedBigIPSession)
		rpc := new(mockedRPCClient)
		rpc.On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).Return(&server.DatacentersResponse{Response: []*rpcmodels.Datacenter{
			{Id: "dc1-uuid", Name: "dc1-name"},
			{Id: "dc2-uuid", Name: "dc2-name"},
		}}, nil)
		rpc.On("GetDomains", mock.Anything, mock.Anything, mock.Anything).Return(&server.DomainsResponse{Response: []*rpcmodels.Domain{
			{Id: "dom1-uuid", Pools: []*rpcmodels.Pool{}},
		}}, nil)
		rpc.On("UpdateMemberStatus", mock.Anything, mock.Anything, mock.Anything).Return(&server.MemberStatusResponse{}, errors.New("RPC failed for UpdateMemberStatus"))
		err := statusSync(config.F5Config{}, session, rpc)
		assert.ErrorContains(err, "RPC failed for UpdateMemberStatus")
		rpc.AssertNumberOfCalls(t, "GetDatacenters", 1)
		rpc.AssertNumberOfCalls(t, "GetDomains", 1)
		rpc.AssertNumberOfCalls(t, "UpdateMemberStatus", 1)
		session.AssertNotCalled(t, "APICall")
	})

	t.Run("Succeeds if it can post the update request over RPC", func(t *testing.T) {
		expectedURLPath := "gtm/pool/a/~domain_dom1-uuid~application~pool_pool1-uuid/members/~Common~cc_andromeda_srv_10.10.0.11_dc1-name:10.10.0.11:80/stats"
		session := new(mockedBigIPSession)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPath}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"status.availabilityState": {"description": "available"}}}}}}`), nil)
		rpc := new(mockedRPCClient)
		rpc.On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).Return(&server.DatacentersResponse{Response: []*rpcmodels.Datacenter{
			{Id: "dc1-uuid", Name: "dc1-name"},
			{Id: "dc2-uuid", Name: "dc2-name"}}}, nil)
		rpc.On("GetDomains", mock.Anything, mock.Anything, mock.Anything).Return(
			&server.DomainsResponse{
				Response: []*rpcmodels.Domain{
					{Id: "dom1-uuid",
						Pools: []*rpcmodels.Pool{
							{Id: "pool1-uuid",
								Members: []*rpcmodels.Member{
									{
										Id:           "member1-uuid",
										Address:      "10.10.0.11",
										Port:         80,
										DatacenterId: "dc1-uuid",
									}}}}}}}, nil)
		expectedReq := &server.MemberStatusRequest{
			MemberStatus: []*server.MemberStatusRequest_MemberStatus{
				{
					Id:     "member1-uuid",
					Status: server.MemberStatusRequest_MemberStatus_ONLINE,
				}}}
		rpc.On("UpdateMemberStatus", mock.Anything, mock.Anything, mock.Anything).Return(&server.MemberStatusResponse{}, nil)
		err := statusSync(config.F5Config{}, session, rpc)
		assert.Nil(err)
		session.AssertNumberOfCalls(t, "APICall", 1)
		rpc.AssertNumberOfCalls(t, "GetDatacenters", 1)
		rpc.AssertNumberOfCalls(t, "GetDomains", 1)
		rpc.AssertCalled(t, "UpdateMemberStatus", mock.Anything, expectedReq, mock.Anything)
	})
}

func TestMetricsSync(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if metrics collection fails", func(t *testing.T) {
		session := new(mockedBigIPSession)
		rpc := new(mockedRPCClient)
		rpc.
			On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).
			Return(
				&server.DatacentersResponse{},
				errors.New("RPC GetDatacenters() failed"),
			)
		err := metricsSync(config.F5Config{}, session, rpc)
		assert.ErrorContains(err, "RPC GetDatacenters() failed")
	})

	t.Run("Succeeds otherwise", func(t *testing.T) {
		session := new(mockedBigIPSession)
		rpc := new(mockedRPCClient)
		rpc.
			On("GetDatacenters", mock.Anything, mock.Anything, mock.Anything).
			Return(
				&server.DatacentersResponse{
					Response: []*rpcmodels.Datacenter{
						{Id: "dc1-uuid", Name: "dc1-name"},
					},
				},
				nil,
			)
		rpc.
			On("GetDomains", mock.Anything, mock.Anything, mock.Anything).
			Return(
				&server.DomainsResponse{
					Response: []*rpcmodels.Domain{
						{Id: "domain1-uuid"},
					},
				},
				nil,
			)
		err := metricsSync(config.F5Config{}, session, rpc)
		assert.Nil(err)
	})
}
