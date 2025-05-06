// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/actatum/stormrpc"
	"github.com/go-openapi/runtime/middleware"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type SyncController struct {
	CommonController
}

// PostSync POST /sync
func (c SyncController) PostSync(params administrative.PostSyncParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return administrative.NewPostSyncDefault(403).WithPayload(utils.PolicyForbidden)
	}

	var domainIDs []string
	for _, domainUUID := range params.Domains.Domains {
		domainIDs = append(domainIDs, domainUUID.String())
	}

	r, err := stormrpc.NewRequest("andromeda.sync", domainIDs)
	if err != nil {
		panic(err)
	}

	if resp := c.rpc.Do(params.HTTPRequest.Context(), r); resp.Err != nil {
		panic(resp.Err)
	}
	return administrative.NewPostSyncAccepted()
}
