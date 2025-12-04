// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/models"
)

var errEntityPendingDeletion = errors.New("this entity has been marked as either PENDING_DELETE or DELETED and therefore must be excluded from the AS3 declaration")

type as3CommonTenantBuilderFunc func(s AndromedaF5Store, datacenters []*rpcmodels.Datacenter, domains []*rpcmodels.Domain) (
	as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error)

type as3DomainTenantBuilderFunc func(f5Config config.F5Config, datacentersByID map[string]*rpcmodels.Datacenter, domain *rpcmodels.Domain) (
	as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error)

func buildAS3Declaration(f5Config config.F5Config, s AndromedaF5Store, ctbFunc as3CommonTenantBuilderFunc, dtbFunc as3DomainTenantBuilderFunc) (
	as3.ADC, *server.ProvisioningStatusRequest, error) {
	adc := as3.NewADC()
	rpcRequest := &server.ProvisioningStatusRequest{}
	rpcUpdates := []*server.ProvisioningStatusRequest_ProvisioningStatus{}
	datacenters, err := s.GetDatacenters()
	if err != nil {
		return adc, rpcRequest, err
	}
	datacentersByID := make(map[string]*rpcmodels.Datacenter, len(datacenters))
	for _, dc := range datacenters {
		datacentersByID[dc.Id] = dc
	}
	// build the /Common key
	domains, err := s.GetDomains()
	if err != nil {
		return adc, rpcRequest, err
	}
	commonTenant, commonTenantRPCUpdates, err := ctbFunc(s, datacenters, domains)
	if err != nil {
		return adc, rpcRequest, err
	}
	rpcUpdates = append(rpcUpdates, commonTenantRPCUpdates...)
	adc.AddTenant("Common", commonTenant)
	// build all /domain_{domainID} keys
	for _, domain := range domains {
		domainTenant, domainTenantRPCUpdates, err := dtbFunc(f5Config, datacentersByID, domain)
		// not a soft error: the declaration cannot be built
		if err != nil && !errors.Is(err, errEntityPendingDeletion) {
			return adc, rpcRequest, err
		}
		// only the active domains may be included in the declaration.
		// domains pending deletion are excluded (and thus deleted by AS3).
		if !errors.Is(err, errEntityPendingDeletion) {
			adc.AddTenant(as3DeclarationGSLBDomainTenantKey(domain.Id), domainTenant)
		}
		rpcUpdates = append(rpcUpdates, domainTenantRPCUpdates...)
	}
	rpcRequest.ProvisioningStatus = rpcUpdates
	return adc, rpcRequest, nil
}

