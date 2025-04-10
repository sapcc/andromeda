// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: internal/rpcmodels/rpc_models.proto

package rpcmodels

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Monitor_MonitorType int32

const (
	Monitor_HTTP  Monitor_MonitorType = 0
	Monitor_HTTPS Monitor_MonitorType = 1
	Monitor_ICMP  Monitor_MonitorType = 2
	Monitor_TCP   Monitor_MonitorType = 3
	Monitor_UDP   Monitor_MonitorType = 4
	Monitor_POP   Monitor_MonitorType = 5
	Monitor_SMTP  Monitor_MonitorType = 7
)

// Enum value maps for Monitor_MonitorType.
var (
	Monitor_MonitorType_name = map[int32]string{
		0: "HTTP",
		1: "HTTPS",
		2: "ICMP",
		3: "TCP",
		4: "UDP",
		5: "POP",
		7: "SMTP",
	}
	Monitor_MonitorType_value = map[string]int32{
		"HTTP":  0,
		"HTTPS": 1,
		"ICMP":  2,
		"TCP":   3,
		"UDP":   4,
		"POP":   5,
		"SMTP":  7,
	}
)

func (x Monitor_MonitorType) Enum() *Monitor_MonitorType {
	p := new(Monitor_MonitorType)
	*p = x
	return p
}

func (x Monitor_MonitorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Monitor_MonitorType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_rpcmodels_rpc_models_proto_enumTypes[0].Descriptor()
}

func (Monitor_MonitorType) Type() protoreflect.EnumType {
	return &file_internal_rpcmodels_rpc_models_proto_enumTypes[0]
}

func (x Monitor_MonitorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Monitor_MonitorType.Descriptor instead.
func (Monitor_MonitorType) EnumDescriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{6, 0}
}

type Monitor_HttpMethod int32

const (
	Monitor_GET     Monitor_HttpMethod = 0
	Monitor_POST    Monitor_HttpMethod = 1
	Monitor_PUT     Monitor_HttpMethod = 2
	Monitor_HEAD    Monitor_HttpMethod = 3
	Monitor_PATCH   Monitor_HttpMethod = 4
	Monitor_DELETE  Monitor_HttpMethod = 5
	Monitor_OPTIONS Monitor_HttpMethod = 6
)

// Enum value maps for Monitor_HttpMethod.
var (
	Monitor_HttpMethod_name = map[int32]string{
		0: "GET",
		1: "POST",
		2: "PUT",
		3: "HEAD",
		4: "PATCH",
		5: "DELETE",
		6: "OPTIONS",
	}
	Monitor_HttpMethod_value = map[string]int32{
		"GET":     0,
		"POST":    1,
		"PUT":     2,
		"HEAD":    3,
		"PATCH":   4,
		"DELETE":  5,
		"OPTIONS": 6,
	}
)

func (x Monitor_HttpMethod) Enum() *Monitor_HttpMethod {
	p := new(Monitor_HttpMethod)
	*p = x
	return p
}

func (x Monitor_HttpMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Monitor_HttpMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_rpcmodels_rpc_models_proto_enumTypes[1].Descriptor()
}

func (Monitor_HttpMethod) Type() protoreflect.EnumType {
	return &file_internal_rpcmodels_rpc_models_proto_enumTypes[1]
}

func (x Monitor_HttpMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Monitor_HttpMethod.Descriptor instead.
func (Monitor_HttpMethod) EnumDescriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{6, 1}
}

type Domain struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AdminStateUp       bool                   `protobuf:"varint,2,opt,name=admin_state_up,json=adminStateUp,proto3" json:"admin_state_up,omitempty"`
	Aliases            []string               `protobuf:"bytes,3,rep,name=aliases,proto3" json:"aliases,omitempty"`
	Fqdn               string                 `protobuf:"bytes,4,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
	Mode               string                 `protobuf:"bytes,5,opt,name=mode,proto3" json:"mode,omitempty"`
	Pools              []*Pool                `protobuf:"bytes,6,rep,name=pools,proto3" json:"pools,omitempty"`
	RecordType         string                 `protobuf:"bytes,7,opt,name=record_type,json=recordType,proto3" json:"record_type,omitempty"`
	Datacenters        []*Datacenter          `protobuf:"bytes,8,rep,name=datacenters,proto3" json:"datacenters,omitempty"`
	ProvisioningStatus string                 `protobuf:"bytes,9,opt,name=provisioning_status,json=provisioningStatus,proto3" json:"provisioning_status,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Domain) Reset() {
	*x = Domain{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Domain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Domain) ProtoMessage() {}

