// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"fmt"

	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

type AndromedaF5Store interface {
	GetDatacenters() ([]*rpcmodels.Datacenter, error)
	GetDomains() ([]*rpcmodels.Domain, error)
	GetMembers(datacenterID string) ([]*rpcmodels.Member, error)
}

type andromedaF5Store struct {
	rpc server.RPCServerClient
}

func NewAndromedaF5Store(c server.RPCServerClient) AndromedaF5Store {
	return &andromedaF5Store{rpc: c}
}

func (s *andromedaF5Store) GetDatacenters() ([]*rpcmodels.Datacenter, error) {
	// AS3 POST /declare payload must include *all* datacenters
	res, err := s.rpc.GetDatacenters(context.Background(), &server.SearchRequest{
		Provider:       "f5",
		PageNumber:     0,
		ResultPerPage:  1000,
		FullyPopulated: false,
	})
	if err != nil {
		return nil, fmt.Errorf("rpc.GetDatacenters failed: %s", err)
	}
	if res == nil || len(res.GetResponse()) == 0 {
		return nil, fmt.Errorf("no F5 datacenters found")
	}
	return res.GetResponse(), nil
}

func (s *andromedaF5Store) GetDomains() ([]*rpcmodels.Domain, error) {
	// AS3 POST /declare payload must include *all* domains
	res, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "f5",
		PageNumber:     0,
		ResultPerPage:  1000, // TODO: make it possible to go over all results
		FullyPopulated: true,
	})
	if err != nil {
		return nil, fmt.Errorf("rpc.GetDomains failed: %s", err)
	}
	if res == nil || len(res.GetResponse()) == 0 {
		return nil, fmt.Errorf("no F5 domains found")
	}
	return res.GetResponse(), nil
}

func (s *andromedaF5Store) GetMembers(datacenterId string) ([]*rpcmodels.Member, error) {
	// AS3 POST /declare payload must include *all* members (servers)
	res, err := s.rpc.GetMembers(context.Background(), &server.SearchRequest{
		DatacenterId:  datacenterId,
		PageNumber:    0,
		ResultPerPage: 1000, // TODO: make it possible to go over all results
	})
	if err != nil {
		return nil, fmt.Errorf("rpc.GetMembers failed: %s", err)
	}
	if res == nil {
		return nil, fmt.Errorf("rpc.GetMembers response is nil")
	}
	return res.GetResponse(), nil
}
