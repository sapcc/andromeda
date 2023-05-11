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

package members

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

// PutMembersMemberIDHandlerFunc turns a function with the right signature into a put members member ID handler
type PutMembersMemberIDHandlerFunc func(PutMembersMemberIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutMembersMemberIDHandlerFunc) Handle(params PutMembersMemberIDParams) middleware.Responder {
	return fn(params)
}

// PutMembersMemberIDHandler interface for that can handle valid put members member ID params
type PutMembersMemberIDHandler interface {
	Handle(PutMembersMemberIDParams) middleware.Responder
}

// NewPutMembersMemberID creates a new http.Handler for the put members member ID operation
func NewPutMembersMemberID(ctx *middleware.Context, handler PutMembersMemberIDHandler) *PutMembersMemberID {
	return &PutMembersMemberID{Context: ctx, Handler: handler}
}

/*
	PutMembersMemberID swagger:route PUT /members/{member_id} Members putMembersMemberId

Update a member
*/
type PutMembersMemberID struct {
	Context *middleware.Context
	Handler PutMembersMemberIDHandler
}

func (o *PutMembersMemberID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutMembersMemberIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutMembersMemberIDAcceptedBody put members member ID accepted body
//
// swagger:model PutMembersMemberIDAcceptedBody
type PutMembersMemberIDAcceptedBody struct {

	// member
	Member *models.Member `json:"member,omitempty"`
}

// Validate validates this put members member ID accepted body
func (o *PutMembersMemberIDAcceptedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMember(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMembersMemberIDAcceptedBody) validateMember(formats strfmt.Registry) error {
	if swag.IsZero(o.Member) { // not required
		return nil
	}

	if o.Member != nil {
		if err := o.Member.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putMembersMemberIdAccepted" + "." + "member")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putMembersMemberIdAccepted" + "." + "member")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put members member ID accepted body based on the context it is used
func (o *PutMembersMemberIDAcceptedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateMember(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMembersMemberIDAcceptedBody) contextValidateMember(ctx context.Context, formats strfmt.Registry) error {

	if o.Member != nil {
		if err := o.Member.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putMembersMemberIdAccepted" + "." + "member")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putMembersMemberIdAccepted" + "." + "member")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutMembersMemberIDAcceptedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutMembersMemberIDAcceptedBody) UnmarshalBinary(b []byte) error {
	var res PutMembersMemberIDAcceptedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutMembersMemberIDBody put members member ID body
//
// swagger:model PutMembersMemberIDBody
type PutMembersMemberIDBody struct {

	// member
	// Required: true
	Member *models.Member `json:"member"`
}

// Validate validates this put members member ID body
func (o *PutMembersMemberIDBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMember(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMembersMemberIDBody) validateMember(formats strfmt.Registry) error {

	if err := validate.Required("member"+"."+"member", "body", o.Member); err != nil {
		return err
	}

	if o.Member != nil {
		if err := o.Member.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("member" + "." + "member")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("member" + "." + "member")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put members member ID body based on the context it is used
func (o *PutMembersMemberIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateMember(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMembersMemberIDBody) contextValidateMember(ctx context.Context, formats strfmt.Registry) error {

	if o.Member != nil {
		if err := o.Member.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("member" + "." + "member")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("member" + "." + "member")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutMembersMemberIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutMembersMemberIDBody) UnmarshalBinary(b []byte) error {
	var res PutMembersMemberIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}