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
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewPostSyncParams creates a new PostSyncParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostSyncParams() *PostSyncParams {
	return &PostSyncParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostSyncParamsWithTimeout creates a new PostSyncParams object
// with the ability to set a timeout on a request.
func NewPostSyncParamsWithTimeout(timeout time.Duration) *PostSyncParams {
	return &PostSyncParams{
		timeout: timeout,
	}
}

// NewPostSyncParamsWithContext creates a new PostSyncParams object
// with the ability to set a context for a request.
func NewPostSyncParamsWithContext(ctx context.Context) *PostSyncParams {
	return &PostSyncParams{
		Context: ctx,
	}
}

// NewPostSyncParamsWithHTTPClient creates a new PostSyncParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostSyncParamsWithHTTPClient(client *http.Client) *PostSyncParams {
	return &PostSyncParams{
		HTTPClient: client,
	}
}

/* PostSyncParams contains all the parameters to send to the API endpoint
   for the post sync operation.

   Typically these are written to a http.Request.
*/
type PostSyncParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post sync params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostSyncParams) WithDefaults() *PostSyncParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post sync params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostSyncParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post sync params
func (o *PostSyncParams) WithTimeout(timeout time.Duration) *PostSyncParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post sync params
func (o *PostSyncParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post sync params
func (o *PostSyncParams) WithContext(ctx context.Context) *PostSyncParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post sync params
func (o *PostSyncParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post sync params
func (o *PostSyncParams) WithHTTPClient(client *http.Client) *PostSyncParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post sync params
func (o *PostSyncParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *PostSyncParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
