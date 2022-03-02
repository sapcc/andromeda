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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new administrative API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for administrative API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DeleteQuotasProjectID(params *DeleteQuotasProjectIDParams, opts ...ClientOption) (*DeleteQuotasProjectIDNoContent, error)

	GetQuotas(params *GetQuotasParams, opts ...ClientOption) (*GetQuotasOK, error)

	GetQuotasDefaults(params *GetQuotasDefaultsParams, opts ...ClientOption) (*GetQuotasDefaultsOK, error)

	GetQuotasProjectID(params *GetQuotasProjectIDParams, opts ...ClientOption) (*GetQuotasProjectIDOK, error)

	GetServices(params *GetServicesParams, opts ...ClientOption) (*GetServicesOK, error)

	PostSync(params *PostSyncParams, opts ...ClientOption) (*PostSyncAccepted, error)

	PutQuotasProjectID(params *PutQuotasProjectIDParams, opts ...ClientOption) (*PutQuotasProjectIDAccepted, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  DeleteQuotasProjectID resets all quota of a project
*/
func (a *Client) DeleteQuotasProjectID(params *DeleteQuotasProjectIDParams, opts ...ClientOption) (*DeleteQuotasProjectIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteQuotasProjectIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteQuotasProjectID",
		Method:             "DELETE",
		PathPattern:        "/quotas/{project_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteQuotasProjectIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteQuotasProjectIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DeleteQuotasProjectIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetQuotas lists quotas
*/
func (a *Client) GetQuotas(params *GetQuotasParams, opts ...ClientOption) (*GetQuotasOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetQuotasParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetQuotas",
		Method:             "GET",
		PathPattern:        "/quotas",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetQuotasReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetQuotasOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetQuotasDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetQuotasDefaults shows quota defaults
*/
func (a *Client) GetQuotasDefaults(params *GetQuotasDefaultsParams, opts ...ClientOption) (*GetQuotasDefaultsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetQuotasDefaultsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetQuotasDefaults",
		Method:             "GET",
		PathPattern:        "/quotas/defaults",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetQuotasDefaultsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetQuotasDefaultsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetQuotasDefaultsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetQuotasProjectID shows quota detail
*/
func (a *Client) GetQuotasProjectID(params *GetQuotasProjectIDParams, opts ...ClientOption) (*GetQuotasProjectIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetQuotasProjectIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetQuotasProjectID",
		Method:             "GET",
		PathPattern:        "/quotas/{project_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetQuotasProjectIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetQuotasProjectIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetQuotasProjectIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetServices lists services
*/
func (a *Client) GetServices(params *GetServicesParams, opts ...ClientOption) (*GetServicesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetServicesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetServices",
		Method:             "GET",
		PathPattern:        "/services",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetServicesReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetServicesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetServicesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  PostSync enqueues a full sync
*/
func (a *Client) PostSync(params *PostSyncParams, opts ...ClientOption) (*PostSyncAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostSyncParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostSync",
		Method:             "POST",
		PathPattern:        "/sync",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PostSyncReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostSyncAccepted)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PostSyncDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  PutQuotasProjectID updates quota
*/
func (a *Client) PutQuotasProjectID(params *PutQuotasProjectIDParams, opts ...ClientOption) (*PutQuotasProjectIDAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPutQuotasProjectIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PutQuotasProjectID",
		Method:             "PUT",
		PathPattern:        "/quotas/{project_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PutQuotasProjectIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PutQuotasProjectIDAccepted)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PutQuotasProjectIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