// buildAS3DomainTenant builds one domain tenant ("domain_{domainID}" keys of AS3 declaration).
//
// If the given domain has been either scheduled for deletion or marked as
// deleted, the returned tenant (first return value) should be considered a
// zero value to be discarded. This condition can be checked for by comparing
// the return error (third return value) to errEntityPendingDeletion.
//
// In case the domain provisioning status is set as PENDING_DELETE, the
// returned slice of provisioning status requests will include an entry for
// this domain, which shall be sent over RPC to the Andromeda Server
// immediately after the next successful posting of the AS3 declaration.
func buildAS3DomainTenant(
	f5Config config.F5Config,
	datacentersByID map[string]*rpcmodels.Datacenter,
	domain *rpcmodels.Domain) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
	tenant := as3.Tenant{}
	rpcUpdates := []*server.ProvisioningStatusRequest_ProvisioningStatus{}
	switch domain.ProvisioningStatus {
	case server.ProvisioningStatusRequest_ProvisioningStatus_PENDING_DELETE.String():
		rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
			Id:     domain.Id,
			Model:  server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN,
			Status: server.ProvisioningStatusRequest_ProvisioningStatus_DELETED,
		})
		fallthrough
	case server.ProvisioningStatusRequest_ProvisioningStatus_DELETED.String():
		// by excluding the entity from the AS3 declaration the API will delete it from the F5 device
		return tenant, rpcUpdates, errEntityPendingDeletion
	}
	application := as3.Application{}
	as3PoolReferences := []as3.PointerGSLBPool{}
	for _, p := range domain.Pools {
		switch p.ProvisioningStatus {
		case server.ProvisioningStatusRequest_ProvisioningStatus_PENDING_DELETE.String():
			rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
				Id:     p.Id,
				Model:  server.ProvisioningStatusRequest_ProvisioningStatus_POOL,
				Status: server.ProvisioningStatusRequest_ProvisioningStatus_DELETED,
			})
			fallthrough
		case server.ProvisioningStatusRequest_ProvisioningStatus_DELETED.String():
			// by excluding the entity from the AS3 declaration the API will delete it from the F5 device
			continue
		}
		rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
			Id:     p.Id,
			Model:  server.ProvisioningStatusRequest_ProvisioningStatus_POOL,
			Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
		})
		as3PoolReferences = append(as3PoolReferences, as3.PointerGSLBPool{Use: "pool_" + p.Id})
		as3PoolMembers := []as3.GSLBPoolMember{}
		for _, m := range p.Members {
			switch m.ProvisioningStatus {
			case server.ProvisioningStatusRequest_ProvisioningStatus_PENDING_DELETE.String(),
				server.ProvisioningStatusRequest_ProvisioningStatus_DELETED.String():
				// by excluding the entity from the AS3 declaration the API will delete it from the F5 device.
				// the respective pool member RPC update is covered by buildAS3CommonTenant().
				continue
			}
			if _, exists := datacentersByID[m.DatacenterId]; !exists {
				return tenant, rpcUpdates, fmt.Errorf("invalid datacenter ID for member [datacenter ID = %s, member ID = %s]", m.DatacenterId, m.Id)
			}
			if datacentersByID[m.DatacenterId] == nil {
				return tenant, rpcUpdates, fmt.Errorf("nil datacenter for member [member ID = %s]", m.Id)
			}
			as3PoolMember := as3.GSLBPoolMember{
				Server: as3.PointerGSLBServer{
					Use: "/Common/Shared/" + as3DeclarationGSLBServerKey(m.Address, datacentersByID[m.DatacenterId].Name),
				},
				VirtualServer: as3DeclarationGSLBVirtualServerName(m.Address, m.Port),
			}
			as3PoolMembers = append(as3PoolMembers, as3PoolMember)
		}
		application.SetEntity("pool_"+p.Id, as3.GSLBPool{
			Class:              "GSLB_Pool",
			LBModePreferred:    as3DeclarationPoolMemberLBMode(domain.Mode),
			LBModeAlternate:    "none",
			LBModeFallback:     "none",
			Members:            as3PoolMembers,
			ResourceRecordType: "A",
		})
	}
	rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
		Id:     domain.Id,
		Model:  server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN,
		Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
	})
	as3Domain := as3.GSLBDomain{
		Class:              "GSLB_Domain",
		DomainName:         domain.Fqdn + f5Config.DomainSuffix,
		ResourceRecordType: domain.RecordType,
		PoolLbMode:         as3DeclarationDomainPoolLBMode(),
		Pools:              as3PoolReferences,
	}
	application.SetEntity("wideip", as3Domain)
	tenant.AddApplication("application", application)
	return tenant, rpcUpdates, nil
}

