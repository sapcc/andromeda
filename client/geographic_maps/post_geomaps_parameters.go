// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

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
)

// NewPostGeomapsParams creates a new PostGeomapsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostGeomapsParams() *PostGeomapsParams {
	return &PostGeomapsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostGeomapsParamsWithTimeout creates a new PostGeomapsParams object
// with the ability to set a timeout on a request.
func NewPostGeomapsParamsWithTimeout(timeout time.Duration) *PostGeomapsParams {
	return &PostGeomapsParams{
		timeout: timeout,
	}
}

// NewPostGeomapsParamsWithContext creates a new PostGeomapsParams object
// with the ability to set a context for a request.
func NewPostGeomapsParamsWithContext(ctx context.Context) *PostGeomapsParams {
	return &PostGeomapsParams{
		Context: ctx,
	}
}

// NewPostGeomapsParamsWithHTTPClient creates a new PostGeomapsParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostGeomapsParamsWithHTTPClient(client *http.Client) *PostGeomapsParams {
	return &PostGeomapsParams{
		HTTPClient: client,
	}
}

/*
PostGeomapsParams contains all the parameters to send to the API endpoint

	for the post geomaps operation.

	Typically these are written to a http.Request.
*/
type PostGeomapsParams struct {

	// Geomap.
	Geomap PostGeomapsBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post geomaps params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostGeomapsParams) WithDefaults() *PostGeomapsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post geomaps params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostGeomapsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post geomaps params
func (o *PostGeomapsParams) WithTimeout(timeout time.Duration) *PostGeomapsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post geomaps params
func (o *PostGeomapsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post geomaps params
func (o *PostGeomapsParams) WithContext(ctx context.Context) *PostGeomapsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post geomaps params
func (o *PostGeomapsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post geomaps params
func (o *PostGeomapsParams) WithHTTPClient(client *http.Client) *PostGeomapsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post geomaps params
func (o *PostGeomapsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithGeomap adds the geomap to the post geomaps params
func (o *PostGeomapsParams) WithGeomap(geomap PostGeomapsBody) *PostGeomapsParams {
	o.SetGeomap(geomap)
	return o
}

// SetGeomap adds the geomap to the post geomaps params
func (o *PostGeomapsParams) SetGeomap(geomap PostGeomapsBody) {
	o.Geomap = geomap
}

// WriteToRequest writes these params to a swagger request
func (o *PostGeomapsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Geomap); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
