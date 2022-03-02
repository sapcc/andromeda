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
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/models"
)

// GetDomainsHandlerFunc turns a function with the right signature into a get domains handler
type GetDomainsHandlerFunc func(GetDomainsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDomainsHandlerFunc) Handle(params GetDomainsParams) middleware.Responder {
	return fn(params)
}

// GetDomainsHandler interface for that can handle valid get domains params
type GetDomainsHandler interface {
	Handle(GetDomainsParams) middleware.Responder
}

// NewGetDomains creates a new http.Handler for the get domains operation
func NewGetDomains(ctx *middleware.Context, handler GetDomainsHandler) *GetDomains {
	return &GetDomains{Context: ctx, Handler: handler}
}

/* GetDomains swagger:route GET /domains Domains getDomains

List domains

*/
type GetDomains struct {
	Context *middleware.Context
	Handler GetDomainsHandler
}

func (o *GetDomains) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetDomainsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetDomainsOKBody get domains o k body
//
// swagger:model GetDomainsOKBody
type GetDomainsOKBody struct {

	// domains
	Domains []*models.Domain `json:"domains"`

	// links
	Links []*models.Link `json:"links,omitempty"`
}

// Validate validates this get domains o k body
func (o *GetDomainsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomains(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDomainsOKBody) validateDomains(formats strfmt.Registry) error {
	if swag.IsZero(o.Domains) { // not required
		return nil
	}

	for i := 0; i < len(o.Domains); i++ {
		if swag.IsZero(o.Domains[i]) { // not required
			continue
		}

		if o.Domains[i] != nil {
			if err := o.Domains[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetDomainsOKBody) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(o.Links) { // not required
		return nil
	}

	for i := 0; i < len(o.Links); i++ {
		if swag.IsZero(o.Links[i]) { // not required
			continue
		}

		if o.Links[i] != nil {
			if err := o.Links[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get domains o k body based on the context it is used
func (o *GetDomainsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDomains(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDomainsOKBody) contextValidateDomains(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Domains); i++ {

		if o.Domains[i] != nil {
			if err := o.Domains[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetDomainsOKBody) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Links); i++ {

		if o.Links[i] != nil {
			if err := o.Links[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetDomainsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDomainsOKBody) UnmarshalBinary(b []byte) error {
	var res GetDomainsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
