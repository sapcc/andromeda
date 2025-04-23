/*
 *   Copyright 2021 SAP SE
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

package controller

import (
	"encoding/json"
	"fmt"

	"github.com/actatum/stormrpc"
	"github.com/go-openapi/runtime/middleware"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type cidrBlocks []map[string]any

type CidrBlocksController struct {
	CommonController
	cache map[string]cidrBlocks
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
	var res cidrBlocks
	var ok bool
	if res, ok = c.cache[provider]; !ok {
		subject := fmt.Sprintf("andromeda.get_cidrs.%s", provider)
		r, err := stormrpc.NewRequest(subject, nil)
		if err != nil {
			panic(err)
		}

		resp := c.rpc.Do(params.HTTPRequest.Context(), r)
		if resp.Err != nil {
			panic(resp.Err)
		}

		if err = json.Unmarshal(resp.Data, &res); err != nil {
			panic(err)
		}
		c.cache[provider] = res
	}
	return administrative.NewGetCidrBlocksOK().WithPayload(res)
}
