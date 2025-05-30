// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package monitors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/models"
)

// DeleteMonitorsMonitorIDReader is a Reader for the DeleteMonitorsMonitorID structure.
type DeleteMonitorsMonitorIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteMonitorsMonitorIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteMonitorsMonitorIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeleteMonitorsMonitorIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteMonitorsMonitorIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteMonitorsMonitorIDNoContent creates a DeleteMonitorsMonitorIDNoContent with default headers values
func NewDeleteMonitorsMonitorIDNoContent() *DeleteMonitorsMonitorIDNoContent {
	return &DeleteMonitorsMonitorIDNoContent{}
}

/*
DeleteMonitorsMonitorIDNoContent describes a response with status code 204, with default header values.

Resource successfully deleted.
*/
type DeleteMonitorsMonitorIDNoContent struct {
}

// IsSuccess returns true when this delete monitors monitor Id no content response has a 2xx status code
func (o *DeleteMonitorsMonitorIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete monitors monitor Id no content response has a 3xx status code
func (o *DeleteMonitorsMonitorIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete monitors monitor Id no content response has a 4xx status code
func (o *DeleteMonitorsMonitorIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete monitors monitor Id no content response has a 5xx status code
func (o *DeleteMonitorsMonitorIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this delete monitors monitor Id no content response a status code equal to that given
func (o *DeleteMonitorsMonitorIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the delete monitors monitor Id no content response
func (o *DeleteMonitorsMonitorIDNoContent) Code() int {
	return 204
}

func (o *DeleteMonitorsMonitorIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /monitors/{monitor_id}][%d] deleteMonitorsMonitorIdNoContent ", 204)
}

func (o *DeleteMonitorsMonitorIDNoContent) String() string {
	return fmt.Sprintf("[DELETE /monitors/{monitor_id}][%d] deleteMonitorsMonitorIdNoContent ", 204)
}

func (o *DeleteMonitorsMonitorIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteMonitorsMonitorIDNotFound creates a DeleteMonitorsMonitorIDNotFound with default headers values
func NewDeleteMonitorsMonitorIDNotFound() *DeleteMonitorsMonitorIDNotFound {
	return &DeleteMonitorsMonitorIDNotFound{}
}

/*
DeleteMonitorsMonitorIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type DeleteMonitorsMonitorIDNotFound struct {
	Payload *models.Error
}

// IsSuccess returns true when this delete monitors monitor Id not found response has a 2xx status code
func (o *DeleteMonitorsMonitorIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete monitors monitor Id not found response has a 3xx status code
func (o *DeleteMonitorsMonitorIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete monitors monitor Id not found response has a 4xx status code
func (o *DeleteMonitorsMonitorIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete monitors monitor Id not found response has a 5xx status code
func (o *DeleteMonitorsMonitorIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete monitors monitor Id not found response a status code equal to that given
func (o *DeleteMonitorsMonitorIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete monitors monitor Id not found response
func (o *DeleteMonitorsMonitorIDNotFound) Code() int {
	return 404
}

func (o *DeleteMonitorsMonitorIDNotFound) Error() string {
	return fmt.Sprintf("[DELETE /monitors/{monitor_id}][%d] deleteMonitorsMonitorIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteMonitorsMonitorIDNotFound) String() string {
	return fmt.Sprintf("[DELETE /monitors/{monitor_id}][%d] deleteMonitorsMonitorIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteMonitorsMonitorIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteMonitorsMonitorIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteMonitorsMonitorIDDefault creates a DeleteMonitorsMonitorIDDefault with default headers values
func NewDeleteMonitorsMonitorIDDefault(code int) *DeleteMonitorsMonitorIDDefault {
	return &DeleteMonitorsMonitorIDDefault{
		_statusCode: code,
	}
}

/*
DeleteMonitorsMonitorIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type DeleteMonitorsMonitorIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this delete monitors monitor ID default response has a 2xx status code
func (o *DeleteMonitorsMonitorIDDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete monitors monitor ID default response has a 3xx status code
func (o *DeleteMonitorsMonitorIDDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete monitors monitor ID default response has a 4xx status code
func (o *DeleteMonitorsMonitorIDDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete monitors monitor ID default response has a 5xx status code
func (o *DeleteMonitorsMonitorIDDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete monitors monitor ID default response a status code equal to that given
func (o *DeleteMonitorsMonitorIDDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete monitors monitor ID default response
func (o *DeleteMonitorsMonitorIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteMonitorsMonitorIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /monitors/{monitor_id}][%d] DeleteMonitorsMonitorID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteMonitorsMonitorIDDefault) String() string {
	return fmt.Sprintf("[DELETE /monitors/{monitor_id}][%d] DeleteMonitorsMonitorID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteMonitorsMonitorIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteMonitorsMonitorIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
