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

// NewDeleteQuotasProjectIDParams creates a new DeleteQuotasProjectIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteQuotasProjectIDParams() *DeleteQuotasProjectIDParams {
	return &DeleteQuotasProjectIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteQuotasProjectIDParamsWithTimeout creates a new DeleteQuotasProjectIDParams object
// with the ability to set a timeout on a request.
func NewDeleteQuotasProjectIDParamsWithTimeout(timeout time.Duration) *DeleteQuotasProjectIDParams {
	return &DeleteQuotasProjectIDParams{
		timeout: timeout,
	}
}

// NewDeleteQuotasProjectIDParamsWithContext creates a new DeleteQuotasProjectIDParams object
// with the ability to set a context for a request.
func NewDeleteQuotasProjectIDParamsWithContext(ctx context.Context) *DeleteQuotasProjectIDParams {
	return &DeleteQuotasProjectIDParams{
		Context: ctx,
	}
}

// NewDeleteQuotasProjectIDParamsWithHTTPClient creates a new DeleteQuotasProjectIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteQuotasProjectIDParamsWithHTTPClient(client *http.Client) *DeleteQuotasProjectIDParams {
	return &DeleteQuotasProjectIDParams{
		HTTPClient: client,
	}
}

/*
DeleteQuotasProjectIDParams contains all the parameters to send to the API endpoint

	for the delete quotas project ID operation.

	Typically these are written to a http.Request.
*/
type DeleteQuotasProjectIDParams struct {

	/* ProjectID.

	   The ID of the project to query.
	*/
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete quotas project ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteQuotasProjectIDParams) WithDefaults() *DeleteQuotasProjectIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete quotas project ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteQuotasProjectIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) WithTimeout(timeout time.Duration) *DeleteQuotasProjectIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) WithContext(ctx context.Context) *DeleteQuotasProjectIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) WithHTTPClient(client *http.Client) *DeleteQuotasProjectIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProjectID adds the projectID to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) WithProjectID(projectID string) *DeleteQuotasProjectIDParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the delete quotas project ID params
func (o *DeleteQuotasProjectIDParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteQuotasProjectIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