func buildAS3CommonTenant(
	s AndromedaF5Store,
	datacenters []*rpcmodels.Datacenter,
	domains []*rpcmodels.Domain) (as3.Tenant, []*server.ProvisioningStatusRequest_ProvisioningStatus, error) {
	tenant := as3.Tenant{}
	rpcUpdates := []*server.ProvisioningStatusRequest_ProvisioningStatus{}
	application := as3.Application{Template: "shared"}
	// allows referencing monitors by pool ID when iterating over members
	monitorsByPoolID := map[string][]*rpcmodels.Monitor{}
	// add all monitors under /Common/Shared
	for _, domain := range domains {
		for _, pool := range domain.Pools {
			monitorsByPoolID[pool.Id] = pool.Monitors
			for _, monitor := range pool.Monitors {
				switch monitor.ProvisioningStatus {
				case server.ProvisioningStatusRequest_ProvisioningStatus_PENDING_DELETE.String():
					rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
						Id:     monitor.Id,
						Model:  server.ProvisioningStatusRequest_ProvisioningStatus_MONITOR,
						Status: server.ProvisioningStatusRequest_ProvisioningStatus_DELETED,
					})
					fallthrough
				case server.ProvisioningStatusRequest_ProvisioningStatus_DELETED.String():
					// by excluding the entity from the AS3 declaration the API will delete it from the F5 device
					continue
				}
				monitorKey := as3DeclarationGSLBMonitorKey(monitor.Id)
				application.SetEntity(monitorKey, as3.GSLBMonitor{
					Class:        "GSLB_Monitor",
					MonitorType:  as3DeclarationMonitorType(monitor.Type.String()),
					Interval:     monitor.Interval,
					ProbeTimeout: monitor.Timeout,
					Send:         monitor.Send,
					Receive:      monitor.Receive,
				})
				rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
					Id:     monitor.GetId(),
					Model:  server.ProvisioningStatusRequest_ProvisioningStatus_MONITOR,
					Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
				})
			}
		}
	}
	// add all servers under /Common/Shared
	for _, datacenter := range datacenters {
		members, err := s.GetMembers(datacenter.Id)
		if err != nil {
			return tenant, rpcUpdates, err
		}
		for _, member := range members {
			switch member.ProvisioningStatus {
			case server.ProvisioningStatusRequest_ProvisioningStatus_PENDING_DELETE.String():
				rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
					Id:     member.Id,
					Model:  server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER,
					Status: server.ProvisioningStatusRequest_ProvisioningStatus_DELETED,
				})
				fallthrough
			case server.ProvisioningStatusRequest_ProvisioningStatus_DELETED.String():
				// by excluding the entity from the AS3 declaration the API will delete it from the F5 device
				continue
			}
			monitorPointers := []as3.PointerGSLBMonitor{}
			if monitors, ok := monitorsByPoolID[member.PoolId]; ok {
				for _, monitor := range monitors {
					switch monitor.ProvisioningStatus {
					case server.ProvisioningStatusRequest_ProvisioningStatus_PENDING_DELETE.String(),
						server.ProvisioningStatusRequest_ProvisioningStatus_DELETED.String():
						// these monitors are not part of the declaration (see GSLB_Monitor setup above)
						continue
					}
					monitorPointers = append(monitorPointers, as3.PointerGSLBMonitor{
						Use: as3DeclarationGSLBMonitorKey(monitor.Id),
					})
				}
			}
			memberKey := as3DeclarationGSLBServerKey(member.Address, datacenter.Name)
			entity := application.GetEntity(memberKey)
			switch entity {
			case nil:
				application.SetEntity(memberKey, as3.GSLBServer{
					Class:      "GSLB_Server",
					ServerType: "generic-host",
					DataCenter: as3.PointerGSLBDataCenter{BigIP: "/Common/" + datacenter.Name},
					Devices:    []as3.GSLBServerDevice{{Address: member.Address}},
					Monitors:   monitorPointers,
					VirtualServers: []as3.GSLBVirtualServer{{
						Address: member.Address,
						Port:    member.Port,
						Name:    as3DeclarationGSLBVirtualServerName(member.Address, member.Port),
					}},
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
							Name:    as3DeclarationGSLBVirtualServerName(member.Address, member.Port),
						},
					)
					application.SetEntity(memberKey, gslbServer)
				}
			}
			rpcUpdates = append(rpcUpdates, &server.ProvisioningStatusRequest_ProvisioningStatus{
				Id:     member.GetId(),
				Model:  server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER,
				Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
			})
		}
	}
	tenant.AddApplication("Shared", application)
	return tenant, rpcUpdates, nil
}

func as3DeclarationGSLBDomainTenantKey(domainID string) string {
	return "domain_" + domainID
}

func as3DeclarationGSLBPoolKey(poolID string) string {
	return "pool_" + poolID
}

