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

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/sapcc/andromeda/models"
)

// PutQuotasProjectIDHandlerFunc turns a function with the right signature into a put quotas project ID handler
type PutQuotasProjectIDHandlerFunc func(PutQuotasProjectIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutQuotasProjectIDHandlerFunc) Handle(params PutQuotasProjectIDParams) middleware.Responder {
	return fn(params)
}

// PutQuotasProjectIDHandler interface for that can handle valid put quotas project ID params
type PutQuotasProjectIDHandler interface {
	Handle(PutQuotasProjectIDParams) middleware.Responder
}

// NewPutQuotasProjectID creates a new http.Handler for the put quotas project ID operation
func NewPutQuotasProjectID(ctx *middleware.Context, handler PutQuotasProjectIDHandler) *PutQuotasProjectID {
	return &PutQuotasProjectID{Context: ctx, Handler: handler}
}

/* PutQuotasProjectID swagger:route PUT /quotas/{project_id} Administrative putQuotasProjectId

Update Quota

*/
type PutQuotasProjectID struct {
	Context *middleware.Context
	Handler PutQuotasProjectIDHandler
}

func (o *PutQuotasProjectID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutQuotasProjectIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutQuotasProjectIDAcceptedBody put quotas project ID accepted body
//
// swagger:model PutQuotasProjectIDAcceptedBody
type PutQuotasProjectIDAcceptedBody struct {

	// quota
	Quota *models.Quota `json:"quota,omitempty"`
}

// Validate validates this put quotas project ID accepted body
func (o *PutQuotasProjectIDAcceptedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateQuota(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutQuotasProjectIDAcceptedBody) validateQuota(formats strfmt.Registry) error {
	if swag.IsZero(o.Quota) { // not required
		return nil
	}

	if o.Quota != nil {
		if err := o.Quota.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putQuotasProjectIdAccepted" + "." + "quota")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putQuotasProjectIdAccepted" + "." + "quota")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put quotas project ID accepted body based on the context it is used
func (o *PutQuotasProjectIDAcceptedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateQuota(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutQuotasProjectIDAcceptedBody) contextValidateQuota(ctx context.Context, formats strfmt.Registry) error {

	if o.Quota != nil {
		if err := o.Quota.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putQuotasProjectIdAccepted" + "." + "quota")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putQuotasProjectIdAccepted" + "." + "quota")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutQuotasProjectIDAcceptedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutQuotasProjectIDAcceptedBody) UnmarshalBinary(b []byte) error {
	var res PutQuotasProjectIDAcceptedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutQuotasProjectIDBody put quotas project ID body
//
// swagger:model PutQuotasProjectIDBody
type PutQuotasProjectIDBody struct {

	// quota
	// Required: true
	Quota *models.Quota `json:"quota"`
}

// Validate validates this put quotas project ID body
func (o *PutQuotasProjectIDBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateQuota(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutQuotasProjectIDBody) validateQuota(formats strfmt.Registry) error {

	if err := validate.Required("quota"+"."+"quota", "body", o.Quota); err != nil {
		return err
	}

	if o.Quota != nil {
		if err := o.Quota.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("quota" + "." + "quota")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("quota" + "." + "quota")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put quotas project ID body based on the context it is used
func (o *PutQuotasProjectIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateQuota(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutQuotasProjectIDBody) contextValidateQuota(ctx context.Context, formats strfmt.Registry) error {

	if o.Quota != nil {
		if err := o.Quota.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("quota" + "." + "quota")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("quota" + "." + "quota")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutQuotasProjectIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutQuotasProjectIDBody) UnmarshalBinary(b []byte) error {
	var res PutQuotasProjectIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
