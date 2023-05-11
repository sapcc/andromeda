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

package pools

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// PutPoolsPoolIDAcceptedCode is the HTTP code returned for type PutPoolsPoolIDAccepted
const PutPoolsPoolIDAcceptedCode int = 202

/*
PutPoolsPoolIDAccepted Updated pool.

swagger:response putPoolsPoolIdAccepted
*/
type PutPoolsPoolIDAccepted struct {

	/*
	  In: Body
	*/
	Payload *PutPoolsPoolIDAcceptedBody `json:"body,omitempty"`
}

// NewPutPoolsPoolIDAccepted creates PutPoolsPoolIDAccepted with default headers values
func NewPutPoolsPoolIDAccepted() *PutPoolsPoolIDAccepted {

	return &PutPoolsPoolIDAccepted{}
}

// WithPayload adds the payload to the put pools pool Id accepted response
func (o *PutPoolsPoolIDAccepted) WithPayload(payload *PutPoolsPoolIDAcceptedBody) *PutPoolsPoolIDAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put pools pool Id accepted response
func (o *PutPoolsPoolIDAccepted) SetPayload(payload *PutPoolsPoolIDAcceptedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutPoolsPoolIDAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutPoolsPoolIDNotFoundCode is the HTTP code returned for type PutPoolsPoolIDNotFound
const PutPoolsPoolIDNotFoundCode int = 404

/*
PutPoolsPoolIDNotFound Not Found

swagger:response putPoolsPoolIdNotFound
*/
type PutPoolsPoolIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutPoolsPoolIDNotFound creates PutPoolsPoolIDNotFound with default headers values
func NewPutPoolsPoolIDNotFound() *PutPoolsPoolIDNotFound {

	return &PutPoolsPoolIDNotFound{}
}

// WithPayload adds the payload to the put pools pool Id not found response
func (o *PutPoolsPoolIDNotFound) WithPayload(payload *models.Error) *PutPoolsPoolIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put pools pool Id not found response
func (o *PutPoolsPoolIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutPoolsPoolIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PutPoolsPoolIDDefault Unexpected Error

swagger:response putPoolsPoolIdDefault
*/
type PutPoolsPoolIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutPoolsPoolIDDefault creates PutPoolsPoolIDDefault with default headers values
func NewPutPoolsPoolIDDefault(code int) *PutPoolsPoolIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutPoolsPoolIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put pools pool ID default response
func (o *PutPoolsPoolIDDefault) WithStatusCode(code int) *PutPoolsPoolIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put pools pool ID default response
func (o *PutPoolsPoolIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put pools pool ID default response
func (o *PutPoolsPoolIDDefault) WithPayload(payload *models.Error) *PutPoolsPoolIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put pools pool ID default response
func (o *PutPoolsPoolIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutPoolsPoolIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}