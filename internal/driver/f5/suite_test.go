// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"

	"github.com/actatum/stormrpc"
	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/mock"
)

type mockedAS3Client struct {
	mock.Mock
}

func (c *mockedAS3Client) APICall(options *bigip.APIRequest) ([]byte, error) {
	args := c.Called(options)
	return args.Get(0).([]byte), args.Error(1)
}

type mockedBigIPSession struct {
	bigIP
	mock.Mock
}

func (s *mockedBigIPSession) APICall(options *bigip.APIRequest) ([]byte, error) {
	args := s.Called(options)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *mockedBigIPSession) GetDevices() ([]bigip.Device, error) {
	args := s.Called()
	return args.Get(0).([]bigip.Device), args.Error(1)
}

func (s *mockedBigIPSession) GetHost() string {
	args := s.Called()
	return args.String(0)
}

type mockedStore struct {
	mock.Mock
}

func (s *mockedStore) GetDatacenters() ([]*rpcmodels.Datacenter, error) {
	args := s.Called()
	return args.Get(0).([]*rpcmodels.Datacenter), args.Error(1)
}

func (s *mockedStore) GetDomains() ([]*rpcmodels.Domain, error) {
	args := s.Called()
	return args.Get(0).([]*rpcmodels.Domain), args.Error(1)
}

func (s *mockedStore) GetMembers(datacenterId string) ([]*rpcmodels.Member, error) {
	args := s.Called(datacenterId)
	return args.Get(0).([]*rpcmodels.Member), args.Error(1)
}

type mockedRPCClient struct {
	mock.Mock
}

func (c *mockedRPCClient) UpdateProvisioningStatus(ctx context.Context, in *server.ProvisioningStatusRequest, opts ...stormrpc.CallOption) (*server.ProvisioningStatusResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.ProvisioningStatusResponse), args.Error(1)
}

func (c *mockedRPCClient) UpdateMemberStatus(ctx context.Context, in *server.MemberStatusRequest, opts ...stormrpc.CallOption) (*server.MemberStatusResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.MemberStatusResponse), args.Error(1)
}

func (c *mockedRPCClient) GetDomains(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.DomainsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.DomainsResponse), args.Error(1)
}

func (c *mockedRPCClient) GetPools(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.PoolsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.PoolsResponse), args.Error(1)
}

func (c *mockedRPCClient) GetMonitors(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.MonitorsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.MonitorsResponse), args.Error(1)
}

func (c *mockedRPCClient) GetDatacenters(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.DatacentersResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.DatacentersResponse), args.Error(1)
}

func (c *mockedRPCClient) GetMembers(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.MembersResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.MembersResponse), args.Error(1)
}

func (c *mockedRPCClient) GetGeomaps(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.GeomapsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.GeomapsResponse), args.Error(1)
}

func (c *mockedRPCClient) UpdateDatacenterMeta(ctx context.Context, in *server.DatacenterMetaRequest, opts ...stormrpc.CallOption) (*rpcmodels.Datacenter, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*rpcmodels.Datacenter), args.Error(1)
}