func (x *Domain) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Domain.ProtoReflect.Descriptor instead.
func (*Domain) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{0}
}

func (x *Domain) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Domain) GetAdminStateUp() bool {
	if x != nil {
		return x.AdminStateUp
	}
	return false
}

func (x *Domain) GetAliases() []string {
	if x != nil {
		return x.Aliases
	}
	return nil
}

func (x *Domain) GetFqdn() string {
	if x != nil {
		return x.Fqdn
	}
	return ""
}

func (x *Domain) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

func (x *Domain) GetPools() []*Pool {
	if x != nil {
		return x.Pools
	}
	return nil
}

func (x *Domain) GetRecordType() string {
	if x != nil {
		return x.RecordType
	}
	return ""
}

func (x *Domain) GetDatacenters() []*Datacenter {
	if x != nil {
		return x.Datacenters
	}
	return nil
}

func (x *Domain) GetProvisioningStatus() string {
	if x != nil {
		return x.ProvisioningStatus
	}
	return ""
}

type Pool struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AdminStateUp       bool                   `protobuf:"varint,2,opt,name=admin_state_up,json=adminStateUp,proto3" json:"admin_state_up,omitempty"`
	Members            []*Member              `protobuf:"bytes,3,rep,name=members,proto3" json:"members,omitempty"`
	Monitors           []*Monitor             `protobuf:"bytes,4,rep,name=monitors,proto3" json:"monitors,omitempty"`
	ProvisioningStatus string                 `protobuf:"bytes,5,opt,name=provisioning_status,json=provisioningStatus,proto3" json:"provisioning_status,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Pool) Reset() {
	*x = Pool{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pool) ProtoMessage() {}

func (x *Pool) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pool.ProtoReflect.Descriptor instead.
func (*Pool) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{1}
}

func (x *Pool) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Pool) GetAdminStateUp() bool {
	if x != nil {
		return x.AdminStateUp
	}
	return false
}

func (x *Pool) GetMembers() []*Member {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *Pool) GetMonitors() []*Monitor {
	if x != nil {
		return x.Monitors
	}
	return nil
}

func (x *Pool) GetProvisioningStatus() string {
	if x != nil {
		return x.ProvisioningStatus
	}
	return ""
}

type Datacenter struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AdminStateUp       bool                   `protobuf:"varint,2,opt,name=admin_state_up,json=adminStateUp,proto3" json:"admin_state_up,omitempty"`
	Continent          string                 `protobuf:"bytes,3,opt,name=continent,proto3" json:"continent,omitempty"`
	Country            string                 `protobuf:"bytes,4,opt,name=country,proto3" json:"country,omitempty"`
	StateOrProvince    string                 `protobuf:"bytes,5,opt,name=state_or_province,json=stateOrProvince,proto3" json:"state_or_province,omitempty"`
	City               string                 `protobuf:"bytes,6,opt,name=city,proto3" json:"city,omitempty"`
	Longitude          float64                `protobuf:"fixed64,7,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Latitude           float64                `protobuf:"fixed64,8,opt,name=latitude,proto3" json:"latitude,omitempty"`
	ProvisioningStatus string                 `protobuf:"bytes,9,opt,name=provisioning_status,json=provisioningStatus,proto3" json:"provisioning_status,omitempty"`
	Scope              string                 `protobuf:"bytes,10,opt,name=scope,proto3" json:"scope,omitempty"`
	Provider           string                 `protobuf:"bytes,11,opt,name=provider,proto3" json:"provider,omitempty"`
	Meta               int32                  `protobuf:"varint,12,opt,name=meta,proto3" json:"meta,omitempty"`
	ProjectId          string                 `protobuf:"bytes,13,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Datacenter) Reset() {
	*x = Datacenter{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Datacenter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Datacenter) ProtoMessage() {}

func (x *Datacenter) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Datacenter.ProtoReflect.Descriptor instead.
func (*Datacenter) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{2}
}

func (x *Datacenter) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Datacenter) GetAdminStateUp() bool {
	if x != nil {
		return x.AdminStateUp
	}
	return false
}

func (x *Datacenter) GetContinent() string {
	if x != nil {
		return x.Continent
	}
	return ""
}

func (x *Datacenter) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *Datacenter) GetStateOrProvince() string {
	if x != nil {
		return x.StateOrProvince
	}
	return ""
}

func (x *Datacenter) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *Datacenter) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *Datacenter) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Datacenter) GetProvisioningStatus() string {
	if x != nil {
		return x.ProvisioningStatus
	}
	return ""
}

func (x *Datacenter) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *Datacenter) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *Datacenter) GetMeta() int32 {
	if x != nil {
		return x.Meta
	}
	return 0
}

func (x *Datacenter) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type GeomapAssignment struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Datacenter    string                 `protobuf:"bytes,1,opt,name=datacenter,proto3" json:"datacenter,omitempty"`
	Countries     []string               `protobuf:"bytes,2,rep,name=countries,proto3" json:"countries,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GeomapAssignment) Reset() {
	*x = GeomapAssignment{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GeomapAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeomapAssignment) ProtoMessage() {}

