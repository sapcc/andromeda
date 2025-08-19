// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"testing"

	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDatacenters(t *testing.T) {
	assert := assert.New(t)

	t.Run("When RPC call fails", func(t *testing.T) {
		client := new(mockedRPCClient)
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
		client := new(mockedRPCClient)
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
		client := new(mockedRPCClient)
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
		client := new(mockedRPCClient)
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
