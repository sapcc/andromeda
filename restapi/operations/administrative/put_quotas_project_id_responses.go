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

// PutQuotasProjectIDAcceptedCode is the HTTP code returned for type PutQuotasProjectIDAccepted
const PutQuotasProjectIDAcceptedCode int = 202

/*
PutQuotasProjectIDAccepted Updated quota for a project.

swagger:response putQuotasProjectIdAccepted
*/
type PutQuotasProjectIDAccepted struct {

	/*
	  In: Body
	*/
	Payload *PutQuotasProjectIDAcceptedBody `json:"body,omitempty"`
}

// NewPutQuotasProjectIDAccepted creates PutQuotasProjectIDAccepted with default headers values
func NewPutQuotasProjectIDAccepted() *PutQuotasProjectIDAccepted {

	return &PutQuotasProjectIDAccepted{}
}

// WithPayload adds the payload to the put quotas project Id accepted response
func (o *PutQuotasProjectIDAccepted) WithPayload(payload *PutQuotasProjectIDAcceptedBody) *PutQuotasProjectIDAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put quotas project Id accepted response
func (o *PutQuotasProjectIDAccepted) SetPayload(payload *PutQuotasProjectIDAcceptedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutQuotasProjectIDAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PutQuotasProjectIDDefault Unexpected Error

swagger:response putQuotasProjectIdDefault
*/
type PutQuotasProjectIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutQuotasProjectIDDefault creates PutQuotasProjectIDDefault with default headers values
func NewPutQuotasProjectIDDefault(code int) *PutQuotasProjectIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutQuotasProjectIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put quotas project ID default response
func (o *PutQuotasProjectIDDefault) WithStatusCode(code int) *PutQuotasProjectIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put quotas project ID default response
func (o *PutQuotasProjectIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put quotas project ID default response
func (o *PutQuotasProjectIDDefault) WithPayload(payload *models.Error) *PutQuotasProjectIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put quotas project ID default response
func (o *PutQuotasProjectIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutQuotasProjectIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
