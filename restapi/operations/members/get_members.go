// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package members

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

// GetMembersHandlerFunc turns a function with the right signature into a get members handler
type GetMembersHandlerFunc func(GetMembersParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMembersHandlerFunc) Handle(params GetMembersParams) middleware.Responder {
	return fn(params)
}

// GetMembersHandler interface for that can handle valid get members params
type GetMembersHandler interface {
	Handle(GetMembersParams) middleware.Responder
}

// NewGetMembers creates a new http.Handler for the get members operation
func NewGetMembers(ctx *middleware.Context, handler GetMembersHandler) *GetMembers {
	return &GetMembers{Context: ctx, Handler: handler}
}

/*
	GetMembers swagger:route GET /members Members getMembers

List members
*/
type GetMembers struct {
	Context *middleware.Context
	Handler GetMembersHandler
}

func (o *GetMembers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetMembersParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetMembersOKBody get members o k body
//
// swagger:model GetMembersOKBody
type GetMembersOKBody struct {

	// links
	Links []*models.Link `json:"links,omitempty"`

	// members
	Members []*models.Member `json:"members"`
}

// Validate validates this get members o k body
func (o *GetMembersOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateMembers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetMembersOKBody) validateLinks(formats strfmt.Registry) error {
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
					return ve.ValidateName("getMembersOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getMembersOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetMembersOKBody) validateMembers(formats strfmt.Registry) error {
	if swag.IsZero(o.Members) { // not required
		return nil
	}

	for i := 0; i < len(o.Members); i++ {
		if swag.IsZero(o.Members[i]) { // not required
			continue
		}

		if o.Members[i] != nil {
			if err := o.Members[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getMembersOK" + "." + "members" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getMembersOK" + "." + "members" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get members o k body based on the context it is used
func (o *GetMembersOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateMembers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetMembersOKBody) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Links); i++ {

		if o.Links[i] != nil {
			if err := o.Links[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getMembersOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getMembersOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetMembersOKBody) contextValidateMembers(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Members); i++ {

		if o.Members[i] != nil {
			if err := o.Members[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getMembersOK" + "." + "members" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getMembersOK" + "." + "members" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetMembersOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetMembersOKBody) UnmarshalBinary(b []byte) error {
	var res GetMembersOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
