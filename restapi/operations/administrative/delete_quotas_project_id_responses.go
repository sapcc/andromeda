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
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// DeleteQuotasProjectIDNoContentCode is the HTTP code returned for type DeleteQuotasProjectIDNoContent
const DeleteQuotasProjectIDNoContentCode int = 204

/*
DeleteQuotasProjectIDNoContent Resource successfully reseted.

swagger:response deleteQuotasProjectIdNoContent
*/
type DeleteQuotasProjectIDNoContent struct {
}

// NewDeleteQuotasProjectIDNoContent creates DeleteQuotasProjectIDNoContent with default headers values
func NewDeleteQuotasProjectIDNoContent() *DeleteQuotasProjectIDNoContent {

	return &DeleteQuotasProjectIDNoContent{}
}

// WriteResponse to the client
func (o *DeleteQuotasProjectIDNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteQuotasProjectIDNotFoundCode is the HTTP code returned for type DeleteQuotasProjectIDNotFound
const DeleteQuotasProjectIDNotFoundCode int = 404

/*
DeleteQuotasProjectIDNotFound Not Found

swagger:response deleteQuotasProjectIdNotFound
*/
type DeleteQuotasProjectIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteQuotasProjectIDNotFound creates DeleteQuotasProjectIDNotFound with default headers values
func NewDeleteQuotasProjectIDNotFound() *DeleteQuotasProjectIDNotFound {

	return &DeleteQuotasProjectIDNotFound{}
}

// WithPayload adds the payload to the delete quotas project Id not found response
func (o *DeleteQuotasProjectIDNotFound) WithPayload(payload *models.Error) *DeleteQuotasProjectIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete quotas project Id not found response
func (o *DeleteQuotasProjectIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQuotasProjectIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
DeleteQuotasProjectIDDefault Unexpected Error

swagger:response deleteQuotasProjectIdDefault
*/
type DeleteQuotasProjectIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteQuotasProjectIDDefault creates DeleteQuotasProjectIDDefault with default headers values
func NewDeleteQuotasProjectIDDefault(code int) *DeleteQuotasProjectIDDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteQuotasProjectIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete quotas project ID default response
func (o *DeleteQuotasProjectIDDefault) WithStatusCode(code int) *DeleteQuotasProjectIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete quotas project ID default response
func (o *DeleteQuotasProjectIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete quotas project ID default response
func (o *DeleteQuotasProjectIDDefault) WithPayload(payload *models.Error) *DeleteQuotasProjectIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete quotas project ID default response
func (o *DeleteQuotasProjectIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQuotasProjectIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
