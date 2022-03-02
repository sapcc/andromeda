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

package domains

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// GetDomainsOKCode is the HTTP code returned for type GetDomainsOK
const GetDomainsOKCode int = 200

/*GetDomainsOK A JSON array of domains

swagger:response getDomainsOK
*/
type GetDomainsOK struct {

	/*
	  In: Body
	*/
	Payload *GetDomainsOKBody `json:"body,omitempty"`
}

// NewGetDomainsOK creates GetDomainsOK with default headers values
func NewGetDomainsOK() *GetDomainsOK {

	return &GetDomainsOK{}
}

// WithPayload adds the payload to the get domains o k response
func (o *GetDomainsOK) WithPayload(payload *GetDomainsOKBody) *GetDomainsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get domains o k response
func (o *GetDomainsOK) SetPayload(payload *GetDomainsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDomainsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetDomainsBadRequestCode is the HTTP code returned for type GetDomainsBadRequest
const GetDomainsBadRequestCode int = 400

/*GetDomainsBadRequest Bad request

swagger:response getDomainsBadRequest
*/
type GetDomainsBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetDomainsBadRequest creates GetDomainsBadRequest with default headers values
func NewGetDomainsBadRequest() *GetDomainsBadRequest {

	return &GetDomainsBadRequest{}
}

// WithPayload adds the payload to the get domains bad request response
func (o *GetDomainsBadRequest) WithPayload(payload *models.Error) *GetDomainsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get domains bad request response
func (o *GetDomainsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDomainsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetDomainsDefault Unexpected Error

swagger:response getDomainsDefault
*/
type GetDomainsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetDomainsDefault creates GetDomainsDefault with default headers values
func NewGetDomainsDefault(code int) *GetDomainsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetDomainsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get domains default response
func (o *GetDomainsDefault) WithStatusCode(code int) *GetDomainsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get domains default response
func (o *GetDomainsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get domains default response
func (o *GetDomainsDefault) WithPayload(payload *models.Error) *GetDomainsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get domains default response
func (o *GetDomainsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDomainsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
