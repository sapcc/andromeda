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

package geographic_maps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sapcc/andromeda/models"
)

// DeleteGeomapsGeomapIDNoContentCode is the HTTP code returned for type DeleteGeomapsGeomapIDNoContent
const DeleteGeomapsGeomapIDNoContentCode int = 204

/*
DeleteGeomapsGeomapIDNoContent Resource successfully deleted.

swagger:response deleteGeomapsGeomapIdNoContent
*/
type DeleteGeomapsGeomapIDNoContent struct {
}

// NewDeleteGeomapsGeomapIDNoContent creates DeleteGeomapsGeomapIDNoContent with default headers values
func NewDeleteGeomapsGeomapIDNoContent() *DeleteGeomapsGeomapIDNoContent {

	return &DeleteGeomapsGeomapIDNoContent{}
}

// WriteResponse to the client
func (o *DeleteGeomapsGeomapIDNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteGeomapsGeomapIDNotFoundCode is the HTTP code returned for type DeleteGeomapsGeomapIDNotFound
const DeleteGeomapsGeomapIDNotFoundCode int = 404

/*
DeleteGeomapsGeomapIDNotFound Not Found

swagger:response deleteGeomapsGeomapIdNotFound
*/
type DeleteGeomapsGeomapIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteGeomapsGeomapIDNotFound creates DeleteGeomapsGeomapIDNotFound with default headers values
func NewDeleteGeomapsGeomapIDNotFound() *DeleteGeomapsGeomapIDNotFound {

	return &DeleteGeomapsGeomapIDNotFound{}
}

// WithPayload adds the payload to the delete geomaps geomap Id not found response
func (o *DeleteGeomapsGeomapIDNotFound) WithPayload(payload *models.Error) *DeleteGeomapsGeomapIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete geomaps geomap Id not found response
func (o *DeleteGeomapsGeomapIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteGeomapsGeomapIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
DeleteGeomapsGeomapIDDefault Unexpected Error

swagger:response deleteGeomapsGeomapIdDefault
*/
type DeleteGeomapsGeomapIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteGeomapsGeomapIDDefault creates DeleteGeomapsGeomapIDDefault with default headers values
func NewDeleteGeomapsGeomapIDDefault(code int) *DeleteGeomapsGeomapIDDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteGeomapsGeomapIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete geomaps geomap ID default response
func (o *DeleteGeomapsGeomapIDDefault) WithStatusCode(code int) *DeleteGeomapsGeomapIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete geomaps geomap ID default response
func (o *DeleteGeomapsGeomapIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete geomaps geomap ID default response
func (o *DeleteGeomapsGeomapIDDefault) WithPayload(payload *models.Error) *DeleteGeomapsGeomapIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete geomaps geomap ID default response
func (o *DeleteGeomapsGeomapIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteGeomapsGeomapIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}