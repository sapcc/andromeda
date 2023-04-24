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
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetMembersParams creates a new GetMembersParams object
//
// There are no default values defined in the spec.
func NewGetMembersParams() GetMembersParams {

	return GetMembersParams{}
}

// GetMembersParams contains all the bound params for the get members operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetMembers
type GetMembersParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Sets the page size.
	  In: query
	*/
	Limit *int64
	/*Pagination ID of the last item in the previous list.
	  In: query
	*/
	Marker *strfmt.UUID
	/*Filter for resources not having tags, multiple not-tags are considered as logical AND.
	Should be provided in a comma separated list.

	  In: query
	*/
	NotTags []string
	/*Filter for resources not having tags, multiple tags are considered as logical OR.
	Should be provided in a comma separated list.

	  In: query
	*/
	NotTagsAny []string
	/*Sets the page direction.
	  In: query
	*/
	PageReverse *bool
	/*Pool ID of the members to fetch
	  In: query
	*/
	PoolID *strfmt.UUID
	/*Comma-separated list of sort keys, optinally prefix with - to reverse sort order.
	  In: query
	*/
	Sort *string
	/*Filter for tags, multiple tags are considered as logical AND.
	Should be provided in a comma separated list.

	  In: query
	*/
	Tags []string
	/*Filter for tags, multiple tags are considered as logical OR.
	Should be provided in a comma separated list.

	  In: query
	*/
	TagsAny []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetMembersParams() beforehand.
func (o *GetMembersParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qMarker, qhkMarker, _ := qs.GetOK("marker")
	if err := o.bindMarker(qMarker, qhkMarker, route.Formats); err != nil {
		res = append(res, err)
	}

	qNotTags, qhkNotTags, _ := qs.GetOK("not-tags")
	if err := o.bindNotTags(qNotTags, qhkNotTags, route.Formats); err != nil {
		res = append(res, err)
	}

	qNotTagsAny, qhkNotTagsAny, _ := qs.GetOK("not-tags-any")
	if err := o.bindNotTagsAny(qNotTagsAny, qhkNotTagsAny, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageReverse, qhkPageReverse, _ := qs.GetOK("page_reverse")
	if err := o.bindPageReverse(qPageReverse, qhkPageReverse, route.Formats); err != nil {
		res = append(res, err)
	}

	qPoolID, qhkPoolID, _ := qs.GetOK("pool_id")
	if err := o.bindPoolID(qPoolID, qhkPoolID, route.Formats); err != nil {
		res = append(res, err)
	}

	qSort, qhkSort, _ := qs.GetOK("sort")
	if err := o.bindSort(qSort, qhkSort, route.Formats); err != nil {
		res = append(res, err)
	}

	qTags, qhkTags, _ := qs.GetOK("tags")
	if err := o.bindTags(qTags, qhkTags, route.Formats); err != nil {
		res = append(res, err)
	}

	qTagsAny, qhkTagsAny, _ := qs.GetOK("tags-any")
	if err := o.bindTagsAny(qTagsAny, qhkTagsAny, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *GetMembersParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetMembersParams) bindMarker(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetMembersParams) validateMarker(formats strfmt.Registry) error {

	if err := validate.FormatOf("marker", "query", "uuid", o.Marker.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindNotTags binds and validates array parameter NotTags from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetMembersParams) bindNotTags(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvNotTags string
	if len(rawData) > 0 {
		qvNotTags = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	notTagsIC := swag.SplitByFormat(qvNotTags, "")
	if len(notTagsIC) == 0 {
		return nil
	}

	var notTagsIR []string
	for _, notTagsIV := range notTagsIC {
		notTagsI := notTagsIV

		notTagsIR = append(notTagsIR, notTagsI)
	}

	o.NotTags = notTagsIR

	return nil
}

// bindNotTagsAny binds and validates array parameter NotTagsAny from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetMembersParams) bindNotTagsAny(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvNotTagsAny string
	if len(rawData) > 0 {
		qvNotTagsAny = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	notTagsAnyIC := swag.SplitByFormat(qvNotTagsAny, "")
	if len(notTagsAnyIC) == 0 {
		return nil
	}

	var notTagsAnyIR []string
	for _, notTagsAnyIV := range notTagsAnyIC {
		notTagsAnyI := notTagsAnyIV

		notTagsAnyIR = append(notTagsAnyIR, notTagsAnyI)
	}

	o.NotTagsAny = notTagsAnyIR

	return nil
}

// bindPageReverse binds and validates parameter PageReverse from query.
func (o *GetMembersParams) bindPageReverse(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindPoolID binds and validates parameter PoolID from query.
func (o *GetMembersParams) bindPoolID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
		return errors.InvalidType("pool_id", "query", "strfmt.UUID", raw)
	}
	o.PoolID = (value.(*strfmt.UUID))

	if err := o.validatePoolID(formats); err != nil {
		return err
	}

	return nil
}

// validatePoolID carries on validations for parameter PoolID
func (o *GetMembersParams) validatePoolID(formats strfmt.Registry) error {

	if err := validate.FormatOf("pool_id", "query", "uuid", o.PoolID.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindSort binds and validates parameter Sort from query.
func (o *GetMembersParams) bindSort(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindTags binds and validates array parameter Tags from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetMembersParams) bindTags(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvTags string
	if len(rawData) > 0 {
		qvTags = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	tagsIC := swag.SplitByFormat(qvTags, "")
	if len(tagsIC) == 0 {
		return nil
	}

	var tagsIR []string
	for _, tagsIV := range tagsIC {
		tagsI := tagsIV

		tagsIR = append(tagsIR, tagsI)
	}

	o.Tags = tagsIR

	return nil
}

// bindTagsAny binds and validates array parameter TagsAny from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetMembersParams) bindTagsAny(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvTagsAny string
	if len(rawData) > 0 {
		qvTagsAny = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	tagsAnyIC := swag.SplitByFormat(qvTagsAny, "")
	if len(tagsAnyIC) == 0 {
		return nil
	}

	var tagsAnyIR []string
	for _, tagsAnyIV := range tagsAnyIC {
		tagsAnyI := tagsAnyIV

		tagsAnyIR = append(tagsAnyIR, tagsAnyI)
	}

	o.TagsAny = tagsAnyIR

	return nil
}
