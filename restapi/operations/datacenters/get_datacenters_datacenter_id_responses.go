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

// GetDatacentersDatacenterIDOKCode is the HTTP code returned for type GetDatacentersDatacenterIDOK
const GetDatacentersDatacenterIDOKCode int = 200

/*
GetDatacentersDatacenterIDOK Shows the details of a specific datacenter.

swagger:response getDatacentersDatacenterIdOK
*/
type GetDatacentersDatacenterIDOK struct {

	/*
	  In: Body
	*/
	Payload *GetDatacentersDatacenterIDOKBody `json:"body,omitempty"`
}

// NewGetDatacentersDatacenterIDOK creates GetDatacentersDatacenterIDOK with default headers values
func NewGetDatacentersDatacenterIDOK() *GetDatacentersDatacenterIDOK {

	return &GetDatacentersDatacenterIDOK{}
}

// WithPayload adds the payload to the get datacenters datacenter Id o k response
func (o *GetDatacentersDatacenterIDOK) WithPayload(payload *GetDatacentersDatacenterIDOKBody) *GetDatacentersDatacenterIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get datacenters datacenter Id o k response
func (o *GetDatacentersDatacenterIDOK) SetPayload(payload *GetDatacentersDatacenterIDOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDatacentersDatacenterIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetDatacentersDatacenterIDNotFoundCode is the HTTP code returned for type GetDatacentersDatacenterIDNotFound
const GetDatacentersDatacenterIDNotFoundCode int = 404

/*
GetDatacentersDatacenterIDNotFound Not Found

swagger:response getDatacentersDatacenterIdNotFound
*/
type GetDatacentersDatacenterIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetDatacentersDatacenterIDNotFound creates GetDatacentersDatacenterIDNotFound with default headers values
func NewGetDatacentersDatacenterIDNotFound() *GetDatacentersDatacenterIDNotFound {

	return &GetDatacentersDatacenterIDNotFound{}
}

// WithPayload adds the payload to the get datacenters datacenter Id not found response
func (o *GetDatacentersDatacenterIDNotFound) WithPayload(payload *models.Error) *GetDatacentersDatacenterIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get datacenters datacenter Id not found response
func (o *GetDatacentersDatacenterIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDatacentersDatacenterIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
GetDatacentersDatacenterIDDefault Unexpected Error

swagger:response getDatacentersDatacenterIdDefault
*/
type GetDatacentersDatacenterIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetDatacentersDatacenterIDDefault creates GetDatacentersDatacenterIDDefault with default headers values
func NewGetDatacentersDatacenterIDDefault(code int) *GetDatacentersDatacenterIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetDatacentersDatacenterIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get datacenters datacenter ID default response
func (o *GetDatacentersDatacenterIDDefault) WithStatusCode(code int) *GetDatacentersDatacenterIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get datacenters datacenter ID default response
func (o *GetDatacentersDatacenterIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get datacenters datacenter ID default response
func (o *GetDatacentersDatacenterIDDefault) WithPayload(payload *models.Error) *GetDatacentersDatacenterIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get datacenters datacenter ID default response
func (o *GetDatacentersDatacenterIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDatacentersDatacenterIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
