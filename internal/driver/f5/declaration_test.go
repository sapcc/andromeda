// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"testing"

	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
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
		store.On("GetMembers", "dc1").Return([]*rpcmodels.Member{
			{Id: "member1", Address: "200.10.0.1", Port: 80},
		}, nil)
		store.On("GetMembers", "dc2").Return([]*rpcmodels.Member{}, nil)
		b := as3DeclarationBuilder{store: store}
		declaration, err := b.Build()
		assert.Nil(err, "failed to build the declaration")
		expected := as3.ADC{SchemaVersion: "3.22.0"}
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
		expected := as3.ADC{SchemaVersion: "3.22.0"}
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
