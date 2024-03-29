// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domains

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteDomainsDomainIDHandlerFunc turns a function with the right signature into a delete domains domain ID handler
type DeleteDomainsDomainIDHandlerFunc func(DeleteDomainsDomainIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteDomainsDomainIDHandlerFunc) Handle(params DeleteDomainsDomainIDParams) middleware.Responder {
	return fn(params)
}

// DeleteDomainsDomainIDHandler interface for that can handle valid delete domains domain ID params
type DeleteDomainsDomainIDHandler interface {
	Handle(DeleteDomainsDomainIDParams) middleware.Responder
}

// NewDeleteDomainsDomainID creates a new http.Handler for the delete domains domain ID operation
func NewDeleteDomainsDomainID(ctx *middleware.Context, handler DeleteDomainsDomainIDHandler) *DeleteDomainsDomainID {
	return &DeleteDomainsDomainID{Context: ctx, Handler: handler}
}

/*
	DeleteDomainsDomainID swagger:route DELETE /domains/{domain_id} Domains deleteDomainsDomainId

Delete a domain
*/
type DeleteDomainsDomainID struct {
	Context *middleware.Context
	Handler DeleteDomainsDomainIDHandler
}

func (o *DeleteDomainsDomainID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteDomainsDomainIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
