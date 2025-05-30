// Code generated by go-swagger; DO NOT EDIT.

// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package geographic_maps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/models"
)

// GetGeomapsReader is a Reader for the GetGeomaps structure.
type GetGeomapsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetGeomapsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetGeomapsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetGeomapsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetGeomapsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetGeomapsOK creates a GetGeomapsOK with default headers values
func NewGetGeomapsOK() *GetGeomapsOK {
	return &GetGeomapsOK{}
}

/*
GetGeomapsOK describes a response with status code 200, with default header values.

A JSON array of geographic maps
*/
type GetGeomapsOK struct {
	Payload *GetGeomapsOKBody
}

// IsSuccess returns true when this get geomaps o k response has a 2xx status code
func (o *GetGeomapsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get geomaps o k response has a 3xx status code
func (o *GetGeomapsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get geomaps o k response has a 4xx status code
func (o *GetGeomapsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get geomaps o k response has a 5xx status code
func (o *GetGeomapsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get geomaps o k response a status code equal to that given
func (o *GetGeomapsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get geomaps o k response
func (o *GetGeomapsOK) Code() int {
	return 200
}

func (o *GetGeomapsOK) Error() string {
	return fmt.Sprintf("[GET /geomaps][%d] getGeomapsOK  %+v", 200, o.Payload)
}

func (o *GetGeomapsOK) String() string {
	return fmt.Sprintf("[GET /geomaps][%d] getGeomapsOK  %+v", 200, o.Payload)
}

func (o *GetGeomapsOK) GetPayload() *GetGeomapsOKBody {
	return o.Payload
}

func (o *GetGeomapsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetGeomapsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetGeomapsBadRequest creates a GetGeomapsBadRequest with default headers values
func NewGetGeomapsBadRequest() *GetGeomapsBadRequest {
	return &GetGeomapsBadRequest{}
}

/*
GetGeomapsBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type GetGeomapsBadRequest struct {
	Payload *models.Error
}

// IsSuccess returns true when this get geomaps bad request response has a 2xx status code
func (o *GetGeomapsBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get geomaps bad request response has a 3xx status code
func (o *GetGeomapsBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get geomaps bad request response has a 4xx status code
func (o *GetGeomapsBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get geomaps bad request response has a 5xx status code
func (o *GetGeomapsBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get geomaps bad request response a status code equal to that given
func (o *GetGeomapsBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get geomaps bad request response
func (o *GetGeomapsBadRequest) Code() int {
	return 400
}

func (o *GetGeomapsBadRequest) Error() string {
	return fmt.Sprintf("[GET /geomaps][%d] getGeomapsBadRequest  %+v", 400, o.Payload)
}

func (o *GetGeomapsBadRequest) String() string {
	return fmt.Sprintf("[GET /geomaps][%d] getGeomapsBadRequest  %+v", 400, o.Payload)
}

func (o *GetGeomapsBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetGeomapsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetGeomapsDefault creates a GetGeomapsDefault with default headers values
func NewGetGeomapsDefault(code int) *GetGeomapsDefault {
	return &GetGeomapsDefault{
		_statusCode: code,
	}
}

/*
GetGeomapsDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type GetGeomapsDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get geomaps default response has a 2xx status code
func (o *GetGeomapsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get geomaps default response has a 3xx status code
func (o *GetGeomapsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get geomaps default response has a 4xx status code
func (o *GetGeomapsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get geomaps default response has a 5xx status code
func (o *GetGeomapsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get geomaps default response a status code equal to that given
func (o *GetGeomapsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get geomaps default response
func (o *GetGeomapsDefault) Code() int {
	return o._statusCode
}

func (o *GetGeomapsDefault) Error() string {
	return fmt.Sprintf("[GET /geomaps][%d] GetGeomaps default  %+v", o._statusCode, o.Payload)
}

func (o *GetGeomapsDefault) String() string {
	return fmt.Sprintf("[GET /geomaps][%d] GetGeomaps default  %+v", o._statusCode, o.Payload)
}

func (o *GetGeomapsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetGeomapsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetGeomapsOKBody get geomaps o k body
swagger:model GetGeomapsOKBody
*/
type GetGeomapsOKBody struct {

	// geomaps
	Geomaps []*models.Geomap `json:"geomaps"`

	// links
	Links []*models.Link `json:"links,omitempty"`
}

// Validate validates this get geomaps o k body
func (o *GetGeomapsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGeomaps(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetGeomapsOKBody) validateGeomaps(formats strfmt.Registry) error {
	if swag.IsZero(o.Geomaps) { // not required
		return nil
	}

	for i := 0; i < len(o.Geomaps); i++ {
		if swag.IsZero(o.Geomaps[i]) { // not required
			continue
		}

		if o.Geomaps[i] != nil {
			if err := o.Geomaps[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetGeomapsOKBody) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(o.Links) { // not required
		return nil
	}

	for i := 0; i < len(o.Links); i++ {
		if swag.IsZero(o.Links[i]) { // not required
			continue
		}

		if o.Links[i] != nil {
			if err := o.Links[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get geomaps o k body based on the context it is used
func (o *GetGeomapsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateGeomaps(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetGeomapsOKBody) contextValidateGeomaps(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Geomaps); i++ {

		if o.Geomaps[i] != nil {
			if err := o.Geomaps[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "geomaps" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetGeomapsOKBody) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Links); i++ {

		if o.Links[i] != nil {
			if err := o.Links[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGeomapsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetGeomapsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetGeomapsOKBody) UnmarshalBinary(b []byte) error {
	var res GetGeomapsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
