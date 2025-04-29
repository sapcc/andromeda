// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package as3

import "encoding/json"

// AS3
type AS3 struct {
	Persist     bool        `json:"persist"`
	Class       string      `json:"class"`
	Action      string      `json:"action,omitempty"`
	Declaration interface{} `json:"declaration,omitempty"`
}

// ADC
type ADC struct {
	Label         string `json:"label,omitempty"`
	Remark        string `json:"remark,omitempty"`
	SchemaVersion string `json:"schemaVersion"`
	Id            string `json:"id,omitempty"`

	tenants map[string]Tenant
}

func (a ADC) MarshalJSON() ([]byte, error) {
	adc := make(map[string]interface{})
	adc["class"] = "ADC"
	adc["schemaVersion"] = a.SchemaVersion
	adc["id"] = a.Id

	for name, tenant := range a.tenants {
		adc[name] = tenant
	}
	return json.Marshal(adc)
}

func (a *ADC) AddTenant(name string, tenant Tenant) {
	if a.tenants == nil {
		a.tenants = make(map[string]Tenant, 1)
	}
	a.tenants[name] = tenant
}

// Tenants
type Tenant struct {
	Label  string `json:"label,omitempty"`
	Remark string `json:"remark,omitempty"`

	applications Applications
}

func (t Tenant) MarshalJSON() ([]byte, error) {
	tenant := make(map[string]interface{})
	tenant["class"] = "Tenant"
	tenant["label"] = t.Label
	tenant["remark"] = t.Remark

	for name, application := range t.applications {
		tenant[name] = application
	}
	return json.Marshal(tenant)
}

func (t *Tenant) AddApplication(name string, application Application) {
	if t.applications == nil {
		t.applications = make(Applications, 1)
	}
	t.applications[name] = application
}

// Applications
type Applications map[string]Application
type Application struct {
	Label    string `json:"label,omitempty"`
	Remark   string `json:"remark,omitempty"`
	Template string `json:"template,omitempty"`

	// Application entities
	entities map[string]interface{}
}

func (a Application) MarshalJSON() ([]byte, error) {
	application := make(map[string]interface{})
	application["class"] = "Application"
	application["label"] = a.Label
	application["remark"] = a.Remark
	application["template"] = a.Template

	for name, entity := range a.entities {
		application[name] = entity
	}
	return json.Marshal(application)
}

func (a *Application) AddEntity(name string, entity interface{}) {
	if a.entities == nil {
		a.entities = make(map[string]interface{}, 1)
	}
	a.entities[name] = entity
}

// GSLB Entities
type GSLBDatacenter struct {
	Class  string `json:"class"`
	Label  string `json:"label,omitempty"`
	Remark string `json:"remark,omitempty"`
}

type GSLBDomain struct {
	Class              string            `json:"class"`
	Label              string            `json:"label,omitempty"`
	Remark             string            `json:"remark,omitempty"`
	DomainName         string            `json:"domainName,omitempty"`
	Aliases            []string          `json:"aliases,omitempty"`
	ResourceRecordType string            `json:"resourceRecordType,omitempty"`
	PoolLbMode         string            `json:"poolLbMode,omitempty"`
	Pools              []PointerGSLBPool `json:"pools,omitempty"`
}

type GSLBPool struct {
	Class              string               `json:"class"`
	Label              string               `json:"label,omitempty"`
	Remark             string               `json:"remark,omitempty"`
	Enabled            bool                 `json:"enabled,omitempty"`
	ResourceRecordType string               `json:"resourceRecordType,omitempty"`
	Members            []GSLBPoolMember     `json:"members,omitempty"`
	Monitors           []PointerGSLBMonitor `json:"monitors,omitempty"`
	TTL                int                  `json:"ttl,omitempty"`
}

type GSLBPoolMember struct {
	Label         string            `json:"label,omitempty"`
	Remark        string            `json:"remark,omitempty"`
	DependsOn     string            `json:"depends_on,omitempty"`
	Ratio         int               `json:"ratio,omitempty"`
	Server        PointerGSLBServer `json:"server,omitempty"`
	VirtualServer string            `json:"virtualServer,omitempty"`
	DomainName    string            `json:"domainName,omitempty"`
}

type GSLBPoolMemberA GSLBPoolMember
type GSLBPoolMemberAAAA GSLBPoolMember
type GSLBPoolMemberCNAME GSLBPoolMember
type GSLBPoolMemberMX GSLBPoolMember

type GSLBMonitor struct {
	Class       string `json:"class"`
	Label       string `json:"label,omitempty"`
	Remark      string `json:"remark,omitempty"`
	MonitorType string `json:"monitorType"`
	Interval    int64  `json:"interval,omitempty"`
	Timeout     int64  `json:"timeout,omitempty"`
	Receive     string `json:"receive,omitempty"`
	Send        string `json:"send,omitempty"`
}

type GSLBMonitorHTTP GSLBMonitor
type GSLBMonitorHTTS GSLBMonitor
type GSLBMonitorICMP GSLBMonitor
type GSLBMonitorTCP GSLBMonitor
type GSLBMonitorUDP GSLBMonitor

type GSLBServer struct {
	Class                    string                `json:"class"`
	Label                    string                `json:"label,omitempty"`
	Remark                   string                `json:"remark,omitempty"`
	DataCenter               PointerGSLBDataCenter `json:"dataCenter,omitempty"`
	Devices                  []GSLBServerDevice    `json:"devices,omitempty"`
	VirtualServers           []GSLBVirtualServer   `json:"virtualServers,omitempty"`
	Monitors                 []PointerGSLBMonitor  `json:"monitors,omitempty"`
	ServiceCheckProbeEnabled bool                  `json:"serviceCheckProbeEnabled"`
	SnmpProbeEnabled         bool                  `json:"snmpProbeEnabled"`
	PathProbeEnabled         bool                  `json:"pathProbeEnabled"`
}

type GSLBServerDevice struct {
	Label   string `json:"label,omitempty"`
	Remark  string `json:"remark,omitempty"`
	Address string `json:"address,omitempty"`
}

type GSLBVirtualServer struct {
	Label    string               `json:"label,omitempty"`
	Remark   string               `json:"remark,omitempty"`
	Address  string               `json:"address,omitempty"`
	Enabled  bool                 `json:"enabled,omitempty"`
	Monitors []PointerGSLBMonitor `json:"monitors,omitempty"`
	Name     string               `json:"name,omitempty"`
	Port     uint32               `json:"port,omitempty"`
}

type Pointer struct {
	Use   string `json:"use,omitempty"`
	BigIP string `json:"bigip,omitempty"`
}

type PointerGSLBDataCenter Pointer
type PointerGSLBDomainA Pointer
type PointerGSLBDomainAAAA Pointer
type PointerGSLBDomainCNAME Pointer
type PointerGSLBDomainMX Pointer
type PointerGSLBMonitor Pointer
type PointerGSLBPool Pointer
type PointerGSLBPoolA Pointer
type PointerGSLBPoolAAAA Pointer
type PointerGSLBPoolCNAME Pointer
type PointerGSLBPoolMemberA Pointer
type PointerGSLBPoolMemberAAAA Pointer
type PointerGSLBPoolMemberCNAME Pointer
type PointerGSLBPoolMemberMX Pointer
type PointerGSLBPoolMX Pointer
type PointerGSLBProberPool Pointer
type PointerGSLBProberPoolMember Pointer
type PointerGSLBServer Pointer
type PointerGSLBServerDevice Pointer
type PointerGSLBTopologyRegion Pointer
type PointerGSLBVirtualServer Pointer
