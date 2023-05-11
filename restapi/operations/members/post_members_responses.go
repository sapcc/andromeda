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

// PostMembersCreatedCode is the HTTP code returned for type PostMembersCreated
const PostMembersCreatedCode int = 201

/*
PostMembersCreated Created member.

swagger:response postMembersCreated
*/
type PostMembersCreated struct {

	/*
	  In: Body
	*/
	Payload *PostMembersCreatedBody `json:"body,omitempty"`
}

// NewPostMembersCreated creates PostMembersCreated with default headers values
func NewPostMembersCreated() *PostMembersCreated {

	return &PostMembersCreated{}
}

// WithPayload adds the payload to the post members created response
func (o *PostMembersCreated) WithPayload(payload *PostMembersCreatedBody) *PostMembersCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post members created response
func (o *PostMembersCreated) SetPayload(payload *PostMembersCreatedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMembersCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostMembersBadRequestCode is the HTTP code returned for type PostMembersBadRequest
const PostMembersBadRequestCode int = 400

/*
PostMembersBadRequest Bad request

swagger:response postMembersBadRequest
*/
type PostMembersBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMembersBadRequest creates PostMembersBadRequest with default headers values
func NewPostMembersBadRequest() *PostMembersBadRequest {

	return &PostMembersBadRequest{}
}

// WithPayload adds the payload to the post members bad request response
func (o *PostMembersBadRequest) WithPayload(payload *models.Error) *PostMembersBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post members bad request response
func (o *PostMembersBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMembersBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostMembersNotFoundCode is the HTTP code returned for type PostMembersNotFound
const PostMembersNotFoundCode int = 404

/*
PostMembersNotFound Not Found

swagger:response postMembersNotFound
*/
type PostMembersNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMembersNotFound creates PostMembersNotFound with default headers values
func NewPostMembersNotFound() *PostMembersNotFound {

	return &PostMembersNotFound{}
}

// WithPayload adds the payload to the post members not found response
func (o *PostMembersNotFound) WithPayload(payload *models.Error) *PostMembersNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post members not found response
func (o *PostMembersNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMembersNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PostMembersDefault Unexpected Error

swagger:response postMembersDefault
*/
type PostMembersDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMembersDefault creates PostMembersDefault with default headers values
func NewPostMembersDefault(code int) *PostMembersDefault {
	if code <= 0 {
		code = 500
	}

	return &PostMembersDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post members default response
func (o *PostMembersDefault) WithStatusCode(code int) *PostMembersDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post members default response
func (o *PostMembersDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post members default response
func (o *PostMembersDefault) WithPayload(payload *models.Error) *PostMembersDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post members default response
func (o *PostMembersDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMembersDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}