func as3DeclarationGSLBMonitorKey(monitorID string) string {
	return fmt.Sprintf("cc_andromeda_monitor_%s", monitorID)
}

func as3DeclarationGSLBServerKey(memberAddress, datacenterName string) string {
	return fmt.Sprintf("cc_andromeda_srv_%s_%s", memberAddress, datacenterName)
}

func as3DeclarationGSLBVirtualServerName(memberAddress string, memberPort uint32) string {
	return memberAddress + ":" + strconv.FormatUint(uint64(memberPort), 10)
}

// as3DeclarationDomainPoolLBMode refers to valid values for GSLB_Domain.poolLbMode.
func as3DeclarationDomainPoolLBMode() string {
	return "global-availability"
}

// as3DeclarationPoolMemberLBMode refers to valid values for:
//
// - GSLB_Pool.lbModeAlternate
// - GSLB_Pool.lbModeFallback
// - GSLB_Pool.lbModePreferred
//
// Expected behavior for supported values:
//
//   - global-availability: DNS resolution pick is fixed to the first available
//     virtual server in a pool (i.e. GSLB_Pool.Members[0]).
//
//   - round-robin: DNS resolution pick is both circular and sequential among
//     GSLB_Pool.Members[]. Over time each virtual server in a pool is picked
//     an equal amount of times compared to the other pool members.
func as3DeclarationPoolMemberLBMode(memberMode string) string {
	switch memberMode {
	case models.DomainModeROUNDROBIN:
		return "round-robin"
	case models.DomainModeAVAILABILITY:
		return "global-availability"
	default:
		return "round-robin"
	}
}

// as3DeclarationMonitorType maps an Andromeda monitor type string to an F5/AS3 monitor type string.
//
// Andromeda does not support all F5 monitor types.
//
// For the full list of F5 monitor types, see (allowed values of field "monitorType" of class "Monitor"):
// <https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/refguide/schemaref/Monitor.schema.json.html>
func as3DeclarationMonitorType(monitorType string) string {
	switch monitorType {
	case models.MonitorTypeHTTP:
		return "http"
	case models.MonitorTypeHTTPS:
		return "https"
	case models.MonitorTypeICMP:
		return "gateway-icmp"
	case models.MonitorTypeTCP:
		return "tcp"
	case models.MonitorTypeUDP:
		return "udp"
	default:
		return "tcp"
	}
}

// as3Client is a subset of interface bigIPSession
type as3Client interface {
	PostAs3Bigip(as3NewJson, tenantFilter, queryParam string) (error, string, string)
}

func postAS3Declaration(decl as3.ADC, client as3Client, declChecker func(as3.ADC) error) error {
	if err := declChecker(decl); err != nil {
		return err
	}
	jsonDecl, err := json.Marshal(decl)
	if err != nil {
		return err
	}
	log.Debugf("AS3 declaration: %s", string(jsonDecl))
	if err, _, _ := client.PostAs3Bigip(string(jsonDecl), "", ""); err != nil {
		return fmt.Errorf("failed to post AS3 declaration: %w", err)
	}
	return nil
}

var errUnexpectedADCSchemaVersion = errors.New("unexpected AS3 ADC.SchemaVersion")
var errUnexpectedADCUpdateMode = errors.New("unexpected AS3 ADC.UpdateMode")
var errMissingCommonTenant = errors.New("missing required tenant /Common")
var errMissingApplicationCommonShared = errors.New("missing required application /Common/Shared")

func sanityCheckAS3Declaration(decl as3.ADC) error {
	if decl.SchemaVersion != as3.ADCSchemaVersion {
		return errUnexpectedADCSchemaVersion
	}
	if decl.UpdateMode != as3.ADCUpdateMode {
		return errUnexpectedADCUpdateMode
	}
	commonT, err := decl.GetTenant("Common")
	if err != nil {
		return errMissingCommonTenant
	}
	_, err = commonT.GetApplication("Shared")
	if err != nil {
		return errMissingApplicationCommonShared
	}
	return nil
}
