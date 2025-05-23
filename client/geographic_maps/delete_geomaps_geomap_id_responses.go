// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package geographic_maps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/models"
)

// DeleteGeomapsGeomapIDReader is a Reader for the DeleteGeomapsGeomapID structure.
type DeleteGeomapsGeomapIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteGeomapsGeomapIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteGeomapsGeomapIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeleteGeomapsGeomapIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteGeomapsGeomapIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteGeomapsGeomapIDNoContent creates a DeleteGeomapsGeomapIDNoContent with default headers values
func NewDeleteGeomapsGeomapIDNoContent() *DeleteGeomapsGeomapIDNoContent {
	return &DeleteGeomapsGeomapIDNoContent{}
}

/*
DeleteGeomapsGeomapIDNoContent describes a response with status code 204, with default header values.

Resource successfully deleted.
*/
type DeleteGeomapsGeomapIDNoContent struct {
}

// IsSuccess returns true when this delete geomaps geomap Id no content response has a 2xx status code
func (o *DeleteGeomapsGeomapIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete geomaps geomap Id no content response has a 3xx status code
func (o *DeleteGeomapsGeomapIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete geomaps geomap Id no content response has a 4xx status code
func (o *DeleteGeomapsGeomapIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete geomaps geomap Id no content response has a 5xx status code
func (o *DeleteGeomapsGeomapIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this delete geomaps geomap Id no content response a status code equal to that given
func (o *DeleteGeomapsGeomapIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the delete geomaps geomap Id no content response
func (o *DeleteGeomapsGeomapIDNoContent) Code() int {
	return 204
}

func (o *DeleteGeomapsGeomapIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /geomaps/{geomap_id}][%d] deleteGeomapsGeomapIdNoContent ", 204)
}

func (o *DeleteGeomapsGeomapIDNoContent) String() string {
	return fmt.Sprintf("[DELETE /geomaps/{geomap_id}][%d] deleteGeomapsGeomapIdNoContent ", 204)
}

func (o *DeleteGeomapsGeomapIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteGeomapsGeomapIDNotFound creates a DeleteGeomapsGeomapIDNotFound with default headers values
func NewDeleteGeomapsGeomapIDNotFound() *DeleteGeomapsGeomapIDNotFound {
	return &DeleteGeomapsGeomapIDNotFound{}
}

/*
DeleteGeomapsGeomapIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type DeleteGeomapsGeomapIDNotFound struct {
	Payload *models.Error
}

// IsSuccess returns true when this delete geomaps geomap Id not found response has a 2xx status code
func (o *DeleteGeomapsGeomapIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete geomaps geomap Id not found response has a 3xx status code
func (o *DeleteGeomapsGeomapIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete geomaps geomap Id not found response has a 4xx status code
func (o *DeleteGeomapsGeomapIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete geomaps geomap Id not found response has a 5xx status code
func (o *DeleteGeomapsGeomapIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete geomaps geomap Id not found response a status code equal to that given
func (o *DeleteGeomapsGeomapIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete geomaps geomap Id not found response
func (o *DeleteGeomapsGeomapIDNotFound) Code() int {
	return 404
}

func (o *DeleteGeomapsGeomapIDNotFound) Error() string {
	return fmt.Sprintf("[DELETE /geomaps/{geomap_id}][%d] deleteGeomapsGeomapIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteGeomapsGeomapIDNotFound) String() string {
	return fmt.Sprintf("[DELETE /geomaps/{geomap_id}][%d] deleteGeomapsGeomapIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteGeomapsGeomapIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteGeomapsGeomapIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteGeomapsGeomapIDDefault creates a DeleteGeomapsGeomapIDDefault with default headers values
func NewDeleteGeomapsGeomapIDDefault(code int) *DeleteGeomapsGeomapIDDefault {
	return &DeleteGeomapsGeomapIDDefault{
		_statusCode: code,
	}
}

/*
DeleteGeomapsGeomapIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type DeleteGeomapsGeomapIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this delete geomaps geomap ID default response has a 2xx status code
func (o *DeleteGeomapsGeomapIDDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete geomaps geomap ID default response has a 3xx status code
func (o *DeleteGeomapsGeomapIDDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete geomaps geomap ID default response has a 4xx status code
func (o *DeleteGeomapsGeomapIDDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete geomaps geomap ID default response has a 5xx status code
func (o *DeleteGeomapsGeomapIDDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete geomaps geomap ID default response a status code equal to that given
func (o *DeleteGeomapsGeomapIDDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete geomaps geomap ID default response
func (o *DeleteGeomapsGeomapIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteGeomapsGeomapIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /geomaps/{geomap_id}][%d] DeleteGeomapsGeomapID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteGeomapsGeomapIDDefault) String() string {
	return fmt.Sprintf("[DELETE /geomaps/{geomap_id}][%d] DeleteGeomapsGeomapID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteGeomapsGeomapIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteGeomapsGeomapIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
