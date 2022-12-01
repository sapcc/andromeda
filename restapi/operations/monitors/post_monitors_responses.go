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

package monitors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// PostMonitorsCreatedCode is the HTTP code returned for type PostMonitorsCreated
const PostMonitorsCreatedCode int = 201

/*
PostMonitorsCreated Created monitor.

swagger:response postMonitorsCreated
*/
type PostMonitorsCreated struct {

	/*
	  In: Body
	*/
	Payload *PostMonitorsCreatedBody `json:"body,omitempty"`
}

// NewPostMonitorsCreated creates PostMonitorsCreated with default headers values
func NewPostMonitorsCreated() *PostMonitorsCreated {

	return &PostMonitorsCreated{}
}

// WithPayload adds the payload to the post monitors created response
func (o *PostMonitorsCreated) WithPayload(payload *PostMonitorsCreatedBody) *PostMonitorsCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post monitors created response
func (o *PostMonitorsCreated) SetPayload(payload *PostMonitorsCreatedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMonitorsCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostMonitorsBadRequestCode is the HTTP code returned for type PostMonitorsBadRequest
const PostMonitorsBadRequestCode int = 400

/*
PostMonitorsBadRequest Bad request

swagger:response postMonitorsBadRequest
*/
type PostMonitorsBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMonitorsBadRequest creates PostMonitorsBadRequest with default headers values
func NewPostMonitorsBadRequest() *PostMonitorsBadRequest {

	return &PostMonitorsBadRequest{}
}

// WithPayload adds the payload to the post monitors bad request response
func (o *PostMonitorsBadRequest) WithPayload(payload *models.Error) *PostMonitorsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post monitors bad request response
func (o *PostMonitorsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMonitorsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostMonitorsNotFoundCode is the HTTP code returned for type PostMonitorsNotFound
const PostMonitorsNotFoundCode int = 404

/*
PostMonitorsNotFound Not Found

swagger:response postMonitorsNotFound
*/
type PostMonitorsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMonitorsNotFound creates PostMonitorsNotFound with default headers values
func NewPostMonitorsNotFound() *PostMonitorsNotFound {

	return &PostMonitorsNotFound{}
}

// WithPayload adds the payload to the post monitors not found response
func (o *PostMonitorsNotFound) WithPayload(payload *models.Error) *PostMonitorsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post monitors not found response
func (o *PostMonitorsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMonitorsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PostMonitorsDefault Unexpected Error

swagger:response postMonitorsDefault
*/
type PostMonitorsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMonitorsDefault creates PostMonitorsDefault with default headers values
func NewPostMonitorsDefault(code int) *PostMonitorsDefault {
	if code <= 0 {
		code = 500
	}

	return &PostMonitorsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post monitors default response
func (o *PostMonitorsDefault) WithStatusCode(code int) *PostMonitorsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post monitors default response
func (o *PostMonitorsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post monitors default response
func (o *PostMonitorsDefault) WithPayload(payload *models.Error) *PostMonitorsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post monitors default response
func (o *PostMonitorsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMonitorsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
