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
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewPutQuotasProjectIDParams creates a new PutQuotasProjectIDParams object
//
// There are no default values defined in the spec.
func NewPutQuotasProjectIDParams() PutQuotasProjectIDParams {

	return PutQuotasProjectIDParams{}
}

// PutQuotasProjectIDParams contains all the bound params for the put quotas project ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutQuotasProjectID
type PutQuotasProjectIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The ID of the project to query.
	  Required: true
	  In: path
	*/
	ProjectID string
	/*
	  Required: true
	  In: body
	*/
	Quota PutQuotasProjectIDBody
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutQuotasProjectIDParams() beforehand.
func (o *PutQuotasProjectIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rProjectID, rhkProjectID, _ := route.Params.GetOK("project_id")
	if err := o.bindProjectID(rProjectID, rhkProjectID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body PutQuotasProjectIDBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("quota", "body", ""))
			} else {
				res = append(res, errors.NewParseError("quota", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Quota = body
			}
		}
	} else {
		res = append(res, errors.Required("quota", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindProjectID binds and validates parameter ProjectID from path.
func (o *PutQuotasProjectIDParams) bindProjectID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.ProjectID = raw

	return nil
}
