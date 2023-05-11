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

package administrative

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostSyncHandlerFunc turns a function with the right signature into a post sync handler
type PostSyncHandlerFunc func(PostSyncParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostSyncHandlerFunc) Handle(params PostSyncParams) middleware.Responder {
	return fn(params)
}

// PostSyncHandler interface for that can handle valid post sync params
type PostSyncHandler interface {
	Handle(PostSyncParams) middleware.Responder
}

// NewPostSync creates a new http.Handler for the post sync operation
func NewPostSync(ctx *middleware.Context, handler PostSyncHandler) *PostSync {
	return &PostSync{Context: ctx, Handler: handler}
}

/*
	PostSync swagger:route POST /sync Administrative postSync

Enqueue a full sync
*/
type PostSync struct {
	Context *middleware.Context
	Handler PostSyncHandler
}

func (o *PostSync) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostSyncParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostSyncBody post sync body
//
// swagger:model PostSyncBody
type PostSyncBody struct {

	// domains
	// Required: true
	Domains []strfmt.UUID `json:"domains"`
}

// Validate validates this post sync body
func (o *PostSyncBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomains(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostSyncBody) validateDomains(formats strfmt.Registry) error {

	if err := validate.Required("domains"+"."+"domains", "body", o.Domains); err != nil {
		return err
	}

	for i := 0; i < len(o.Domains); i++ {

		if err := validate.FormatOf("domains"+"."+"domains"+"."+strconv.Itoa(i), "body", "uuid", o.Domains[i].String(), formats); err != nil {
			return err
		}

	}

	return nil
}

// ContextValidate validates this post sync body based on context it is used
func (o *PostSyncBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostSyncBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostSyncBody) UnmarshalBinary(b []byte) error {
	var res PostSyncBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}