func (x *GeomapAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeomapAssignment.ProtoReflect.Descriptor instead.
func (*GeomapAssignment) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{3}
}

func (x *GeomapAssignment) GetDatacenter() string {
	if x != nil {
		return x.Datacenter
	}
	return ""
}

func (x *GeomapAssignment) GetCountries() []string {
	if x != nil {
		return x.Countries
	}
	return nil
}

type Geomap struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DefaultDatacenter  string                 `protobuf:"bytes,2,opt,name=default_datacenter,json=defaultDatacenter,proto3" json:"default_datacenter,omitempty"`
	Assignment         []*GeomapAssignment    `protobuf:"bytes,3,rep,name=assignment,proto3" json:"assignment,omitempty"`
	ProvisioningStatus string                 `protobuf:"bytes,4,opt,name=provisioning_status,json=provisioningStatus,proto3" json:"provisioning_status,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Geomap) Reset() {
	*x = Geomap{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Geomap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Geomap) ProtoMessage() {}

func (x *Geomap) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Geomap.ProtoReflect.Descriptor instead.
func (*Geomap) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{4}
}

func (x *Geomap) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Geomap) GetDefaultDatacenter() string {
	if x != nil {
		return x.DefaultDatacenter
	}
	return ""
}

func (x *Geomap) GetAssignment() []*GeomapAssignment {
	if x != nil {
		return x.Assignment
	}
	return nil
}

func (x *Geomap) GetProvisioningStatus() string {
	if x != nil {
		return x.ProvisioningStatus
	}
	return ""
}

type Member struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AdminStateUp       bool                   `protobuf:"varint,2,opt,name=admin_state_up,json=adminStateUp,proto3" json:"admin_state_up,omitempty"`
	Address            uint32                 `protobuf:"fixed32,3,opt,name=address,proto3" json:"address,omitempty"`
	Port               uint32                 `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
	Datacenter         string                 `protobuf:"bytes,5,opt,name=datacenter,proto3" json:"datacenter,omitempty"`
	ProvisioningStatus string                 `protobuf:"bytes,6,opt,name=provisioning_status,json=provisioningStatus,proto3" json:"provisioning_status,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Member) Reset() {
	*x = Member{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Member) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Member) ProtoMessage() {}

func (x *Member) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Member.ProtoReflect.Descriptor instead.
func (*Member) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{5}
}

func (x *Member) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Member) GetAdminStateUp() bool {
	if x != nil {
		return x.AdminStateUp
	}
	return false
}

func (x *Member) GetAddress() uint32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *Member) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Member) GetDatacenter() string {
	if x != nil {
		return x.Datacenter
	}
	return ""
}

func (x *Member) GetProvisioningStatus() string {
	if x != nil {
		return x.ProvisioningStatus
	}
	return ""
}

type Monitor struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AdminStateUp       bool                   `protobuf:"varint,2,opt,name=admin_state_up,json=adminStateUp,proto3" json:"admin_state_up,omitempty"`
	Interval           int64                  `protobuf:"varint,3,opt,name=interval,proto3" json:"interval,omitempty"`
	PoolId             string                 `protobuf:"bytes,4,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Send               string                 `protobuf:"bytes,5,opt,name=send,proto3" json:"send,omitempty"`
	Receive            string                 `protobuf:"bytes,6,opt,name=receive,proto3" json:"receive,omitempty"`
	Timeout            int64                  `protobuf:"varint,7,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Type               Monitor_MonitorType    `protobuf:"varint,8,opt,name=type,proto3,enum=Monitor_MonitorType" json:"type,omitempty"`
	ProvisioningStatus string                 `protobuf:"bytes,9,opt,name=provisioning_status,json=provisioningStatus,proto3" json:"provisioning_status,omitempty"`
	Method             Monitor_HttpMethod     `protobuf:"varint,10,opt,name=method,proto3,enum=Monitor_HttpMethod" json:"method,omitempty"`
	DomainName         string                 `protobuf:"bytes,11,opt,name=domain_name,json=domainName,proto3" json:"domain_name,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Monitor) Reset() {
	*x = Monitor{}
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Monitor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Monitor) ProtoMessage() {}

