// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package domains

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewDeleteDomainsDomainIDParams creates a new DeleteDomainsDomainIDParams object
//
// There are no default values defined in the spec.
func NewDeleteDomainsDomainIDParams() DeleteDomainsDomainIDParams {

	return DeleteDomainsDomainIDParams{}
}

// DeleteDomainsDomainIDParams contains all the bound params for the delete domains domain ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteDomainsDomainID
type DeleteDomainsDomainIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The UUID of the domain
	  Required: true
	  In: path
	*/
	DomainID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteDomainsDomainIDParams() beforehand.
func (o *DeleteDomainsDomainIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rDomainID, rhkDomainID, _ := route.Params.GetOK("domain_id")
	if err := o.bindDomainID(rDomainID, rhkDomainID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDomainID binds and validates parameter DomainID from path.
func (o *DeleteDomainsDomainIDParams) bindDomainID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("domain_id", "path", "strfmt.UUID", raw)
	}
	o.DomainID = *(value.(*strfmt.UUID))

	if err := o.validateDomainID(formats); err != nil {
		return err
	}

	return nil
}

// validateDomainID carries on validations for parameter DomainID
func (o *DeleteDomainsDomainIDParams) validateDomainID(formats strfmt.Registry) error {

	if err := validate.FormatOf("domain_id", "path", "uuid", o.DomainID.String(), formats); err != nil {
		return err
	}
	return nil
}
