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

package datacenters

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
)

// NewPostDatacentersParams creates a new PostDatacentersParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostDatacentersParams() *PostDatacentersParams {
	return &PostDatacentersParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostDatacentersParamsWithTimeout creates a new PostDatacentersParams object
// with the ability to set a timeout on a request.
func NewPostDatacentersParamsWithTimeout(timeout time.Duration) *PostDatacentersParams {
	return &PostDatacentersParams{
		timeout: timeout,
	}
}

// NewPostDatacentersParamsWithContext creates a new PostDatacentersParams object
// with the ability to set a context for a request.
func NewPostDatacentersParamsWithContext(ctx context.Context) *PostDatacentersParams {
	return &PostDatacentersParams{
		Context: ctx,
	}
}

// NewPostDatacentersParamsWithHTTPClient creates a new PostDatacentersParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostDatacentersParamsWithHTTPClient(client *http.Client) *PostDatacentersParams {
	return &PostDatacentersParams{
		HTTPClient: client,
	}
}

/*
PostDatacentersParams contains all the parameters to send to the API endpoint

	for the post datacenters operation.

	Typically these are written to a http.Request.
*/
type PostDatacentersParams struct {

	// Datacenter.
	Datacenter PostDatacentersBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post datacenters params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostDatacentersParams) WithDefaults() *PostDatacentersParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post datacenters params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostDatacentersParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post datacenters params
func (o *PostDatacentersParams) WithTimeout(timeout time.Duration) *PostDatacentersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post datacenters params
func (o *PostDatacentersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post datacenters params
func (o *PostDatacentersParams) WithContext(ctx context.Context) *PostDatacentersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post datacenters params
func (o *PostDatacentersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post datacenters params
func (o *PostDatacentersParams) WithHTTPClient(client *http.Client) *PostDatacentersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post datacenters params
func (o *PostDatacentersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDatacenter adds the datacenter to the post datacenters params
func (o *PostDatacentersParams) WithDatacenter(datacenter PostDatacentersBody) *PostDatacentersParams {
	o.SetDatacenter(datacenter)
	return o
}

// SetDatacenter adds the datacenter to the post datacenters params
func (o *PostDatacentersParams) SetDatacenter(datacenter PostDatacentersBody) {
	o.Datacenter = datacenter
}

// WriteToRequest writes these params to a swagger request
func (o *PostDatacentersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Datacenter); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
