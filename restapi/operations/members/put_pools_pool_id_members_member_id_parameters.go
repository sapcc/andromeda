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

// NewPutPoolsPoolIDMembersMemberIDParams creates a new PutPoolsPoolIDMembersMemberIDParams object
//
// There are no default values defined in the spec.
func NewPutPoolsPoolIDMembersMemberIDParams() PutPoolsPoolIDMembersMemberIDParams {

	return PutPoolsPoolIDMembersMemberIDParams{}
}

// PutPoolsPoolIDMembersMemberIDParams contains all the bound params for the put pools pool ID members member ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutPoolsPoolIDMembersMemberID
type PutPoolsPoolIDMembersMemberIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Member PutPoolsPoolIDMembersMemberIDBody
	/*The UUID of the member
	  Required: true
	  In: path
	*/
	MemberID strfmt.UUID
	/*The UUID of the pool
	  Required: true
	  In: path
	*/
	PoolID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutPoolsPoolIDMembersMemberIDParams() beforehand.
func (o *PutPoolsPoolIDMembersMemberIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body PutPoolsPoolIDMembersMemberIDBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("member", "body", ""))
			} else {
				res = append(res, errors.NewParseError("member", "body", "", err))
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
				o.Member = body
			}
		}
	} else {
		res = append(res, errors.Required("member", "body", ""))
	}

	rMemberID, rhkMemberID, _ := route.Params.GetOK("member_id")
	if err := o.bindMemberID(rMemberID, rhkMemberID, route.Formats); err != nil {
		res = append(res, err)
	}

	rPoolID, rhkPoolID, _ := route.Params.GetOK("pool_id")
	if err := o.bindPoolID(rPoolID, rhkPoolID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindMemberID binds and validates parameter MemberID from path.
func (o *PutPoolsPoolIDMembersMemberIDParams) bindMemberID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("member_id", "path", "strfmt.UUID", raw)
	}
	o.MemberID = *(value.(*strfmt.UUID))

	if err := o.validateMemberID(formats); err != nil {
		return err
	}

	return nil
}

// validateMemberID carries on validations for parameter MemberID
func (o *PutPoolsPoolIDMembersMemberIDParams) validateMemberID(formats strfmt.Registry) error {

	if err := validate.FormatOf("member_id", "path", "uuid", o.MemberID.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindPoolID binds and validates parameter PoolID from path.
func (o *PutPoolsPoolIDMembersMemberIDParams) bindPoolID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("pool_id", "path", "strfmt.UUID", raw)
	}
	o.PoolID = *(value.(*strfmt.UUID))

	if err := o.validatePoolID(formats); err != nil {
		return err
	}

	return nil
}

// validatePoolID carries on validations for parameter PoolID
func (o *PutPoolsPoolIDMembersMemberIDParams) validatePoolID(formats strfmt.Registry) error {

	if err := validate.FormatOf("pool_id", "path", "uuid", o.PoolID.String(), formats); err != nil {
		return err
	}
	return nil
}
