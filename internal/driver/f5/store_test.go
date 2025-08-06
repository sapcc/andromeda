// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"errors"
	"testing"

	"github.com/actatum/stormrpc"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedClient struct {
	mock.Mock
}

func (c *mockedClient) UpdateProvisioningStatus(ctx context.Context, in *server.ProvisioningStatusRequest, opts ...stormrpc.CallOption) (*server.ProvisioningStatusResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.ProvisioningStatusResponse), args.Error(1)
}

func (c *mockedClient) UpdateMemberStatus(ctx context.Context, in *server.MemberStatusRequest, opts ...stormrpc.CallOption) (*server.MemberStatusResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.MemberStatusResponse), args.Error(1)
}

func (c *mockedClient) GetDomains(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.DomainsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.DomainsResponse), args.Error(1)
}

func (c *mockedClient) GetPools(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.PoolsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.PoolsResponse), args.Error(1)
}

func (c *mockedClient) GetMonitors(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.MonitorsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.MonitorsResponse), args.Error(1)
}

func (c *mockedClient) GetDatacenters(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.DatacentersResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.DatacentersResponse), args.Error(1)
}

func (c *mockedClient) GetMembers(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.MembersResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.MembersResponse), args.Error(1)
}

func (c *mockedClient) GetGeomaps(ctx context.Context, in *server.SearchRequest, opts ...stormrpc.CallOption) (*server.GeomapsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*server.GeomapsResponse), args.Error(1)
}

func (c *mockedClient) UpdateDatacenterMeta(ctx context.Context, in *server.DatacenterMetaRequest, opts ...stormrpc.CallOption) (*rpcmodels.Datacenter, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*rpcmodels.Datacenter), args.Error(1)
}

func TestGetDatacenters(t *testing.T) {
	assert := assert.New(t)
	t.Run("When RPC call fails", func(t *testing.T) {
		client := new(mockedClient)
		client.
			On("GetDatacenters",
				mock.Anything,
				mock.Anything,
				mock.Anything).
			Return(&server.DatacentersResponse{},
				errors.New("RPC failed"))
		store := andromedaF5Store{rpc: client}
		_, err := store.GetDatacenters()
		assert.NotNil(err, "Expected store.GetDatacenters() to have returned an error")
	})
	t.Run("When RPC call returns at least 1 datacenter", func(t *testing.T) {
		client := new(mockedClient)
		client.
			On("GetDatacenters",
				mock.Anything,
				mock.Anything,
				mock.Anything).
			Return(&server.DatacentersResponse{
				Response: []*rpcmodels.Datacenter{
					{Id: "dc1"},
				},
			}, nil)
		store := andromedaF5Store{rpc: client}
		expected := []*rpcmodels.Datacenter{
			{Id: "dc1"},
		}
		datacenters, err := store.GetDatacenters()
		assert.Nil(err, "Expected store.GetDatacenters() to not have returned an error")
		assert.Equal(expected, datacenters)
	})
}

func TestGetMembers(t *testing.T) {
	assert := assert.New(t)
	t.Run("When RPC call fails", func(t *testing.T) {
		client := new(mockedClient)
		client.
			On("GetMembers",
				mock.Anything,
				mock.Anything,
				mock.Anything).
			Return(&server.MembersResponse{},
				errors.New("RPC failed"))
		store := andromedaF5Store{rpc: client}
		_, err := store.GetMembers("dc1")
		assert.NotNil(err, "Expected store.GetMembers() to have returned an error")
	})
	t.Run("When RPC call returns at least 1 member", func(t *testing.T) {
		client := new(mockedClient)
		client.
			On("GetMembers",
				mock.Anything,
				mock.MatchedBy(func(req *server.SearchRequest) bool { return req.DatacenterId == "dc1" }),
				mock.Anything,
			).
			Return(&server.MembersResponse{
				Response: []*rpcmodels.Member{
					{Id: "member1"},
				},
			}, nil)
		store := andromedaF5Store{rpc: client}
		expected := []*rpcmodels.Member{
			{Id: "member1"},
		}
		datacenters, err := store.GetMembers("dc1")
		assert.Nil(err, "Expected store.GetMembers() to not have returned an error")
		assert.Equal(expected, datacenters)
	})
}
