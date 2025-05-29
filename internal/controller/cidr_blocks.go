// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"

	"github.com/sapcc/andromeda/internal/auth"
	akamai "github.com/sapcc/andromeda/internal/rpc/agent/akamai"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type CidrBlocksController struct {
	CommonController
	cache map[string]*rpcmodels.GetCidrsResponse
	agent akamai.RPCAgentAkamaiClient
}

func NewCidrBlocksController(cc CommonController) CidrBlocksController {
	// Initialize the cache for CIDR blocks
	cache := make(map[string]*rpcmodels.GetCidrsResponse)
	return CidrBlocksController{
		CommonController: cc,
		cache:            cache,
		agent:            akamai.NewRPCAgentAkamaiClient(cc.rpc),
	}
}

// GetCidrBlocks GET /cidr-blocks
func (c CidrBlocksController) GetCidrBlocks(params administrative.GetCidrBlocksParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return administrative.NewGetCidrBlocksDefault(403).WithPayload(utils.PolicyForbidden)
	}

	provider := "akamai"
	if params.Provider != nil {
		provider = *params.Provider
	}

	// Check if the CIDR blocks are already cached
	var res *rpcmodels.GetCidrsResponse
	var ok bool
	if res, ok = c.cache[provider]; !ok {
		ctx, cancel := context.WithTimeout(params.HTTPRequest.Context(), 5*time.Second)
		defer cancel()

		var err error
		if res, err = c.agent.GetCidrs(ctx, &rpcmodels.GetCidrsRequest{}); err != nil {
			log.WithError(err).Error("failed to get cidr blocks")
			return administrative.NewGetCidrBlocksDefault(408).WithPayload(utils.TryAgainLater)
		}

		c.cache[provider] = res
	}
	return administrative.NewGetCidrBlocksOK().WithPayload(res)
}
