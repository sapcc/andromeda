// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuildAS3Declaration(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if RPC GetDatacenters() fails", func(t *testing.T) {
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{}, errors.New("RPC GetDatacenters() failed"))
		_, _, err := buildAS3Declaration(config.F5Config{}, store, buildAS3CommonTenant, buildAS3DomainTenant)
		assert.ErrorContains(err, "RPC GetDatacenters() failed")
	})

	t.Run("Fails if AS3 Common tenant builder function fails", func(t *testing.T) {
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{{Id: "dc1-uuid", Name: "dc1"}}, nil)
		as3CommonTenantBuilder := func(s AndromedaF5Store, datacenters []*rpcmodels.Datacenter) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
			return as3.Tenant{}, []*server.ProvisioningStatusRequest_ProvisioningStatus{}, errors.New("ctbFunc failed")
		}
		_, _, err := buildAS3Declaration(config.F5Config{}, store, as3CommonTenantBuilder, buildAS3DomainTenant)
		assert.ErrorContains(err, "ctbFunc failed")
	})

	t.Run("Fails if RPC GetDomains() fails", func(t *testing.T) {
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{{Id: "dc1-uuid", Name: "dc1"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{}, errors.New("RPC GetDomains() failed"))
		as3CommonTenantBuilder := func(s AndromedaF5Store, datacenters []*rpcmodels.Datacenter) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
			return as3.Tenant{}, []*server.ProvisioningStatusRequest_ProvisioningStatus{}, nil
		}
		_, _, err := buildAS3Declaration(config.F5Config{}, store, as3CommonTenantBuilder, buildAS3DomainTenant)
		assert.ErrorContains(err, "RPC GetDomains() failed")
	})

	t.Run("Fails if AS3 Domain tenant builder function fails", func(t *testing.T) {
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{{Id: "dc1-uuid", Name: "dc1"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{{Id: "dom1-uuid"}}, nil)
		as3CommonTenantBuilder := func(s AndromedaF5Store, datacenters []*rpcmodels.Datacenter) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
			return as3.Tenant{}, []*server.ProvisioningStatusRequest_ProvisioningStatus{}, nil
		}
		as3DomainTenantBuilder := func(f5Config config.F5Config, datacentersByID map[string]*rpcmodels.Datacenter, domain *rpcmodels.Domain) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
			return as3.Tenant{}, []*server.ProvisioningStatusRequest_ProvisioningStatus{}, errors.New("dtbFunc failed")
		}
		_, _, err := buildAS3Declaration(config.F5Config{}, store, as3CommonTenantBuilder, as3DomainTenantBuilder)
		assert.ErrorContains(err, "dtbFunc failed")
	})

	t.Run("Succeeds by creating the full AS3 declaration", func(t *testing.T) {
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{{Id: "dc1-uuid", Name: "dc1"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{{Id: "dom1-uuid"}}, nil)
		expectedCommonTenant := as3.Tenant{}
		expectedDomainTenant := as3.Tenant{}
		as3CommonTenantBuilder := func(s AndromedaF5Store, datacenters []*rpcmodels.Datacenter) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
			return expectedCommonTenant, []*server.ProvisioningStatusRequest_ProvisioningStatus{
				{Id: "member1", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			}, nil
		}
		as3DomainTenantBuilder := func(f5Config config.F5Config, datacentersByID map[string]*rpcmodels.Datacenter, domain *rpcmodels.Domain) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
			return expectedDomainTenant, []*server.ProvisioningStatusRequest_ProvisioningStatus{
				{Id: "pool1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "pool2-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "dom1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			}, nil
		}

		declaration, req, err := buildAS3Declaration(config.F5Config{}, store, as3CommonTenantBuilder, as3DomainTenantBuilder)
		assert.Nil(err)

		t.Run("Adds the common tenant to the AS3 declaration", func(t *testing.T) {
			tenant, err := declaration.GetTenant("Common")
			assert.Nil(err)
			assert.Equal(expectedCommonTenant, tenant)
		})

		t.Run("Adds the domain tenant to the AS3 declaration", func(t *testing.T) {
			tenant, err := declaration.GetTenant("domain_dom1-uuid")
			assert.Nil(err)
			assert.Equal(expectedDomainTenant, tenant)
		})

		t.Run("Combines the RPC update requests returned by both the common and domain tenant builders", func(t *testing.T) {
			expectedReq := &server.ProvisioningStatusRequest{
				ProvisioningStatus: []*server.ProvisioningStatusRequest_ProvisioningStatus{
					{Id: "member1", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
					{Id: "pool1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
					{Id: "pool2-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
					{Id: "dom1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				},
			}
			assert.Equal(expectedReq, req)
		})
	})
}

func TestBuildAS3DomainTenant(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if datacentersByID map is missing a datacenter referenced by a member", func(t *testing.T) {
		datacentersByID := map[string]*rpcmodels.Datacenter{"dc50-uuid": {Id: "dc50-uuid", Name: "dc50"}}
		domain := &rpcmodels.Domain{
			Id:         "dom1-uuid",
			Fqdn:       "test1.internal",
			Mode:       "global",
			RecordType: "A",
			Pools: []*rpcmodels.Pool{
				{Id: "pool1-uuid", Members: []*rpcmodels.Member{{Id: "member1", Address: "200.10.0.1", Port: 80, DatacenterId: "dc1-uuid"}}},
			},
		}
		_, _, err := buildAS3DomainTenant(config.F5Config{}, datacentersByID, domain)
		assert.ErrorContains(err, "invalid datacenter ID for member")
	})

	t.Run("Fails if datacentersByID map has a nil pointer datacenter referenced by a member", func(t *testing.T) {
		datacentersByID := map[string]*rpcmodels.Datacenter{"dc1-uuid": nil}
		domain := &rpcmodels.Domain{
			Id:         "dom1-uuid",
			Fqdn:       "test1.internal",
			Mode:       "global",
			RecordType: "A",
			Pools: []*rpcmodels.Pool{
				{Id: "pool1-uuid", Members: []*rpcmodels.Member{{Id: "member1", Address: "200.10.0.1", Port: 80, DatacenterId: "dc1-uuid"}}},
			},
		}
		_, _, err := buildAS3DomainTenant(config.F5Config{}, datacentersByID, domain)
		assert.ErrorContains(err, "nil datacenter for member")
	})

	t.Run("Succeeds by creating a direct mapping between each Andromeda domain's pools and members and their F5 counterpart entity", func(t *testing.T) {
		datacentersByID := map[string]*rpcmodels.Datacenter{
			"dc1-uuid": {Id: "dc1-uuid", Name: "dc1"},
		}
		domain := &rpcmodels.Domain{
			Id:         "dom1-uuid",
			Fqdn:       "test1",
			Mode:       "AVAILABILITY",
			RecordType: "A",
			Pools: []*rpcmodels.Pool{
				{
					Id: "pool1-uuid",
					Members: []*rpcmodels.Member{
						{Id: "member1", Address: "200.10.0.1", Port: 80, DatacenterId: "dc1-uuid"},
						{Id: "member2", Address: "200.10.0.2", Port: 80, DatacenterId: "dc1-uuid"},
						{Id: "member3", Address: "200.10.0.3", Port: 80, DatacenterId: "dc1-uuid"},
					},
				},
				{
					Id: "pool2-uuid",
					Members: []*rpcmodels.Member{
						{Id: "member4", Address: "200.10.0.4", Port: 80, DatacenterId: "dc1-uuid"},
						{Id: "member5", Address: "200.10.0.5", Port: 80, DatacenterId: "dc1-uuid"},
					},
				},
			},
		}
		tenant, req, err := buildAS3DomainTenant(config.F5Config{DomainSuffix: ".internal"}, datacentersByID, domain)
		assert.Nil(err)

		t.Run("Builds tenant correctly", func(t *testing.T) {
			expectedTenant := as3.Tenant{}
			application := as3.Application{}
			application.SetEntity("wideip", as3.GSLBDomain{
				Class:              "GSLB_Domain",
				DomainName:         "test1.internal",
				ResourceRecordType: "A",
				PoolLbMode:         "global-availability",
				Pools: []as3.PointerGSLBPool{
					{Use: "pool_pool1-uuid"},
					{Use: "pool_pool2-uuid"},
				},
			})
			application.SetEntity("pool_pool1-uuid", as3.GSLBPool{
				Class:           "GSLB_Pool",
				LBModePreferred: "global-availability",
				LBModeAlternate: "none",
				LBModeFallback:  "none",
				Members: []as3.GSLBPoolMember{
					{Server: as3.PointerGSLBServer{Use: "/Common/Shared/cc_andromeda_srv_200.10.0.1_dc1"}, VirtualServer: "200.10.0.1:80"},
					{Server: as3.PointerGSLBServer{Use: "/Common/Shared/cc_andromeda_srv_200.10.0.2_dc1"}, VirtualServer: "200.10.0.2:80"},
					{Server: as3.PointerGSLBServer{Use: "/Common/Shared/cc_andromeda_srv_200.10.0.3_dc1"}, VirtualServer: "200.10.0.3:80"},
				},
				ResourceRecordType: "A",
			})
			application.SetEntity("pool_pool2-uuid", as3.GSLBPool{
				Class:           "GSLB_Pool",
				LBModePreferred: "global-availability",
				LBModeAlternate: "none",
				LBModeFallback:  "none",
				Members: []as3.GSLBPoolMember{
					{Server: as3.PointerGSLBServer{Use: "/Common/Shared/cc_andromeda_srv_200.10.0.4_dc1"}, VirtualServer: "200.10.0.4:80"},
					{Server: as3.PointerGSLBServer{Use: "/Common/Shared/cc_andromeda_srv_200.10.0.5_dc1"}, VirtualServer: "200.10.0.5:80"},
				},
				ResourceRecordType: "A",
			})
			expectedTenant.AddApplication("application", application)
			assert.Equal(expectedTenant, tenant)
		})

		t.Run("Builds RPC updates request correctly", func(t *testing.T) {
			expectedRPCUpdates := []*server.ProvisioningStatusRequest_ProvisioningStatus{
				{Id: "pool1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "pool2-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_POOL, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
				{Id: "dom1-uuid", Model: server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			}
			assert.Equal(expectedRPCUpdates, req)
		})
	})
}

func TestBuildAS3CommonTenant(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if RPC GetMembers() fails", func(t *testing.T) {
		datacenters := []*rpcmodels.Datacenter{
			{Id: "dc1-uuid", Name: "dc1"},
			{Id: "dc2-uuid", Name: "dc2"},
		}
		s := new(mockedStore)
		s.On("GetMembers", "dc1-uuid").Return([]*rpcmodels.Member{}, errors.New("RPC GetDomains() failed"))
		_, _, err := buildAS3CommonTenant(s, datacenters)
		s.AssertCalled(t, "GetMembers", "dc1-uuid")
		s.AssertNotCalled(t, "GetMembers", "dc2-uuid")
		assert.ErrorContains(err, "RPC GetDomains() failed")
	})

	t.Run("Creates one GSLB Server entity for each datacenter member that does not share an IP address", func(t *testing.T) {
		datacenters := []*rpcmodels.Datacenter{
			{Id: "dc1-uuid", Name: "dc1"},
			{Id: "dc2-uuid", Name: "dc2"},
		}
		store := new(mockedStore)
		store.On("GetMembers", "dc1-uuid").Return([]*rpcmodels.Member{
			{Id: "member1", Address: "200.10.0.1", Port: 80},
			{Id: "member2", Address: "200.10.0.2", Port: 80},   // shares DC IP with member 3
			{Id: "member3", Address: "200.10.0.2", Port: 8080}, // shares DC IP with member 2
		}, nil)
		store.On("GetMembers", "dc2-uuid").Return([]*rpcmodels.Member{
			{Id: "member4", Address: "200.10.0.4", Port: 80},
			{Id: "member5", Address: "200.10.0.5", Port: 80},
		}, nil)
		tenant, req, err := buildAS3CommonTenant(store, datacenters)
		expectedTenant := as3.Tenant{}
		application := as3.Application{Template: "shared"}
		application.SetEntity("cc_andromeda_srv_200.10.0.1_dc1", as3.GSLBServer{
			Class:          "GSLB_Server",
			ServerType:     "generic-host",
			DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/dc1"},
			Devices:        []as3.GSLBServerDevice{{Address: "200.10.0.1"}},
			Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
			VirtualServers: []as3.GSLBVirtualServer{{Address: "200.10.0.1", Port: 80, Name: "200.10.0.1:80"}},
		})
		application.SetEntity("cc_andromeda_srv_200.10.0.2_dc1", as3.GSLBServer{
			Class:      "GSLB_Server",
			ServerType: "generic-host",
			DataCenter: as3.PointerGSLBDataCenter{BigIP: "/Common/dc1"},
			Devices:    []as3.GSLBServerDevice{{Address: "200.10.0.2"}},
			Monitors:   []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
			VirtualServers: []as3.GSLBVirtualServer{
				{Address: "200.10.0.2", Port: 80, Name: "200.10.0.2:80"},
				{Address: "200.10.0.2", Port: 8080, Name: "200.10.0.2:8080"},
			},
		})
		application.SetEntity("cc_andromeda_srv_200.10.0.4_dc2", as3.GSLBServer{
			Class:          "GSLB_Server",
			ServerType:     "generic-host",
			DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/dc2"},
			Devices:        []as3.GSLBServerDevice{{Address: "200.10.0.4"}},
			Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
			VirtualServers: []as3.GSLBVirtualServer{{Address: "200.10.0.4", Port: 80, Name: "200.10.0.4:80"}},
		})
		application.SetEntity("cc_andromeda_srv_200.10.0.5_dc2", as3.GSLBServer{
			Class:          "GSLB_Server",
			ServerType:     "generic-host",
			DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/dc2"},
			Devices:        []as3.GSLBServerDevice{{Address: "200.10.0.5"}},
			Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
			VirtualServers: []as3.GSLBVirtualServer{{Address: "200.10.0.5", Port: 80, Name: "200.10.0.5:80"}},
		})
		expectedTenant.AddApplication("Shared", application)
		expectedRPCUpdates := []*server.ProvisioningStatusRequest_ProvisioningStatus{
			{Id: "member1", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			{Id: "member2", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			{Id: "member3", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			{Id: "member4", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
			{Id: "member5", Model: server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER, Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE},
		}
		assert.Nil(err)
		assert.Equal(expectedTenant, tenant)
		assert.Equal(expectedRPCUpdates, req)
		store.AssertCalled(t, "GetMembers", "dc1-uuid")
		store.AssertCalled(t, "GetMembers", "dc2-uuid")
	})
}

func TestPostAS3Declaration(t *testing.T) {
	assert := assert.New(t)

	t.Run("it should fail (without posting) if checker func fails", func(t *testing.T) {
		decl := as3.ADC{Label: "I should get rejected"}
		client := new(mockedAS3Client)
		declChecker := func(d as3.ADC) error {
			assert.Equal(d.Label, decl.Label)
			return errors.New("nope")
		}
		err := postAS3Declaration(decl, client, declChecker)
		assert.NotNil(err, "it should have failed")
		client.AssertNotCalled(t, "APICall")
	})

	t.Run("if the checker func accepts the declaration then it should post it", func(t *testing.T) {
		decl := as3.NewADC()
		tenant := as3.Tenant{}
		application := as3.Application{Template: "shared"}
		tenant.AddApplication("Shared", application)
		decl.AddTenant("Common", tenant)

		t.Run("it should encode the declaration as JSON before posting", func(t *testing.T) {
			expectedJSONDecl := string(`
			{ "schemaVersion": "3.36.0",
			  "id" :"",
			  "class": "ADC",
			  "updateMode": "complete",
			  "Common": {
			    "class": "Tenant",
			    "label": "",
			    "remark": "",
			    "Shared": { "class": "Application", "label": "", "remark": "", "template": "shared" }
			  }
		        }`)
			client := new(mockedAS3Client)
			client.
				On("APICall",
					mock.MatchedBy(func(req *bigip.APIRequest) bool {
						return assert.JSONEq(expectedJSONDecl, req.Body, "unexpected JSON declaration")
					}),
				).
				Return([]byte(`{"code": 200}`), nil)
			declChecker := func(d as3.ADC) error { return nil }
			err := postAS3Declaration(decl, client, declChecker)
			assert.Nil(err, "it should have succeeded")
			client.AssertCalled(t, "APICall", mock.Anything)
		})

		t.Run("it should fail if posting fails (checking APICall() error return value is enough)", func(t *testing.T) {
			client := new(mockedAS3Client)
			client.On("APICall", mock.Anything).Return([]byte(""), errors.New("it failed, please let the caller now"))
			declChecker := func(d as3.ADC) error { return nil }
			err := postAS3Declaration(decl, client, declChecker)
			assert.NotNil(err, "it should have failed")
			client.AssertCalled(t, "APICall", mock.Anything)
		})

		t.Run("it should succeed if posting succeeds (checking APICall() error return value is enough)", func(t *testing.T) {
			client := new(mockedAS3Client)
			client.On("APICall", mock.Anything).Return([]byte(""), nil)
			declChecker := func(d as3.ADC) error { return nil }
			err := postAS3Declaration(decl, client, declChecker)
			assert.Nil(err, "it should have succeeded")
			client.AssertCalled(t, "APICall", mock.Anything)
		})
	})
}

func TestSanityCheckAS3Declaration(t *testing.T) {
	assert := assert.New(t)

	t.Run("it should reject a standard declaration if the SchemaVersion is tainted", func(t *testing.T) {
		decl := as3.NewADC()
		decl.SchemaVersion = "1.0.0"
		assert.Same(sanityCheckAS3Declaration(decl), errUnexpectedADCSchemaVersion)
	})

	t.Run("it should reject a standard declaration if the UpdateMode is tainted", func(t *testing.T) {
		decl := as3.NewADC()
		decl.UpdateMode = "selective"
		assert.Same(sanityCheckAS3Declaration(decl), errUnexpectedADCUpdateMode)
	})

	t.Run("it should reject a standard declaration if the /Common tenant is missing", func(t *testing.T) {
		decl := as3.NewADC()
		assert.Same(sanityCheckAS3Declaration(decl), errMissingCommonTenant)
	})

	t.Run("it should reject a standard declaration if the /Common/Shared application is missing", func(t *testing.T) {
		decl := as3.NewADC()
		tenant := as3.Tenant{}
		decl.AddTenant("Common", tenant)
		assert.Same(sanityCheckAS3Declaration(decl), errMissingApplicationCommonShared)
	})

	t.Run("it should accept a standard declaration if the /Common/Shared application is present", func(t *testing.T) {
		decl := as3.NewADC()
		tenant := as3.Tenant{}
		application := as3.Application{Template: "shared"}
		tenant.AddApplication("Shared", application)
		decl.AddTenant("Common", tenant)
		assert.Nil(sanityCheckAS3Declaration(decl), "declaration should have been accepted")
	})
}