func (x *Monitor) ProtoReflect() protoreflect.Message {
	mi := &file_internal_rpcmodels_rpc_models_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Monitor.ProtoReflect.Descriptor instead.
func (*Monitor) Descriptor() ([]byte, []int) {
	return file_internal_rpcmodels_rpc_models_proto_rawDescGZIP(), []int{6}
}

func (x *Monitor) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Monitor) GetAdminStateUp() bool {
	if x != nil {
		return x.AdminStateUp
	}
	return false
}

func (x *Monitor) GetInterval() int64 {
	if x != nil {
		return x.Interval
	}
	return 0
}

func (x *Monitor) GetPoolId() string {
	if x != nil {
		return x.PoolId
	}
	return ""
}

func (x *Monitor) GetSend() string {
	if x != nil {
		return x.Send
	}
	return ""
}

func (x *Monitor) GetReceive() string {
	if x != nil {
		return x.Receive
	}
	return ""
}

func (x *Monitor) GetTimeout() int64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *Monitor) GetType() Monitor_MonitorType {
	if x != nil {
		return x.Type
	}
	return Monitor_HTTP
}

func (x *Monitor) GetProvisioningStatus() string {
	if x != nil {
		return x.ProvisioningStatus
	}
	return ""
}

func (x *Monitor) GetMethod() Monitor_HttpMethod {
	if x != nil {
		return x.Method
	}
	return Monitor_GET
}

func (x *Monitor) GetDomainName() string {
	if x != nil {
		return x.DomainName
	}
	return ""
}

var File_internal_rpcmodels_rpc_models_proto protoreflect.FileDescriptor

