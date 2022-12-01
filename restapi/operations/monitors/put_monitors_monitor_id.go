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

package monitors

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

// PutMonitorsMonitorIDHandlerFunc turns a function with the right signature into a put monitors monitor ID handler
type PutMonitorsMonitorIDHandlerFunc func(PutMonitorsMonitorIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutMonitorsMonitorIDHandlerFunc) Handle(params PutMonitorsMonitorIDParams) middleware.Responder {
	return fn(params)
}

// PutMonitorsMonitorIDHandler interface for that can handle valid put monitors monitor ID params
type PutMonitorsMonitorIDHandler interface {
	Handle(PutMonitorsMonitorIDParams) middleware.Responder
}

// NewPutMonitorsMonitorID creates a new http.Handler for the put monitors monitor ID operation
func NewPutMonitorsMonitorID(ctx *middleware.Context, handler PutMonitorsMonitorIDHandler) *PutMonitorsMonitorID {
	return &PutMonitorsMonitorID{Context: ctx, Handler: handler}
}

/*
	PutMonitorsMonitorID swagger:route PUT /monitors/{monitor_id} Monitors putMonitorsMonitorId

Update a monitor
*/
type PutMonitorsMonitorID struct {
	Context *middleware.Context
	Handler PutMonitorsMonitorIDHandler
}

func (o *PutMonitorsMonitorID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutMonitorsMonitorIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutMonitorsMonitorIDAcceptedBody put monitors monitor ID accepted body
//
// swagger:model PutMonitorsMonitorIDAcceptedBody
type PutMonitorsMonitorIDAcceptedBody struct {

	// monitor
	Monitor *models.Monitor `json:"monitor,omitempty"`
}

// Validate validates this put monitors monitor ID accepted body
func (o *PutMonitorsMonitorIDAcceptedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMonitor(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMonitorsMonitorIDAcceptedBody) validateMonitor(formats strfmt.Registry) error {
	if swag.IsZero(o.Monitor) { // not required
		return nil
	}

	if o.Monitor != nil {
		if err := o.Monitor.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putMonitorsMonitorIdAccepted" + "." + "monitor")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putMonitorsMonitorIdAccepted" + "." + "monitor")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put monitors monitor ID accepted body based on the context it is used
func (o *PutMonitorsMonitorIDAcceptedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateMonitor(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMonitorsMonitorIDAcceptedBody) contextValidateMonitor(ctx context.Context, formats strfmt.Registry) error {

	if o.Monitor != nil {
		if err := o.Monitor.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putMonitorsMonitorIdAccepted" + "." + "monitor")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putMonitorsMonitorIdAccepted" + "." + "monitor")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutMonitorsMonitorIDAcceptedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutMonitorsMonitorIDAcceptedBody) UnmarshalBinary(b []byte) error {
	var res PutMonitorsMonitorIDAcceptedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutMonitorsMonitorIDBody put monitors monitor ID body
//
// swagger:model PutMonitorsMonitorIDBody
type PutMonitorsMonitorIDBody struct {

	// monitor
	// Required: true
	Monitor *models.Monitor `json:"monitor"`
}

// Validate validates this put monitors monitor ID body
func (o *PutMonitorsMonitorIDBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMonitor(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMonitorsMonitorIDBody) validateMonitor(formats strfmt.Registry) error {

	if err := validate.Required("monitor"+"."+"monitor", "body", o.Monitor); err != nil {
		return err
	}

	if o.Monitor != nil {
		if err := o.Monitor.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("monitor" + "." + "monitor")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("monitor" + "." + "monitor")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put monitors monitor ID body based on the context it is used
func (o *PutMonitorsMonitorIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateMonitor(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutMonitorsMonitorIDBody) contextValidateMonitor(ctx context.Context, formats strfmt.Registry) error {

	if o.Monitor != nil {
		if err := o.Monitor.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("monitor" + "." + "monitor")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("monitor" + "." + "monitor")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutMonitorsMonitorIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutMonitorsMonitorIDBody) UnmarshalBinary(b []byte) error {
	var res PutMonitorsMonitorIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
