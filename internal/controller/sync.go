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
	"context"

	"github.com/go-openapi/runtime/middleware"
	"go-micro.dev/v4"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/rpc/worker"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type SyncController struct {
	sv micro.Service
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

	ev := &worker.SyncRequest{
		DomainIds: domainIDs,
	}

	pub1 := micro.NewEvent("andromeda.sync", c.sv.Client())
	if err := pub1.Publish(context.Background(), ev); err != nil {
		panic(err)
	}
	return administrative.NewPostSyncAccepted()
}
