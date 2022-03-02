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

// GetServicesOKCode is the HTTP code returned for type GetServicesOK
const GetServicesOKCode int = 200

/*GetServicesOK A JSON array of services

swagger:response getServicesOK
*/
type GetServicesOK struct {

	/*
	  In: Body
	*/
	Payload *GetServicesOKBody `json:"body,omitempty"`
}

// NewGetServicesOK creates GetServicesOK with default headers values
func NewGetServicesOK() *GetServicesOK {

	return &GetServicesOK{}
}

// WithPayload adds the payload to the get services o k response
func (o *GetServicesOK) WithPayload(payload *GetServicesOKBody) *GetServicesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services o k response
func (o *GetServicesOK) SetPayload(payload *GetServicesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetServicesDefault Unexpected Error

swagger:response getServicesDefault
*/
type GetServicesDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetServicesDefault creates GetServicesDefault with default headers values
func NewGetServicesDefault(code int) *GetServicesDefault {
	if code <= 0 {
		code = 500
	}

	return &GetServicesDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get services default response
func (o *GetServicesDefault) WithStatusCode(code int) *GetServicesDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get services default response
func (o *GetServicesDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get services default response
func (o *GetServicesDefault) WithPayload(payload *models.Error) *GetServicesDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services default response
func (o *GetServicesDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
