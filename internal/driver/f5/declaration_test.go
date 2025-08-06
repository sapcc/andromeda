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

func TestBuildsAS3Declaration(t *testing.T) {
	assert := assert.New(t)
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
	application.AddEntity("cc_andromeda_srv_member1_dc1", as3.GSLBServer{
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
}
