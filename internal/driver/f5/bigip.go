// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/apex/log"
	"github.com/f5devcentral/go-bigip"

	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

func (f5 *F5Agent) GetCommonTenantDeclaration(datacenters []*rpcmodels.Datacenter) (*as3.Tenant, error) {
	tenant := as3.Tenant{}
	application := as3.Application{Template: "shared"}
	for _, datacenter := range datacenters {
		log.Infof("Creating GSLBServer declaration for members of datacenter [id = %q] [name = %s]", datacenter.Id, datacenter.Name)
		members, err := f5.getMembers(datacenter.Id)
		if err != nil {
			return nil, err
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

func (f5 *F5Agent) GetProjectTenantDeclaration(domains []*rpcmodels.Domain) (*as3.Tenant, error) {
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
