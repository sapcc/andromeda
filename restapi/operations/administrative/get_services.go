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

	"github.com/sapcc/andromeda/models"
)

// GetServicesHandlerFunc turns a function with the right signature into a get services handler
type GetServicesHandlerFunc func(GetServicesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetServicesHandlerFunc) Handle(params GetServicesParams) middleware.Responder {
	return fn(params)
}

// GetServicesHandler interface for that can handle valid get services params
type GetServicesHandler interface {
	Handle(GetServicesParams) middleware.Responder
}

// NewGetServices creates a new http.Handler for the get services operation
func NewGetServices(ctx *middleware.Context, handler GetServicesHandler) *GetServices {
	return &GetServices{Context: ctx, Handler: handler}
}

/*
	GetServices swagger:route GET /services Administrative getServices

List Services
*/
type GetServices struct {
	Context *middleware.Context
	Handler GetServicesHandler
}

func (o *GetServices) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetServicesParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetServicesOKBody get services o k body
//
// swagger:model GetServicesOKBody
type GetServicesOKBody struct {

	// services
	Services []*models.Service `json:"services"`
}

// Validate validates this get services o k body
func (o *GetServicesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateServices(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetServicesOKBody) validateServices(formats strfmt.Registry) error {
	if swag.IsZero(o.Services) { // not required
		return nil
	}

	for i := 0; i < len(o.Services); i++ {
		if swag.IsZero(o.Services[i]) { // not required
			continue
		}

		if o.Services[i] != nil {
			if err := o.Services[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get services o k body based on the context it is used
func (o *GetServicesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateServices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetServicesOKBody) contextValidateServices(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Services); i++ {

		if o.Services[i] != nil {
			if err := o.Services[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetServicesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetServicesOKBody) UnmarshalBinary(b []byte) error {
	var res GetServicesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
