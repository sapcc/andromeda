// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: internal/rpc/server/rpc_server.proto

package server

import (
	fmt "fmt"
	rpcmodels "github.com/sapcc/andromeda/internal/rpcmodels"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for RPCServer service

func NewRPCServerEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for RPCServer service

type RPCServerService interface {
	UpdateProvisioningStatus(ctx context.Context, in *ProvisioningStatusRequest, opts ...client.CallOption) (*ProvisioningStatusResponse, error)
	UpdateMemberStatus(ctx context.Context, in *MemberStatusRequest, opts ...client.CallOption) (*MemberStatusResponse, error)
	GetDomains(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*DomainsResponse, error)
	GetPools(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*PoolsResponse, error)
	GetMonitors(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*MonitorsResponse, error)
	GetDatacenters(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*DatacentersResponse, error)
	GetMembers(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*MembersResponse, error)
	GetGeomaps(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*GeomapsResponse, error)
	UpdateDatacenterMeta(ctx context.Context, in *DatacenterMetaRequest, opts ...client.CallOption) (*rpcmodels.Datacenter, error)
}

type rPCServerService struct {
	c    client.Client
	name string
}

func NewRPCServerService(name string, c client.Client) RPCServerService {
	return &rPCServerService{
		c:    c,
		name: name,
	}
}

func (c *rPCServerService) UpdateProvisioningStatus(ctx context.Context, in *ProvisioningStatusRequest, opts ...client.CallOption) (*ProvisioningStatusResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.UpdateProvisioningStatus", in)
	out := new(ProvisioningStatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) UpdateMemberStatus(ctx context.Context, in *MemberStatusRequest, opts ...client.CallOption) (*MemberStatusResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.UpdateMemberStatus", in)
	out := new(MemberStatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) GetDomains(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*DomainsResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetDomains", in)
	out := new(DomainsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) GetPools(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*PoolsResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetPools", in)
	out := new(PoolsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) GetMonitors(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*MonitorsResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetMonitors", in)
	out := new(MonitorsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) GetDatacenters(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*DatacentersResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetDatacenters", in)
	out := new(DatacentersResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) GetMembers(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*MembersResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetMembers", in)
	out := new(MembersResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) GetGeomaps(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*GeomapsResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetGeomaps", in)
	out := new(GeomapsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerService) UpdateDatacenterMeta(ctx context.Context, in *DatacenterMetaRequest, opts ...client.CallOption) (*rpcmodels.Datacenter, error) {
	req := c.c.NewRequest(c.name, "RPCServer.UpdateDatacenterMeta", in)
	out := new(rpcmodels.Datacenter)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RPCServer service

type RPCServerHandler interface {
	UpdateProvisioningStatus(context.Context, *ProvisioningStatusRequest, *ProvisioningStatusResponse) error
	UpdateMemberStatus(context.Context, *MemberStatusRequest, *MemberStatusResponse) error
	GetDomains(context.Context, *SearchRequest, *DomainsResponse) error
	GetPools(context.Context, *SearchRequest, *PoolsResponse) error
	GetMonitors(context.Context, *SearchRequest, *MonitorsResponse) error
	GetDatacenters(context.Context, *SearchRequest, *DatacentersResponse) error
	GetMembers(context.Context, *SearchRequest, *MembersResponse) error
	GetGeomaps(context.Context, *SearchRequest, *GeomapsResponse) error
	UpdateDatacenterMeta(context.Context, *DatacenterMetaRequest, *rpcmodels.Datacenter) error
}

func RegisterRPCServerHandler(s server.Server, hdlr RPCServerHandler, opts ...server.HandlerOption) error {
	type rPCServer interface {
		UpdateProvisioningStatus(ctx context.Context, in *ProvisioningStatusRequest, out *ProvisioningStatusResponse) error
		UpdateMemberStatus(ctx context.Context, in *MemberStatusRequest, out *MemberStatusResponse) error
		GetDomains(ctx context.Context, in *SearchRequest, out *DomainsResponse) error
		GetPools(ctx context.Context, in *SearchRequest, out *PoolsResponse) error
		GetMonitors(ctx context.Context, in *SearchRequest, out *MonitorsResponse) error
		GetDatacenters(ctx context.Context, in *SearchRequest, out *DatacentersResponse) error
		GetMembers(ctx context.Context, in *SearchRequest, out *MembersResponse) error
		GetGeomaps(ctx context.Context, in *SearchRequest, out *GeomapsResponse) error
		UpdateDatacenterMeta(ctx context.Context, in *DatacenterMetaRequest, out *rpcmodels.Datacenter) error
	}
	type RPCServer struct {
		rPCServer
	}
	h := &rPCServerHandler{hdlr}
	return s.Handle(s.NewHandler(&RPCServer{h}, opts...))
}

type rPCServerHandler struct {
	RPCServerHandler
}

func (h *rPCServerHandler) UpdateProvisioningStatus(ctx context.Context, in *ProvisioningStatusRequest, out *ProvisioningStatusResponse) error {
	return h.RPCServerHandler.UpdateProvisioningStatus(ctx, in, out)
}

func (h *rPCServerHandler) UpdateMemberStatus(ctx context.Context, in *MemberStatusRequest, out *MemberStatusResponse) error {
	return h.RPCServerHandler.UpdateMemberStatus(ctx, in, out)
}

func (h *rPCServerHandler) GetDomains(ctx context.Context, in *SearchRequest, out *DomainsResponse) error {
	return h.RPCServerHandler.GetDomains(ctx, in, out)
}

func (h *rPCServerHandler) GetPools(ctx context.Context, in *SearchRequest, out *PoolsResponse) error {
	return h.RPCServerHandler.GetPools(ctx, in, out)
}

func (h *rPCServerHandler) GetMonitors(ctx context.Context, in *SearchRequest, out *MonitorsResponse) error {
	return h.RPCServerHandler.GetMonitors(ctx, in, out)
}

func (h *rPCServerHandler) GetDatacenters(ctx context.Context, in *SearchRequest, out *DatacentersResponse) error {
	return h.RPCServerHandler.GetDatacenters(ctx, in, out)
}

func (h *rPCServerHandler) GetMembers(ctx context.Context, in *SearchRequest, out *MembersResponse) error {
	return h.RPCServerHandler.GetMembers(ctx, in, out)
}

func (h *rPCServerHandler) GetGeomaps(ctx context.Context, in *SearchRequest, out *GeomapsResponse) error {
	return h.RPCServerHandler.GetGeomaps(ctx, in, out)
}

func (h *rPCServerHandler) UpdateDatacenterMeta(ctx context.Context, in *DatacenterMetaRequest, out *rpcmodels.Datacenter) error {
	return h.RPCServerHandler.UpdateDatacenterMeta(ctx, in, out)
}
