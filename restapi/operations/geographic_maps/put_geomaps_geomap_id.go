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

package geographic_maps

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

// PutGeomapsGeomapIDHandlerFunc turns a function with the right signature into a put geomaps geomap ID handler
type PutGeomapsGeomapIDHandlerFunc func(PutGeomapsGeomapIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutGeomapsGeomapIDHandlerFunc) Handle(params PutGeomapsGeomapIDParams) middleware.Responder {
	return fn(params)
}

// PutGeomapsGeomapIDHandler interface for that can handle valid put geomaps geomap ID params
type PutGeomapsGeomapIDHandler interface {
	Handle(PutGeomapsGeomapIDParams) middleware.Responder
}

// NewPutGeomapsGeomapID creates a new http.Handler for the put geomaps geomap ID operation
func NewPutGeomapsGeomapID(ctx *middleware.Context, handler PutGeomapsGeomapIDHandler) *PutGeomapsGeomapID {
	return &PutGeomapsGeomapID{Context: ctx, Handler: handler}
}

/*
	PutGeomapsGeomapID swagger:route PUT /geomaps/{geomap_id} Geographic maps putGeomapsGeomapId

Update a geographic map
*/
type PutGeomapsGeomapID struct {
	Context *middleware.Context
	Handler PutGeomapsGeomapIDHandler
}

func (o *PutGeomapsGeomapID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutGeomapsGeomapIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutGeomapsGeomapIDAcceptedBody put geomaps geomap ID accepted body
//
// swagger:model PutGeomapsGeomapIDAcceptedBody
type PutGeomapsGeomapIDAcceptedBody struct {

	// geomap
	Geomap *models.Geomap `json:"geomap,omitempty"`
}

// Validate validates this put geomaps geomap ID accepted body
func (o *PutGeomapsGeomapIDAcceptedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGeomap(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutGeomapsGeomapIDAcceptedBody) validateGeomap(formats strfmt.Registry) error {
	if swag.IsZero(o.Geomap) { // not required
		return nil
	}

	if o.Geomap != nil {
		if err := o.Geomap.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putGeomapsGeomapIdAccepted" + "." + "geomap")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putGeomapsGeomapIdAccepted" + "." + "geomap")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put geomaps geomap ID accepted body based on the context it is used
func (o *PutGeomapsGeomapIDAcceptedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateGeomap(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutGeomapsGeomapIDAcceptedBody) contextValidateGeomap(ctx context.Context, formats strfmt.Registry) error {

	if o.Geomap != nil {
		if err := o.Geomap.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putGeomapsGeomapIdAccepted" + "." + "geomap")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putGeomapsGeomapIdAccepted" + "." + "geomap")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutGeomapsGeomapIDAcceptedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutGeomapsGeomapIDAcceptedBody) UnmarshalBinary(b []byte) error {
	var res PutGeomapsGeomapIDAcceptedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutGeomapsGeomapIDBody put geomaps geomap ID body
//
// swagger:model PutGeomapsGeomapIDBody
type PutGeomapsGeomapIDBody struct {

	// geomap
	// Required: true
	Geomap *models.Geomap `json:"geomap"`
}

// Validate validates this put geomaps geomap ID body
func (o *PutGeomapsGeomapIDBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGeomap(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutGeomapsGeomapIDBody) validateGeomap(formats strfmt.Registry) error {

	if err := validate.Required("geomap"+"."+"geomap", "body", o.Geomap); err != nil {
		return err
	}

	if o.Geomap != nil {
		if err := o.Geomap.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("geomap" + "." + "geomap")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("geomap" + "." + "geomap")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put geomaps geomap ID body based on the context it is used
func (o *PutGeomapsGeomapIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateGeomap(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutGeomapsGeomapIDBody) contextValidateGeomap(ctx context.Context, formats strfmt.Registry) error {

	if o.Geomap != nil {
		if err := o.Geomap.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("geomap" + "." + "geomap")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("geomap" + "." + "geomap")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutGeomapsGeomapIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutGeomapsGeomapIDBody) UnmarshalBinary(b []byte) error {
	var res PutGeomapsGeomapIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
