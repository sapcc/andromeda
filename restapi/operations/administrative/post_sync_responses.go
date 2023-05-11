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

// PostSyncAcceptedCode is the HTTP code returned for type PostSyncAccepted
const PostSyncAcceptedCode int = 202

/*
PostSyncAccepted Full sync has been enqueued.

swagger:response postSyncAccepted
*/
type PostSyncAccepted struct {
}

// NewPostSyncAccepted creates PostSyncAccepted with default headers values
func NewPostSyncAccepted() *PostSyncAccepted {

	return &PostSyncAccepted{}
}

// WriteResponse to the client
func (o *PostSyncAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

/*
PostSyncDefault Unexpected Error

swagger:response postSyncDefault
*/
type PostSyncDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostSyncDefault creates PostSyncDefault with default headers values
func NewPostSyncDefault(code int) *PostSyncDefault {
	if code <= 0 {
		code = 500
	}

	return &PostSyncDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post sync default response
func (o *PostSyncDefault) WithStatusCode(code int) *PostSyncDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post sync default response
func (o *PostSyncDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post sync default response
func (o *PostSyncDefault) WithPayload(payload *models.Error) *PostSyncDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post sync default response
func (o *PostSyncDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostSyncDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}