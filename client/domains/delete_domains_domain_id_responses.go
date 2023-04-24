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
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/models"
)

// DeleteDomainsDomainIDReader is a Reader for the DeleteDomainsDomainID structure.
type DeleteDomainsDomainIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteDomainsDomainIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteDomainsDomainIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeleteDomainsDomainIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteDomainsDomainIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteDomainsDomainIDNoContent creates a DeleteDomainsDomainIDNoContent with default headers values
func NewDeleteDomainsDomainIDNoContent() *DeleteDomainsDomainIDNoContent {
	return &DeleteDomainsDomainIDNoContent{}
}

/*
DeleteDomainsDomainIDNoContent describes a response with status code 204, with default header values.

Resource successfully deleted.
*/
type DeleteDomainsDomainIDNoContent struct {
}

// IsSuccess returns true when this delete domains domain Id no content response has a 2xx status code
func (o *DeleteDomainsDomainIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete domains domain Id no content response has a 3xx status code
func (o *DeleteDomainsDomainIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete domains domain Id no content response has a 4xx status code
func (o *DeleteDomainsDomainIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete domains domain Id no content response has a 5xx status code
func (o *DeleteDomainsDomainIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this delete domains domain Id no content response a status code equal to that given
func (o *DeleteDomainsDomainIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the delete domains domain Id no content response
func (o *DeleteDomainsDomainIDNoContent) Code() int {
	return 204
}

func (o *DeleteDomainsDomainIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /domains/{domain_id}][%d] deleteDomainsDomainIdNoContent ", 204)
}

func (o *DeleteDomainsDomainIDNoContent) String() string {
	return fmt.Sprintf("[DELETE /domains/{domain_id}][%d] deleteDomainsDomainIdNoContent ", 204)
}

func (o *DeleteDomainsDomainIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDomainsDomainIDNotFound creates a DeleteDomainsDomainIDNotFound with default headers values
func NewDeleteDomainsDomainIDNotFound() *DeleteDomainsDomainIDNotFound {
	return &DeleteDomainsDomainIDNotFound{}
}

/*
DeleteDomainsDomainIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type DeleteDomainsDomainIDNotFound struct {
	Payload *models.Error
}

// IsSuccess returns true when this delete domains domain Id not found response has a 2xx status code
func (o *DeleteDomainsDomainIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete domains domain Id not found response has a 3xx status code
func (o *DeleteDomainsDomainIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete domains domain Id not found response has a 4xx status code
func (o *DeleteDomainsDomainIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete domains domain Id not found response has a 5xx status code
func (o *DeleteDomainsDomainIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete domains domain Id not found response a status code equal to that given
func (o *DeleteDomainsDomainIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete domains domain Id not found response
func (o *DeleteDomainsDomainIDNotFound) Code() int {
	return 404
}

func (o *DeleteDomainsDomainIDNotFound) Error() string {
	return fmt.Sprintf("[DELETE /domains/{domain_id}][%d] deleteDomainsDomainIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteDomainsDomainIDNotFound) String() string {
	return fmt.Sprintf("[DELETE /domains/{domain_id}][%d] deleteDomainsDomainIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteDomainsDomainIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteDomainsDomainIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteDomainsDomainIDDefault creates a DeleteDomainsDomainIDDefault with default headers values
func NewDeleteDomainsDomainIDDefault(code int) *DeleteDomainsDomainIDDefault {
	return &DeleteDomainsDomainIDDefault{
		_statusCode: code,
	}
}

/*
DeleteDomainsDomainIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type DeleteDomainsDomainIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this delete domains domain ID default response has a 2xx status code
func (o *DeleteDomainsDomainIDDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete domains domain ID default response has a 3xx status code
func (o *DeleteDomainsDomainIDDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete domains domain ID default response has a 4xx status code
func (o *DeleteDomainsDomainIDDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete domains domain ID default response has a 5xx status code
func (o *DeleteDomainsDomainIDDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete domains domain ID default response a status code equal to that given
func (o *DeleteDomainsDomainIDDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete domains domain ID default response
func (o *DeleteDomainsDomainIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteDomainsDomainIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /domains/{domain_id}][%d] DeleteDomainsDomainID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteDomainsDomainIDDefault) String() string {
	return fmt.Sprintf("[DELETE /domains/{domain_id}][%d] DeleteDomainsDomainID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteDomainsDomainIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteDomainsDomainIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
