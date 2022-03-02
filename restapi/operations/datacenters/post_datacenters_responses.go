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
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// PostDatacentersCreatedCode is the HTTP code returned for type PostDatacentersCreated
const PostDatacentersCreatedCode int = 201

/*PostDatacentersCreated Created datacenter.

swagger:response postDatacentersCreated
*/
type PostDatacentersCreated struct {

	/*
	  In: Body
	*/
	Payload *PostDatacentersCreatedBody `json:"body,omitempty"`
}

// NewPostDatacentersCreated creates PostDatacentersCreated with default headers values
func NewPostDatacentersCreated() *PostDatacentersCreated {

	return &PostDatacentersCreated{}
}

// WithPayload adds the payload to the post datacenters created response
func (o *PostDatacentersCreated) WithPayload(payload *PostDatacentersCreatedBody) *PostDatacentersCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post datacenters created response
func (o *PostDatacentersCreated) SetPayload(payload *PostDatacentersCreatedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostDatacentersCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostDatacentersNotFoundCode is the HTTP code returned for type PostDatacentersNotFound
const PostDatacentersNotFoundCode int = 404

/*PostDatacentersNotFound Not Found

swagger:response postDatacentersNotFound
*/
type PostDatacentersNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostDatacentersNotFound creates PostDatacentersNotFound with default headers values
func NewPostDatacentersNotFound() *PostDatacentersNotFound {

	return &PostDatacentersNotFound{}
}

// WithPayload adds the payload to the post datacenters not found response
func (o *PostDatacentersNotFound) WithPayload(payload *models.Error) *PostDatacentersNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post datacenters not found response
func (o *PostDatacentersNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostDatacentersNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostDatacentersDefault Unexpected Error

swagger:response postDatacentersDefault
*/
type PostDatacentersDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostDatacentersDefault creates PostDatacentersDefault with default headers values
func NewPostDatacentersDefault(code int) *PostDatacentersDefault {
	if code <= 0 {
		code = 500
	}

	return &PostDatacentersDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post datacenters default response
func (o *PostDatacentersDefault) WithStatusCode(code int) *PostDatacentersDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post datacenters default response
func (o *PostDatacentersDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post datacenters default response
func (o *PostDatacentersDefault) WithPayload(payload *models.Error) *PostDatacentersDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post datacenters default response
func (o *PostDatacentersDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostDatacentersDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
