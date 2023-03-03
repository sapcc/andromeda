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

package members

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// PutMembersMemberIDAcceptedCode is the HTTP code returned for type PutMembersMemberIDAccepted
const PutMembersMemberIDAcceptedCode int = 202

/*
PutMembersMemberIDAccepted Updated member.

swagger:response putMembersMemberIdAccepted
*/
type PutMembersMemberIDAccepted struct {

	/*
	  In: Body
	*/
	Payload *PutMembersMemberIDAcceptedBody `json:"body,omitempty"`
}

// NewPutMembersMemberIDAccepted creates PutMembersMemberIDAccepted with default headers values
func NewPutMembersMemberIDAccepted() *PutMembersMemberIDAccepted {

	return &PutMembersMemberIDAccepted{}
}

// WithPayload adds the payload to the put members member Id accepted response
func (o *PutMembersMemberIDAccepted) WithPayload(payload *PutMembersMemberIDAcceptedBody) *PutMembersMemberIDAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put members member Id accepted response
func (o *PutMembersMemberIDAccepted) SetPayload(payload *PutMembersMemberIDAcceptedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMembersMemberIDAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutMembersMemberIDNotFoundCode is the HTTP code returned for type PutMembersMemberIDNotFound
const PutMembersMemberIDNotFoundCode int = 404

/*
PutMembersMemberIDNotFound Not Found

swagger:response putMembersMemberIdNotFound
*/
type PutMembersMemberIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutMembersMemberIDNotFound creates PutMembersMemberIDNotFound with default headers values
func NewPutMembersMemberIDNotFound() *PutMembersMemberIDNotFound {

	return &PutMembersMemberIDNotFound{}
}

// WithPayload adds the payload to the put members member Id not found response
func (o *PutMembersMemberIDNotFound) WithPayload(payload *models.Error) *PutMembersMemberIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put members member Id not found response
func (o *PutMembersMemberIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMembersMemberIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PutMembersMemberIDDefault Unexpected Error

swagger:response putMembersMemberIdDefault
*/
type PutMembersMemberIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutMembersMemberIDDefault creates PutMembersMemberIDDefault with default headers values
func NewPutMembersMemberIDDefault(code int) *PutMembersMemberIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutMembersMemberIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put members member ID default response
func (o *PutMembersMemberIDDefault) WithStatusCode(code int) *PutMembersMemberIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put members member ID default response
func (o *PutMembersMemberIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put members member ID default response
func (o *PutMembersMemberIDDefault) WithPayload(payload *models.Error) *PutMembersMemberIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put members member ID default response
func (o *PutMembersMemberIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMembersMemberIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
