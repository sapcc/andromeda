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
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/models"
)

// DeletePoolsPoolIDReader is a Reader for the DeletePoolsPoolID structure.
type DeletePoolsPoolIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeletePoolsPoolIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeletePoolsPoolIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeletePoolsPoolIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeletePoolsPoolIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeletePoolsPoolIDNoContent creates a DeletePoolsPoolIDNoContent with default headers values
func NewDeletePoolsPoolIDNoContent() *DeletePoolsPoolIDNoContent {
	return &DeletePoolsPoolIDNoContent{}
}

/*
DeletePoolsPoolIDNoContent describes a response with status code 204, with default header values.

Resource successfully deleted.
*/
type DeletePoolsPoolIDNoContent struct {
}

// IsSuccess returns true when this delete pools pool Id no content response has a 2xx status code
func (o *DeletePoolsPoolIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete pools pool Id no content response has a 3xx status code
func (o *DeletePoolsPoolIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete pools pool Id no content response has a 4xx status code
func (o *DeletePoolsPoolIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete pools pool Id no content response has a 5xx status code
func (o *DeletePoolsPoolIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this delete pools pool Id no content response a status code equal to that given
func (o *DeletePoolsPoolIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the delete pools pool Id no content response
func (o *DeletePoolsPoolIDNoContent) Code() int {
	return 204
}

func (o *DeletePoolsPoolIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /pools/{pool_id}][%d] deletePoolsPoolIdNoContent ", 204)
}

func (o *DeletePoolsPoolIDNoContent) String() string {
	return fmt.Sprintf("[DELETE /pools/{pool_id}][%d] deletePoolsPoolIdNoContent ", 204)
}

func (o *DeletePoolsPoolIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeletePoolsPoolIDNotFound creates a DeletePoolsPoolIDNotFound with default headers values
func NewDeletePoolsPoolIDNotFound() *DeletePoolsPoolIDNotFound {
	return &DeletePoolsPoolIDNotFound{}
}

/*
DeletePoolsPoolIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type DeletePoolsPoolIDNotFound struct {
	Payload *models.Error
}

// IsSuccess returns true when this delete pools pool Id not found response has a 2xx status code
func (o *DeletePoolsPoolIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete pools pool Id not found response has a 3xx status code
func (o *DeletePoolsPoolIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete pools pool Id not found response has a 4xx status code
func (o *DeletePoolsPoolIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete pools pool Id not found response has a 5xx status code
func (o *DeletePoolsPoolIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete pools pool Id not found response a status code equal to that given
func (o *DeletePoolsPoolIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete pools pool Id not found response
func (o *DeletePoolsPoolIDNotFound) Code() int {
	return 404
}

func (o *DeletePoolsPoolIDNotFound) Error() string {
	return fmt.Sprintf("[DELETE /pools/{pool_id}][%d] deletePoolsPoolIdNotFound  %+v", 404, o.Payload)
}

func (o *DeletePoolsPoolIDNotFound) String() string {
	return fmt.Sprintf("[DELETE /pools/{pool_id}][%d] deletePoolsPoolIdNotFound  %+v", 404, o.Payload)
}

func (o *DeletePoolsPoolIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeletePoolsPoolIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeletePoolsPoolIDDefault creates a DeletePoolsPoolIDDefault with default headers values
func NewDeletePoolsPoolIDDefault(code int) *DeletePoolsPoolIDDefault {
	return &DeletePoolsPoolIDDefault{
		_statusCode: code,
	}
}

/*
DeletePoolsPoolIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type DeletePoolsPoolIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this delete pools pool ID default response has a 2xx status code
func (o *DeletePoolsPoolIDDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete pools pool ID default response has a 3xx status code
func (o *DeletePoolsPoolIDDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete pools pool ID default response has a 4xx status code
func (o *DeletePoolsPoolIDDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete pools pool ID default response has a 5xx status code
func (o *DeletePoolsPoolIDDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete pools pool ID default response a status code equal to that given
func (o *DeletePoolsPoolIDDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete pools pool ID default response
func (o *DeletePoolsPoolIDDefault) Code() int {
	return o._statusCode
}

func (o *DeletePoolsPoolIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /pools/{pool_id}][%d] DeletePoolsPoolID default  %+v", o._statusCode, o.Payload)
}

func (o *DeletePoolsPoolIDDefault) String() string {
	return fmt.Sprintf("[DELETE /pools/{pool_id}][%d] DeletePoolsPoolID default  %+v", o._statusCode, o.Payload)
}

func (o *DeletePoolsPoolIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeletePoolsPoolIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}