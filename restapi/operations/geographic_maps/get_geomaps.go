// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package geographic_maps

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

// GetGeomapsHandlerFunc turns a function with the right signature into a get geomaps handler
type GetGeomapsHandlerFunc func(GetGeomapsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetGeomapsHandlerFunc) Handle(params GetGeomapsParams) middleware.Responder {
	return fn(params)
}

// GetGeomapsHandler interface for that can handle valid get geomaps params
type GetGeomapsHandler interface {
	Handle(GetGeomapsParams) middleware.Responder
}

// NewGetGeomaps creates a new http.Handler for the get geomaps operation
func NewGetGeomaps(ctx *middleware.Context, handler GetGeomapsHandler) *GetGeomaps {
	return &GetGeomaps{Context: ctx, Handler: handler}
}

/*
	GetGeomaps swagger:route GET /geomaps Geographic maps getGeomaps

List geographic maps
*/
type GetGeomaps struct {
	Context *middleware.Context
	Handler GetGeomapsHandler
}

func (o *GetGeomaps) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetGeomapsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetGeomapsOKBody get geomaps o k body
//
// swagger:model GetGeomapsOKBody
type GetGeomapsOKBody struct {

	// geomaps
	Geomaps []*models.Geomap `json:"geomaps"`

	// links
	Links []*models.Link `json:"links,omitempty"`
}

// Validate validates this get geomaps o k body
func (o *GetGeomapsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGeomaps(formats); err != nil {
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

func (o *GetGeomapsOKBody) validateGeomaps(formats strfmt.Registry) error {
	if swag.IsZero(o.Geomaps) { // not required
		return nil
	}

	for i := 0; i < len(o.Geomaps); i++ {
		if swag.IsZero(o.Geomaps[i]) { // not required
			continue
		}

		if o.Geomaps[i] != nil {
			if err := o.Geomaps[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetGeomapsOKBody) validateLinks(formats strfmt.Registry) error {
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
					return ve.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get geomaps o k body based on the context it is used
func (o *GetGeomapsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateGeomaps(ctx, formats); err != nil {
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

func (o *GetGeomapsOKBody) contextValidateGeomaps(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Geomaps); i++ {

		if o.Geomaps[i] != nil {
			if err := o.Geomaps[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetGeomapsOKBody) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Links); i++ {

		if o.Links[i] != nil {
			if err := o.Links[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetGeomapsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetGeomapsOKBody) UnmarshalBinary(b []byte) error {
	var res GetGeomapsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
