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
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetDomainsParams creates a new GetDomainsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetDomainsParams() *GetDomainsParams {
	return &GetDomainsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetDomainsParamsWithTimeout creates a new GetDomainsParams object
// with the ability to set a timeout on a request.
func NewGetDomainsParamsWithTimeout(timeout time.Duration) *GetDomainsParams {
	return &GetDomainsParams{
		timeout: timeout,
	}
}

// NewGetDomainsParamsWithContext creates a new GetDomainsParams object
// with the ability to set a context for a request.
func NewGetDomainsParamsWithContext(ctx context.Context) *GetDomainsParams {
	return &GetDomainsParams{
		Context: ctx,
	}
}

// NewGetDomainsParamsWithHTTPClient creates a new GetDomainsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetDomainsParamsWithHTTPClient(client *http.Client) *GetDomainsParams {
	return &GetDomainsParams{
		HTTPClient: client,
	}
}

/*
GetDomainsParams contains all the parameters to send to the API endpoint

	for the get domains operation.

	Typically these are written to a http.Request.
*/
type GetDomainsParams struct {

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

// WithDefaults hydrates default values in the get domains params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDomainsParams) WithDefaults() *GetDomainsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get domains params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDomainsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get domains params
func (o *GetDomainsParams) WithTimeout(timeout time.Duration) *GetDomainsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get domains params
func (o *GetDomainsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get domains params
func (o *GetDomainsParams) WithContext(ctx context.Context) *GetDomainsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get domains params
func (o *GetDomainsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get domains params
func (o *GetDomainsParams) WithHTTPClient(client *http.Client) *GetDomainsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get domains params
func (o *GetDomainsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithLimit adds the limit to the get domains params
func (o *GetDomainsParams) WithLimit(limit *int64) *GetDomainsParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get domains params
func (o *GetDomainsParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithMarker adds the marker to the get domains params
func (o *GetDomainsParams) WithMarker(marker *strfmt.UUID) *GetDomainsParams {
	o.SetMarker(marker)
	return o
}

// SetMarker adds the marker to the get domains params
func (o *GetDomainsParams) SetMarker(marker *strfmt.UUID) {
	o.Marker = marker
}

// WithPageReverse adds the pageReverse to the get domains params
func (o *GetDomainsParams) WithPageReverse(pageReverse *bool) *GetDomainsParams {
	o.SetPageReverse(pageReverse)
	return o
}

// SetPageReverse adds the pageReverse to the get domains params
func (o *GetDomainsParams) SetPageReverse(pageReverse *bool) {
	o.PageReverse = pageReverse
}

// WithSort adds the sort to the get domains params
func (o *GetDomainsParams) WithSort(sort *string) *GetDomainsParams {
	o.SetSort(sort)
	return o
}

// SetSort adds the sort to the get domains params
func (o *GetDomainsParams) SetSort(sort *string) {
	o.Sort = sort
}

// WriteToRequest writes these params to a swagger request
func (o *GetDomainsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
