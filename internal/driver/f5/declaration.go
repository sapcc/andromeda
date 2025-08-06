// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"fmt"

	"github.com/apex/log"
	"github.com/sapcc/andromeda/internal/driver/f5/as3"
)

type AS3DeclarationBuilder interface {
	Build() (as3.ADC, error)
}

type as3DeclarationBuilder struct {
	store AndromedaF5Store
}

func NewAS3DeclarationBuilder(s AndromedaF5Store) AS3DeclarationBuilder {
	return &as3DeclarationBuilder{store: s}
}

func (b *as3DeclarationBuilder) Build() (as3.ADC, error) {
	adc := as3.ADC{SchemaVersion: "3.22.0"}
	common, err := b.getCommonTenant()
	if err != nil {
		return adc, err
	}
	adc.AddTenant("Common", common)
	return adc, nil
}

func (b *as3DeclarationBuilder) getCommonTenant() (as3.Tenant, error) {
	tenant := as3.Tenant{}
	application := as3.Application{Template: "shared"}
	datacenters, err := b.store.GetDatacenters()
	if err != nil {
		return tenant, err
	}
	for _, datacenter := range datacenters {
		log.Infof("Creating GSLBServer declaration for members of datacenter [id = %q] [name = %s]", datacenter.Id, datacenter.Name)
		members, err := b.store.GetMembers(datacenter.Id)
		if err != nil {
			return tenant, err
		}
		log.Infof("Found %d members for datacenter [id = %s] [name = %s]", len(members), datacenter.Id, datacenter.Name)
		for _, member := range members {
			// TODO: members that share `member.Address` must be the same GSLBServer entity
			application.AddEntity(fmt.Sprintf("cc_andromeda_srv_%s_%s", member.Id, datacenter.Name), as3.GSLBServer{
				Class:          "GSLB_Server",
				ServerType:     "generic-host",
				DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/" + datacenter.Name},
				Devices:        []as3.GSLBServerDevice{{Address: member.Address}},
				Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
				VirtualServers: []as3.GSLBVirtualServer{{Address: member.Address, Port: member.Port}},
			})
		}
	}
	tenant.AddApplication("Shared", application)
	return tenant, nil
}
