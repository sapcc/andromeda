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
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetGeomapsGeomapIDParams creates a new GetGeomapsGeomapIDParams object
//
// There are no default values defined in the spec.
func NewGetGeomapsGeomapIDParams() GetGeomapsGeomapIDParams {

	return GetGeomapsGeomapIDParams{}
}

// GetGeomapsGeomapIDParams contains all the bound params for the get geomaps geomap ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetGeomapsGeomapID
type GetGeomapsGeomapIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The UUID of the geomap
	  Required: true
	  In: path
	*/
	GeomapID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetGeomapsGeomapIDParams() beforehand.
func (o *GetGeomapsGeomapIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rGeomapID, rhkGeomapID, _ := route.Params.GetOK("geomap_id")
	if err := o.bindGeomapID(rGeomapID, rhkGeomapID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindGeomapID binds and validates parameter GeomapID from path.
func (o *GetGeomapsGeomapIDParams) bindGeomapID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("geomap_id", "path", "strfmt.UUID", raw)
	}
	o.GeomapID = *(value.(*strfmt.UUID))

	if err := o.validateGeomapID(formats); err != nil {
		return err
	}

	return nil
}

// validateGeomapID carries on validations for parameter GeomapID
func (o *GetGeomapsGeomapIDParams) validateGeomapID(formats strfmt.Registry) error {

	if err := validate.FormatOf("geomap_id", "path", "uuid", o.GeomapID.String(), formats); err != nil {
		return err
	}
	return nil
}
