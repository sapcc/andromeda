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

	"github.com/sapcc/andromeda/internal/rpc/worker"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
	"go-micro.dev/v4"
)

type SyncController struct {
	sv micro.Service
}

//GetServices GET /services
func (c SyncController) PostSync(params administrative.PostSyncParams) middleware.Responder {
	projectID, err := auth.ProjectScopeForRequest(params.HTTPRequest)
	if err != nil {
		panic(err)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, projectID) {
		return GetPolicyForbiddenResponse()
	}

	pub1 := micro.NewEvent("andromeda.sync_all", c.sv.Client())
	if err := pub1.Publish(context.Background(), &worker.SyncRequest{}); err != nil {
		panic(err)
	}
	return administrative.NewPostSyncAccepted()

	/*worker := worker.NewRPCWorkerService("andromeda.agent.*", c.sv.Client())
	res, err := worker.SyncAll(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	if res.Accepted {
		return administrative.NewPostSyncAccepted()
	} else {
		return middleware.NotImplemented("Blub")
	}*/
}
