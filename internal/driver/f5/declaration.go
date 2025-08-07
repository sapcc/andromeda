// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/f5devcentral/go-bigip"
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
		log.Debugf("Creating GSLBServer declaration for members of datacenter [id = %q] [name = %s]", datacenter.Id, datacenter.Name)
		members, err := b.store.GetMembers(datacenter.Id)
		if err != nil {
			return tenant, err
		}
		log.Debugf("Found %d members for datacenter [id = %s] [name = %s]", len(members), datacenter.Id, datacenter.Name)
		for _, member := range members {
			memberKey := as3DeclarationGSLBServerKey(member.Address, datacenter.Name)
			entity := application.GetEntity(memberKey)
			switch entity {
			case nil:
				application.SetEntity(memberKey, as3.GSLBServer{
					Class:          "GSLB_Server",
					ServerType:     "generic-host",
					DataCenter:     as3.PointerGSLBDataCenter{BigIP: "/Common/" + datacenter.Name},
					Devices:        []as3.GSLBServerDevice{{Address: member.Address}},
					Monitors:       []as3.PointerGSLBMonitor{{BigIP: "/Common/tcp"}},
					VirtualServers: []as3.GSLBVirtualServer{{Address: member.Address, Port: member.Port}},
				})
			default:
				gslbServer := entity.(as3.GSLBServer)
				mustAddVirtualServer := true
				for _, vs := range gslbServer.VirtualServers {
					if vs.Port == member.Port {
						mustAddVirtualServer = false
						break
					}
				}
				if mustAddVirtualServer {
					gslbServer.VirtualServers = append(gslbServer.VirtualServers,
						as3.GSLBVirtualServer{
							Address: member.Address,
							Port:    member.Port,
						},
					)
					application.SetEntity(memberKey, gslbServer)
				}
			}
		}
	}
	tenant.AddApplication("Shared", application)
	return tenant, nil
}

func as3DeclarationGSLBServerKey(memberAddress, datacenterName string) string {
	return fmt.Sprintf("cc_andromeda_srv_%s_%s", memberAddress, datacenterName)
}

type as3Client interface {
	APICall(options *bigip.APIRequest) ([]byte, error)
}

func postAS3Declaration(decl as3.ADC, client as3Client, declChecker func(as3.ADC) error) error {
	if err := declChecker(decl); err != nil {
		return err
	}
	jsonDecl, err := json.Marshal(decl)
	if err != nil {
		return err
	}
	if _, err := client.APICall(&bigip.APIRequest{
		Method:      "post",
		URL:         "mgmt/shared/appsvcs/declare",
		Body:        string(jsonDecl),
		ContentType: "application/json",
	}); err != nil {
		log.Errorf("failed posting AS3 declaration: %s", err)
		return err
	}
	return nil
}

func sanityCheckAS3Declaration(decl as3.ADC) error {
	commonT, err := decl.GetTenant("Common")
	if err != nil {
		return errors.New("missing required tenant /Common")
	}
	_, err = commonT.GetApplication("Shared")
	if err != nil {
		return errors.New("missing required application /Common/Shared")
	}
	return nil
}
