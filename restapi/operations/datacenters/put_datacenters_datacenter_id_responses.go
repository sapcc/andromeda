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

// PutDatacentersDatacenterIDAcceptedCode is the HTTP code returned for type PutDatacentersDatacenterIDAccepted
const PutDatacentersDatacenterIDAcceptedCode int = 202

/*
PutDatacentersDatacenterIDAccepted Updated datacenter.

swagger:response putDatacentersDatacenterIdAccepted
*/
type PutDatacentersDatacenterIDAccepted struct {

	/*
	  In: Body
	*/
	Payload *PutDatacentersDatacenterIDAcceptedBody `json:"body,omitempty"`
}

// NewPutDatacentersDatacenterIDAccepted creates PutDatacentersDatacenterIDAccepted with default headers values
func NewPutDatacentersDatacenterIDAccepted() *PutDatacentersDatacenterIDAccepted {

	return &PutDatacentersDatacenterIDAccepted{}
}

// WithPayload adds the payload to the put datacenters datacenter Id accepted response
func (o *PutDatacentersDatacenterIDAccepted) WithPayload(payload *PutDatacentersDatacenterIDAcceptedBody) *PutDatacentersDatacenterIDAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put datacenters datacenter Id accepted response
func (o *PutDatacentersDatacenterIDAccepted) SetPayload(payload *PutDatacentersDatacenterIDAcceptedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDatacentersDatacenterIDAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutDatacentersDatacenterIDNotFoundCode is the HTTP code returned for type PutDatacentersDatacenterIDNotFound
const PutDatacentersDatacenterIDNotFoundCode int = 404

/*
PutDatacentersDatacenterIDNotFound Not Found

swagger:response putDatacentersDatacenterIdNotFound
*/
type PutDatacentersDatacenterIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutDatacentersDatacenterIDNotFound creates PutDatacentersDatacenterIDNotFound with default headers values
func NewPutDatacentersDatacenterIDNotFound() *PutDatacentersDatacenterIDNotFound {

	return &PutDatacentersDatacenterIDNotFound{}
}

// WithPayload adds the payload to the put datacenters datacenter Id not found response
func (o *PutDatacentersDatacenterIDNotFound) WithPayload(payload *models.Error) *PutDatacentersDatacenterIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put datacenters datacenter Id not found response
func (o *PutDatacentersDatacenterIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDatacentersDatacenterIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PutDatacentersDatacenterIDDefault Unexpected Error

swagger:response putDatacentersDatacenterIdDefault
*/
type PutDatacentersDatacenterIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutDatacentersDatacenterIDDefault creates PutDatacentersDatacenterIDDefault with default headers values
func NewPutDatacentersDatacenterIDDefault(code int) *PutDatacentersDatacenterIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutDatacentersDatacenterIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put datacenters datacenter ID default response
func (o *PutDatacentersDatacenterIDDefault) WithStatusCode(code int) *PutDatacentersDatacenterIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put datacenters datacenter ID default response
func (o *PutDatacentersDatacenterIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put datacenters datacenter ID default response
func (o *PutDatacentersDatacenterIDDefault) WithPayload(payload *models.Error) *PutDatacentersDatacenterIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put datacenters datacenter ID default response
func (o *PutDatacentersDatacenterIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDatacentersDatacenterIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}