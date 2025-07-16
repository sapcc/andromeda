// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/sapcc/andromeda/internal/rpc/server"

	"github.com/f5devcentral/go-bigip"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

func (f5 *F5Agent) GetCommonDeclaration() (*as3.Tenant, error) {
	var virtualServers []as3.GSLBVirtualServer
	res, err := f5.rpc.GetMembers(context.Background(), &server.SearchRequest{
		Provider:      "f5",
		PageNumber:    0,
		ResultPerPage: 1000,
	})
	if err != nil {
		return nil, err
	}

	for _, member := range res.GetResponse() {

		virtualServer := as3.GSLBVirtualServer{
			Name:    "member_" + member.GetId(),
			Address: member.Address,
			Port:    member.Port,
			Enabled: member.AdminStateUp,
		}
		virtualServers = append(virtualServers, virtualServer)
	}

	device := as3.GSLBServerDevice{
		Address: config.Global.F5Config.DNSServerAddress,
	}

	gsblServer := as3.GSLBServer{
		Class: "GSLB_Server",
		DataCenter: as3.PointerGSLBDataCenter{
			Use: "testDataCenter",
		},
		Devices:                  []as3.GSLBServerDevice{device},
		VirtualServers:           virtualServers,
		SnmpProbeEnabled:         true,
		PathProbeEnabled:         true,
		ServiceCheckProbeEnabled: true,
		Monitors: []as3.PointerGSLBMonitor{{
			BigIP: "/Common/gateway_icmp",
		}},
	}

	datacenter := as3.GSLBDatacenter{
		Class:  "GSLB_Data_Center",
		Label:  "Test Datacenter",
		Remark: "Datacenter to test",
	}

	application := as3.Application{
		Template: "shared",
	}
	application.AddEntity("testServer", gsblServer)
	application.AddEntity("testDataCenter", datacenter)

	tenant := as3.Tenant{}
	tenant.AddApplication("Shared", application)

	return &tenant, nil
}

func getPoolLbMode(mode string) string {
	switch mode {
	case "ROUND_ROBIN":
		return "round-robin"
	default:
		return "round-robin"
	}
}

func (f5 *F5Agent) GetTenantDeclaration(domains []*rpcmodels.Domain) (*as3.Tenant, error) {
	application := as3.Application{}
	for _, domain := range domains {
		var as3poolsPtr []as3.PointerGSLBPool

		for _, pool := range domain.Pools {
			as3poolsPtr = append(as3poolsPtr, as3.PointerGSLBPool{Use: "pool_" + pool.GetId()})
		}

		as3domain := as3.GSLBDomain{
			Class:              "GSLB_Domain",
			Label:              domain.GetId(),
			Remark:             "Blub",
			DomainName:         domain.GetFqdn(),
			Aliases:            domain.GetAliases(),
			ResourceRecordType: domain.GetRecordType(),
			PoolLbMode:         getPoolLbMode(domain.Mode),
			Pools:              as3poolsPtr,
		}
		application.AddEntity("domain_"+domain.Id, as3domain)

		for _, pool := range domain.GetPools() {
			var as3PoolMembers []as3.GSLBPoolMember
			for _, member := range pool.GetMembers() {
				as3PoolMember := as3.GSLBPoolMember{
					Server: as3.PointerGSLBServer{
						Use: "/Common/Shared/testServer",
					},
					VirtualServer: "member_" + member.GetId(),
				}
				as3PoolMembers = append(as3PoolMembers, as3PoolMember)
			}

			var as3PoolMonitors []as3.PointerGSLBMonitor
			for _, monitor := range pool.GetMonitors() {
				as3Monitor := as3.GSLBMonitor{
					Class:    "GSLB_Monitor",
					Label:    monitor.GetId(),
					Interval: monitor.GetInterval(),
					Timeout:  monitor.GetTimeout(),
				}
				switch monitor.Type {
				case rpcmodels.Monitor_ICMP:
					as3Monitor.MonitorType = "gateway-icmp"
				case rpcmodels.Monitor_HTTP:
					as3Monitor.MonitorType = "http"
					as3Monitor.Send = monitor.GetSend()
					as3Monitor.Receive = monitor.GetReceive()
				case rpcmodels.Monitor_HTTPS:
					as3Monitor.MonitorType = "https"
					as3Monitor.Send = monitor.GetSend()
					as3Monitor.Receive = monitor.GetReceive()
				case rpcmodels.Monitor_TCP:
					as3Monitor.MonitorType = "tcp"
					as3Monitor.Send = monitor.GetSend()
					as3Monitor.Receive = monitor.GetReceive()
				case rpcmodels.Monitor_UDP:
					as3Monitor.MonitorType = "udp"
					as3Monitor.Send = monitor.GetSend()
					as3Monitor.Receive = monitor.GetReceive()
				}
				application.AddEntity("monitor_"+monitor.GetId(), as3Monitor)
				as3PoolMonitors = append(as3PoolMonitors, as3.PointerGSLBMonitor{Use: "monitor_" + monitor.GetId()})
			}

			as3pool := as3.GSLBPool{
				Class:              "GSLB_Pool",
				Enabled:            pool.GetAdminStateUp(),
				Members:            as3PoolMembers,
				Monitors:           as3PoolMonitors,
				ResourceRecordType: "A",
			}
			application.AddEntity("pool_"+pool.GetId(), as3pool)
		}
	}

	tenant := as3.Tenant{}
	tenant.AddApplication("Application", application)

	return &tenant, nil
}

func ConvertPartitionPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	return GetiControlRestPartitionPath(strings.Split(path, "/"))
}

func GetiControlRestPartitionPath(path []string) string {
	var buffer strings.Builder
	for i, p := range path {
		if i < len(path) {
			buffer.WriteString("~")
		}
		buffer.WriteString(p)
	}
	return buffer.String()
}

func GetForEntity(b *bigip.BigIP, e interface{}, path string) (error, bool) {
	req := &bigip.APIRequest{
		Method:      "get",
		URL:         path,
		ContentType: "application/json",
	}

	resp, err := b.APICall(req)
	if err != nil {
		var reqError bigip.RequestError
		_ = json.Unmarshal(resp, &reqError)
		if reqError.Code == 404 {
			return nil, false
		}
		return err, false
	}

	err = json.Unmarshal(resp, e)
	if err != nil {
		return err, false
	}

	return nil, true
}

func GetSessionHostname(session *bigip.BigIP) (string, error) {
	deviceURL, err := url.Parse(session.Host)
	if err != nil {
		return "", err
	}
	if deviceURL.Hostname() != "" {
		return deviceURL.Hostname(), nil
	}
	return session.Host, nil
}

func FilterDeviceMatchingHostname(devices []bigip.Device, hostname string) (*bigip.Device, error) {
	for _, device := range devices {
		if strings.HasSuffix(hostname, device.Hostname) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("device %s not found", hostname)
}
