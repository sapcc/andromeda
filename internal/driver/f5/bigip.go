// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/f5devcentral/go-bigip"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

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
		application.SetEntity("domain_"+domain.Id, as3domain)

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
				application.SetEntity("monitor_"+monitor.GetId(), as3Monitor)
				as3PoolMonitors = append(as3PoolMonitors, as3.PointerGSLBMonitor{Use: "monitor_" + monitor.GetId()})
			}

			as3pool := as3.GSLBPool{
				Class:              "GSLB_Pool",
				Enabled:            pool.GetAdminStateUp(),
				Members:            as3PoolMembers,
				Monitors:           as3PoolMonitors,
				ResourceRecordType: "A",
			}
			application.SetEntity("pool_"+pool.GetId(), as3pool)
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

func GetForEntity(s bigIPSession, e interface{}, path string) (error, bool) {
	req := &bigip.APIRequest{
		Method:      "get",
		URL:         path,
		ContentType: "application/json",
	}

	resp, err := s.APICall(req)
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

// TODO what's above this line is pending refactoring

type bigIPSession interface {
	APICall(options *bigip.APIRequest) ([]byte, error)
	GetDevices() ([]bigip.Device, error)
	GetHost() string
}

type bigIP struct {
	*bigip.BigIP
}

func newBigIPSession(b *bigip.BigIP) *bigIP {
	return &bigIP{BigIP: b}
}

func (b *bigIP) GetHost() string {
	return b.BigIP.Host
}

type activeDeviceMatcher func(bigIPSession) (*bigip.Device, error)
type deviceSessionFactory func(url string) (bigIPSession, error)

func getActiveDeviceSession(
	conf config.F5Config,
	factory deviceSessionFactory,
	matcher activeDeviceMatcher) (bigIPSession, *bigip.Device, error) {
	var s bigIPSession
	var d *bigip.Device
	for _, url := range conf.Devices {
		session, err := factory(url)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create session: %s", err)
		}
		device, err := matchActiveDevice(session)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get active device: %s", err)
		}
		if device != nil {
			s = session
			d = device
			break
		}
	}
	return s, d, nil
}

func matchActiveDevice(session bigIPSession) (*bigip.Device, error) {
	hostname, err := getSessionHostnameFromURL(session.GetHost())
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname from F5 device session: %s", err)
	}
	devices, err := session.GetDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices from F5 device session: %s", err)
	}
	device, err := filterDeviceMatchingHostnameSuffix(devices, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to filter F5 device matching session hostname: %v", err)
	}
	if device.FailoverState != "active" {
		return nil, nil
	}
	return device, nil
}

func getSessionHostnameFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Hostname() != "" {
		return parsedURL.Hostname(), nil
	}
	return rawURL, nil
}

func filterDeviceMatchingHostnameSuffix(devices []bigip.Device, hostname string) (*bigip.Device, error) {
	for _, device := range devices {
		if strings.HasSuffix(hostname, device.Hostname) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("device %s not found", hostname)
}

func getBigIPSession(rawURL string) (bigIPSession, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	user := parsedURL.User.Username()
	if user == "" {
		var ok bool
		user, ok = os.LookupEnv("BIGIP_USER")
		if !ok {
			return nil, fmt.Errorf("BIGIP_USER required for host '%s'", parsedURL.Hostname())
		}
	}
	// check for password
	password, ok := parsedURL.User.Password()
	if !ok {
		password, ok = os.LookupEnv("BIGIP_PASSWORD")
		if !ok {
			return nil, fmt.Errorf("BIGIP_PASSWORD required for host '%s'", parsedURL.Hostname())
		}
	}
	// todo: make configurable
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return newBigIPSession(bigip.NewSession(&bigip.Config{
		Address:           parsedURL.Hostname(),
		Username:          user,
		Password:          password,
		LoginReference:    "tmos",
		CertVerifyDisable: !config.Global.F5Config.ValidateCert,
		ConfigOptions: &bigip.ConfigOptions{
			APICallTimeout: 60 * time.Second,
			TokenTimeout:   1200 * time.Second,
			APICallRetries: int(config.Global.F5Config.MaxRetries),
		},
	})), nil
}
