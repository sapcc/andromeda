// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAS3DeclarationBuilder(t *testing.T) {
	assert := assert.New(t)
	t.Run("Two DCs (dc1 and dc2); dc1 has one member; dc2 has no members", func(t *testing.T) {
		// GIVEN
		// * 2 DCs (dc1 and dc2)
		// * dc1 has 1 member; dc2 has no members
		// * dc1 has 1 domain testapp.sap.com
		// EXPECT
		// * 1 GSLB Server in /Common partition
		// * 1 partition for domain testapp.sap.com
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{
			{Id: "dc1", Name: "dc1"},
			{Id: "dc2", Name: "dc2"},
		}, nil)
		// TODO: what domains should be retrieved? filtered by what?
		// there's member.pool_id, so a server is linked to a single pool
		// of course a pool can contain multiple servers
		// there's domain_pool_relation, so a domain may be in multiple pools
		// I think I GetDomains should retrieve all F5 domains by
		// store.On("GetDomains", "")
		store.On("GetMembers", "dc1").Return([]*rpcmodels.Member{
			{Id: "member1", Address: "200.10.0.1", Port: 80},
		}, nil)
		store.On("GetMembers", "dc2").Return([]*rpcmodels.Member{}, nil)
		b := as3DeclarationBuilder{store: store}
		declaration, err := b.Build()
		assert.Nil(err, "failed to build the declaration")
		expected := as3.NewADC()
		func() {
			tenant := as3.Tenant{}
			application := as3.Application{Template: "shared"}
			application.SetEntity("cc_andromeda_srv_200.10.0.1_dc1", as3.GSLBServer{
				Class:          "GSLB_Server",
				ServerType:     "generic-host",
				DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/dc1"},
				Devices:        []as3.GSLBServerDevice{{Address: "200.10.0.1"}},
				Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
				VirtualServers: []as3.GSLBVirtualServer{{Address: "200.10.0.1", Port: 80}},
			})
			tenant.AddApplication("Shared", application)
			expected.AddTenant("Common", tenant)
		}()
		func() {
			//tenant := as3.Tenant{}
			//expected.AddTenant("testapp.sap.com", tenant)
		}()
		assert.Equal(expected, declaration, "declaration built does not match expectation")
	})
	t.Run(`On datacenter "dc1", 2 members share an IP address on different ports`, func(t *testing.T) {
		// GIVEN
		// * 1 DC (dc1) with 3 members.
		// * 2 (member1 and member3) out of 3 members share the same IP address (with different ports: 80 and 9000)
		// EXPECT
		// * 2 (not 3) GSLB servers
		// * member1 and member3 correspond to the same GSLB Server
		// * member2 gets its own GSLB Server
		// * The member1/member3 GSLB Server has two VirtualServers (:80 and :9000).
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{
			{Id: "dc1", Name: "dc1"},
		}, nil)
		store.On("GetMembers", "dc1").Return([]*rpcmodels.Member{
			{Id: "member1", Address: "200.10.0.100", Port: 80},
			{Id: "member2", Address: "200.10.0.2", Port: 80},
			{Id: "member3", Address: "200.10.0.100", Port: 9000},
		}, nil)
		b := as3DeclarationBuilder{store: store}
		declaration, err := b.Build()
		assert.Nil(err, "failed to build the declaration")
		expected := as3.NewADC()
		tenant := as3.Tenant{}
		application := as3.Application{Template: "shared"}
		application.SetEntity("cc_andromeda_srv_200.10.0.100_dc1", as3.GSLBServer{
			Class:      "GSLB_Server",
			ServerType: "generic-host",
			DataCenter: as3.PointerGSLBDataCenter{BigIP: "/Common/dc1"},
			Devices:    []as3.GSLBServerDevice{{Address: "200.10.0.100"}},
			Monitors:   []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
			VirtualServers: []as3.GSLBVirtualServer{
				{Address: "200.10.0.100", Port: 80},
				{Address: "200.10.0.100", Port: 9000},
			},
		})
		application.SetEntity("cc_andromeda_srv_200.10.0.2_dc1", as3.GSLBServer{
			Class:          "GSLB_Server",
			ServerType:     "generic-host",
			DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/dc1"},
			Devices:        []as3.GSLBServerDevice{{Address: "200.10.0.2"}},
			Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
			VirtualServers: []as3.GSLBVirtualServer{{Address: "200.10.0.2", Port: 80}},
		})
		tenant.AddApplication("Shared", application)
		expected.AddTenant("Common", tenant)
		assert.Equal(expected, declaration, "declaration built does not match expectation")
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
	t.Run("it should reject an ad-hoc declaration if the SchemaVersion isn't supported", func(t *testing.T) {
		decl := as3.NewADC()
		decl.SchemaVersion = "1.0.0"
		assert.Same(sanityCheckAS3Declaration(decl), errUnexpectedADCSchemaVersion)
	})
	t.Run("it should reject an ad-hoc declaration if the UpdateMode isn't 'complete'", func(t *testing.T) {
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
