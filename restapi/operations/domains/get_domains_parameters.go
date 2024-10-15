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

package domains

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetDomainsParams creates a new GetDomainsParams object
//
// There are no default values defined in the spec.
func NewGetDomainsParams() GetDomainsParams {

	return GetDomainsParams{}
}

// GetDomainsParams contains all the bound params for the get domains operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetDomains
type GetDomainsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Filter domains by domain ID
	  In: query
	*/
	DomainID *strfmt.UUID
	/*Sets the page size.
	  In: query
	*/
	Limit *int64
	/*Pagination ID of the last item in the previous list.
	  In: query
	*/
	Marker *strfmt.UUID
	/*Sets the page direction.
	  In: query
	*/
	PageReverse *bool
	/*Comma-separated list of sort keys, optinally prefix with - to reverse sort order.
	  In: query
	*/
	Sort *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetDomainsParams() beforehand.
func (o *GetDomainsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDomainID, qhkDomainID, _ := qs.GetOK("domain_id")
	if err := o.bindDomainID(qDomainID, qhkDomainID, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qMarker, qhkMarker, _ := qs.GetOK("marker")
	if err := o.bindMarker(qMarker, qhkMarker, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageReverse, qhkPageReverse, _ := qs.GetOK("page_reverse")
	if err := o.bindPageReverse(qPageReverse, qhkPageReverse, route.Formats); err != nil {
		res = append(res, err)
	}

	qSort, qhkSort, _ := qs.GetOK("sort")
	if err := o.bindSort(qSort, qhkSort, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDomainID binds and validates parameter DomainID from query.
func (o *GetDomainsParams) bindDomainID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("domain_id", "query", "strfmt.UUID", raw)
	}
	o.DomainID = (value.(*strfmt.UUID))

	if err := o.validateDomainID(formats); err != nil {
		return err
	}

	return nil
}

// validateDomainID carries on validations for parameter DomainID
func (o *GetDomainsParams) validateDomainID(formats strfmt.Registry) error {

	if err := validate.FormatOf("domain_id", "query", "uuid", o.DomainID.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *GetDomainsParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = &value

	return nil
}

// bindMarker binds and validates parameter Marker from query.
func (o *GetDomainsParams) bindMarker(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("marker", "query", "strfmt.UUID", raw)
	}
	o.Marker = (value.(*strfmt.UUID))

	if err := o.validateMarker(formats); err != nil {
		return err
	}

	return nil
}

// validateMarker carries on validations for parameter Marker
func (o *GetDomainsParams) validateMarker(formats strfmt.Registry) error {

	if err := validate.FormatOf("marker", "query", "uuid", o.Marker.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindPageReverse binds and validates parameter PageReverse from query.
func (o *GetDomainsParams) bindPageReverse(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("page_reverse", "query", "bool", raw)
	}
	o.PageReverse = &value

	return nil
}

// bindSort binds and validates parameter Sort from query.
func (o *GetDomainsParams) bindSort(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Sort = &raw

	return nil
}