var file_internal_rpcmodels_rpc_models_proto_rawDesc = string([]byte{
	0x0a, 0x23, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x72, 0x70, 0x63, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9e, 0x02, 0x0a, 0x06, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x24, 0x0a, 0x0e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f,
	0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x55, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x65, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x71, 0x64, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x71, 0x64, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x05, 0x70, 0x6f, 0x6f, 0x6c,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x05,
	0x70, 0x6f, 0x6f, 0x6c, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2d, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x63, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x63, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x73, 0x12, 0x2f, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xb6, 0x01, 0x0a, 0x04, 0x50, 0x6f, 0x6f, 0x6c, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x24, 0x0a, 0x0e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x75,
	0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x55, 0x70, 0x12, 0x21, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52,
	0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x24, 0x0a, 0x08, 0x6d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x6f, 0x6e,
	0x69, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x2f,
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x8a, 0x03, 0x0a, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24,
	0x0a, 0x0e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x70,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x55, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x65,
	0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x2a, 0x0a, 0x11,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x6f, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x6e, 0x63,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x4f, 0x72,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79, 0x12, 0x1c, 0x0a, 0x09,
	0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x2f, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e,
	0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x65, 0x74,
	0x61, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x22, 0x50, 0x0a, 0x10,
	0x47, 0x65, 0x6f, 0x6d, 0x61, 0x70, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0xab,
	0x01, 0x0a, 0x06, 0x47, 0x65, 0x6f, 0x6d, 0x61, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x12, 0x64, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x0a, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x47,
	0x65, 0x6f, 0x6d, 0x61, 0x70, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x0a, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x13, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xbd, 0x01, 0x0a,
	0x06, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0c, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x55, 0x70, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x07, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x64,
	0x61, 0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x64, 0x61, 0x74, 0x61, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x13, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x90, 0x04, 0x0a,
	0x07, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0c, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x55, 0x70, 0x12, 0x1a,
	0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f,
	0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x6f,
	0x6c, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x28, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x4d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x48, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x06, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x51, 0x0a, 0x0b, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x54, 0x54, 0x50, 0x10, 0x00, 0x12, 0x09, 0x0a,
	0x05, 0x48, 0x54, 0x54, 0x50, 0x53, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x43, 0x4d, 0x50,
	0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x43, 0x50, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x55,
	0x44, 0x50, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x4f, 0x50, 0x10, 0x05, 0x12, 0x08, 0x0a,
	0x04, 0x53, 0x4d, 0x54, 0x50, 0x10, 0x07, 0x22, 0x56, 0x0a, 0x0a, 0x48, 0x74, 0x74, 0x70, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x07, 0x0a, 0x03, 0x47, 0x45, 0x54, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x55, 0x54, 0x10,
	0x02, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x45, 0x41, 0x44, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x50,
	0x41, 0x54, 0x43, 0x48, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x10, 0x06, 0x42,
	0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x61,
	0x70, 0x63, 0x63, 0x2f, 0x61, 0x6e, 0x64, 0x72, 0x6f, 0x6d, 0x65, 0x64, 0x61, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x72, 0x70, 0x63, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_internal_rpcmodels_rpc_models_proto_rawDescOnce sync.Once
	file_internal_rpcmodels_rpc_models_proto_rawDescData []byte
)

func file_internal_rpcmodels_rpc_models_proto_rawDescGZIP() []byte {
	file_internal_rpcmodels_rpc_models_proto_rawDescOnce.Do(func() {
		file_internal_rpcmodels_rpc_models_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_rpcmodels_rpc_models_proto_rawDesc), len(file_internal_rpcmodels_rpc_models_proto_rawDesc)))
	})
	return file_internal_rpcmodels_rpc_models_proto_rawDescData
}

var file_internal_rpcmodels_rpc_models_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_internal_rpcmodels_rpc_models_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_rpcmodels_rpc_models_proto_goTypes = []any{
	(Monitor_MonitorType)(0), // 0: Monitor.MonitorType
	(Monitor_HttpMethod)(0),  // 1: Monitor.HttpMethod
	(*Domain)(nil),           // 2: Domain
	(*Pool)(nil),             // 3: Pool
	(*Datacenter)(nil),       // 4: Datacenter
	(*GeomapAssignment)(nil), // 5: GeomapAssignment
	(*Geomap)(nil),           // 6: Geomap
	(*Member)(nil),           // 7: Member
	(*Monitor)(nil),          // 8: Monitor
}
var file_internal_rpcmodels_rpc_models_proto_depIdxs = []int32{
	3, // 0: Domain.pools:type_name -> Pool
	4, // 1: Domain.datacenters:type_name -> Datacenter
	7, // 2: Pool.members:type_name -> Member
	8, // 3: Pool.monitors:type_name -> Monitor
	5, // 4: Geomap.assignment:type_name -> GeomapAssignment
	0, // 5: Monitor.type:type_name -> Monitor.MonitorType
	1, // 6: Monitor.method:type_name -> Monitor.HttpMethod
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_internal_rpcmodels_rpc_models_proto_init() }
func file_internal_rpcmodels_rpc_models_proto_init() {
	if File_internal_rpcmodels_rpc_models_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_rpcmodels_rpc_models_proto_rawDesc), len(file_internal_rpcmodels_rpc_models_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_rpcmodels_rpc_models_proto_goTypes,
		DependencyIndexes: file_internal_rpcmodels_rpc_models_proto_depIdxs,
		EnumInfos:         file_internal_rpcmodels_rpc_models_proto_enumTypes,
		MessageInfos:      file_internal_rpcmodels_rpc_models_proto_msgTypes,
	}.Build()
	File_internal_rpcmodels_rpc_models_proto = out.File
	file_internal_rpcmodels_rpc_models_proto_goTypes = nil
	file_internal_rpcmodels_rpc_models_proto_depIdxs = nil
}
