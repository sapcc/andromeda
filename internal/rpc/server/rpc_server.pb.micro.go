// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: internal/rpc/server/rpc_server.proto

package server

import (
	fmt "fmt"
	_ "github.com/sapcc/andromeda/internal/models"
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
	GetMembers(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*MembersResponse, error)
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

func (c *rPCServerService) GetMembers(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*MembersResponse, error) {
	req := c.c.NewRequest(c.name, "RPCServer.GetMembers", in)
	out := new(MembersResponse)
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
	GetMembers(context.Context, *SearchRequest, *MembersResponse) error
}

func RegisterRPCServerHandler(s server.Server, hdlr RPCServerHandler, opts ...server.HandlerOption) error {
	type rPCServer interface {
		UpdateProvisioningStatus(ctx context.Context, in *ProvisioningStatusRequest, out *ProvisioningStatusResponse) error
		UpdateMemberStatus(ctx context.Context, in *MemberStatusRequest, out *MemberStatusResponse) error
		GetDomains(ctx context.Context, in *SearchRequest, out *DomainsResponse) error
		GetMembers(ctx context.Context, in *SearchRequest, out *MembersResponse) error
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

func (h *rPCServerHandler) GetMembers(ctx context.Context, in *SearchRequest, out *MembersResponse) error {
	return h.RPCServerHandler.GetMembers(ctx, in, out)
}
