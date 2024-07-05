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
	"github.com/go-openapi/runtime/middleware"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type ServiceController struct {
	CommonController
}

// GetServices GET /services
func (c ServiceController) GetServices(params administrative.GetServicesParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return administrative.NewGetServicesDefault(403).WithPayload(utils.PolicyForbidden)
	}

	//goland:noinspection GoPreferNilSlice
	var responseServices = []*models.Service{}

	// Unsupported by stormrpc
	/*_listServices, err := c.rpc.Options().Registry.ListServices()
	if err != nil {
		panic(err)
	}
	for _, _service := range _listServices {
		if _service.Name == "" {
			continue
		}

		_svs, err := c.rpc.Options().Registry.GetService(_service.Name)
		if err != nil {
			panic(err)
		}
		for _, _service := range _svs {
			for _, _node := range _service.Nodes {
				responseServices = append(responseServices, &models.Service{
					ID:         _node.Id,
					RPCAddress: _node.Address,
					Provider:   getMetadata(_node.Metadata, "type"),
					Host:       strfmt.Hostname(getMetadata(_node.Metadata, "host")),
					Type:       _service.Name,
					Version:    _service.Version,
					Metadata:   _node.Metadata,
				})
			}
		}
	}
	*/

	return administrative.NewGetServicesOK().WithPayload(&administrative.GetServicesOKBody{Services: responseServices})
}
