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
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetGeomapsParams creates a new GetGeomapsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetGeomapsParams() *GetGeomapsParams {
	return &GetGeomapsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetGeomapsParamsWithTimeout creates a new GetGeomapsParams object
// with the ability to set a timeout on a request.
func NewGetGeomapsParamsWithTimeout(timeout time.Duration) *GetGeomapsParams {
	return &GetGeomapsParams{
		timeout: timeout,
	}
}

// NewGetGeomapsParamsWithContext creates a new GetGeomapsParams object
// with the ability to set a context for a request.
func NewGetGeomapsParamsWithContext(ctx context.Context) *GetGeomapsParams {
	return &GetGeomapsParams{
		Context: ctx,
	}
}

// NewGetGeomapsParamsWithHTTPClient creates a new GetGeomapsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetGeomapsParamsWithHTTPClient(client *http.Client) *GetGeomapsParams {
	return &GetGeomapsParams{
		HTTPClient: client,
	}
}

/*
GetGeomapsParams contains all the parameters to send to the API endpoint

	for the get geomaps operation.

	Typically these are written to a http.Request.
*/
type GetGeomapsParams struct {

	/* Limit.

	   Sets the page size.
	*/
	Limit *int64

	/* Marker.

	   Pagination ID of the last item in the previous list.

	   Format: uuid
	*/
	Marker *strfmt.UUID

	/* PageReverse.

	   Sets the page direction.
	*/
	PageReverse *bool

	/* Sort.

	   Comma-separated list of sort keys, optinally prefix with - to reverse sort order.
	*/
	Sort *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get geomaps params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetGeomapsParams) WithDefaults() *GetGeomapsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get geomaps params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetGeomapsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get geomaps params
func (o *GetGeomapsParams) WithTimeout(timeout time.Duration) *GetGeomapsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get geomaps params
func (o *GetGeomapsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get geomaps params
func (o *GetGeomapsParams) WithContext(ctx context.Context) *GetGeomapsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get geomaps params
func (o *GetGeomapsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get geomaps params
func (o *GetGeomapsParams) WithHTTPClient(client *http.Client) *GetGeomapsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get geomaps params
func (o *GetGeomapsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithLimit adds the limit to the get geomaps params
func (o *GetGeomapsParams) WithLimit(limit *int64) *GetGeomapsParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get geomaps params
func (o *GetGeomapsParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithMarker adds the marker to the get geomaps params
func (o *GetGeomapsParams) WithMarker(marker *strfmt.UUID) *GetGeomapsParams {
	o.SetMarker(marker)
	return o
}

// SetMarker adds the marker to the get geomaps params
func (o *GetGeomapsParams) SetMarker(marker *strfmt.UUID) {
	o.Marker = marker
}

// WithPageReverse adds the pageReverse to the get geomaps params
func (o *GetGeomapsParams) WithPageReverse(pageReverse *bool) *GetGeomapsParams {
	o.SetPageReverse(pageReverse)
	return o
}

// SetPageReverse adds the pageReverse to the get geomaps params
func (o *GetGeomapsParams) SetPageReverse(pageReverse *bool) {
	o.PageReverse = pageReverse
}

// WithSort adds the sort to the get geomaps params
func (o *GetGeomapsParams) WithSort(sort *string) *GetGeomapsParams {
	o.SetSort(sort)
	return o
}

// SetSort adds the sort to the get geomaps params
func (o *GetGeomapsParams) SetSort(sort *string) {
	o.Sort = sort
}

// WriteToRequest writes these params to a swagger request
func (o *GetGeomapsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Limit != nil {

		// query param limit
		var qrLimit int64

		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {

			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}
	}

	if o.Marker != nil {

		// query param marker
		var qrMarker strfmt.UUID

		if o.Marker != nil {
			qrMarker = *o.Marker
		}
		qMarker := qrMarker.String()
		if qMarker != "" {

			if err := r.SetQueryParam("marker", qMarker); err != nil {
				return err
			}
		}
	}

	if o.PageReverse != nil {

		// query param page_reverse
		var qrPageReverse bool

		if o.PageReverse != nil {
			qrPageReverse = *o.PageReverse
		}
		qPageReverse := swag.FormatBool(qrPageReverse)
		if qPageReverse != "" {

			if err := r.SetQueryParam("page_reverse", qPageReverse); err != nil {
				return err
			}
		}
	}

	if o.Sort != nil {

		// query param sort
		var qrSort string

		if o.Sort != nil {
			qrSort = *o.Sort
		}
		qSort := qrSort
		if qSort != "" {

			if err := r.SetQueryParam("sort", qSort); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
