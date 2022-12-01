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

package datacenters

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/sapcc/andromeda/models"
)

// PostDatacentersHandlerFunc turns a function with the right signature into a post datacenters handler
type PostDatacentersHandlerFunc func(PostDatacentersParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostDatacentersHandlerFunc) Handle(params PostDatacentersParams) middleware.Responder {
	return fn(params)
}

// PostDatacentersHandler interface for that can handle valid post datacenters params
type PostDatacentersHandler interface {
	Handle(PostDatacentersParams) middleware.Responder
}

// NewPostDatacenters creates a new http.Handler for the post datacenters operation
func NewPostDatacenters(ctx *middleware.Context, handler PostDatacentersHandler) *PostDatacenters {
	return &PostDatacenters{Context: ctx, Handler: handler}
}

/*
	PostDatacenters swagger:route POST /datacenters Datacenters postDatacenters

Create new datacenter
*/
type PostDatacenters struct {
	Context *middleware.Context
	Handler PostDatacentersHandler
}

func (o *PostDatacenters) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostDatacentersParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostDatacentersBody post datacenters body
//
// swagger:model PostDatacentersBody
type PostDatacentersBody struct {

	// datacenter
	// Required: true
	Datacenter *models.Datacenter `json:"datacenter"`
}

// Validate validates this post datacenters body
func (o *PostDatacentersBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDatacenter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersBody) validateDatacenter(formats strfmt.Registry) error {

	if err := validate.Required("datacenter"+"."+"datacenter", "body", o.Datacenter); err != nil {
		return err
	}

	if o.Datacenter != nil {
		if err := o.Datacenter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("datacenter" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("datacenter" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post datacenters body based on the context it is used
func (o *PostDatacentersBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDatacenter(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersBody) contextValidateDatacenter(ctx context.Context, formats strfmt.Registry) error {

	if o.Datacenter != nil {
		if err := o.Datacenter.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("datacenter" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("datacenter" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostDatacentersBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDatacentersBody) UnmarshalBinary(b []byte) error {
	var res PostDatacentersBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostDatacentersCreatedBody post datacenters created body
//
// swagger:model PostDatacentersCreatedBody
type PostDatacentersCreatedBody struct {

	// datacenter
	Datacenter *models.Datacenter `json:"datacenter,omitempty"`
}

// Validate validates this post datacenters created body
func (o *PostDatacentersCreatedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDatacenter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersCreatedBody) validateDatacenter(formats strfmt.Registry) error {
	if swag.IsZero(o.Datacenter) { // not required
		return nil
	}

	if o.Datacenter != nil {
		if err := o.Datacenter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postDatacentersCreated" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postDatacentersCreated" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post datacenters created body based on the context it is used
func (o *PostDatacentersCreatedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDatacenter(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersCreatedBody) contextValidateDatacenter(ctx context.Context, formats strfmt.Registry) error {

	if o.Datacenter != nil {
		if err := o.Datacenter.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postDatacentersCreated" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postDatacentersCreated" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostDatacentersCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDatacentersCreatedBody) UnmarshalBinary(b []byte) error {
	var res PostDatacentersCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
