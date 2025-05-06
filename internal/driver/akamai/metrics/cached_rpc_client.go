// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"context"
	"fmt"

	"github.com/actatum/stormrpc"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/sapcc/andromeda/internal/rpc/server"
)

type CachedRPCClient struct {
	server.RPCServerClient
	cache *lru.Cache[string, string]
}

func NewCachedRPCClient(client *stormrpc.Client) *CachedRPCClient {
	cache, _ := lru.New[string, string](100)
	return &CachedRPCClient{
		RPCServerClient: server.NewRPCServerClient(client),
		cache:           cache,
	}
}

func (c *CachedRPCClient) GetProject(datacenterId string) (string, error) {
	if id, ok := c.cache.Get(datacenterId); ok {
		return id, nil
	}

	response, err := c.GetDatacenters(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: false,
		Ids:            []string{datacenterId},
	})
	if err != nil {
		return "", err
	}

	datacenters := response.GetResponse()
	if len(datacenters) == 0 {
		return "", fmt.Errorf("datacenter %s not found", datacenterId)
	}

	c.cache.Add(datacenterId, datacenters[0].ProjectId)
	return datacenters[0].ProjectId, nil
}
