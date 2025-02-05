/*
 *   Copyright 2025 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

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
