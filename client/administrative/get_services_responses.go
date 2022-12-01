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

// GetServicesReader is a Reader for the GetServices structure.
type GetServicesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetServicesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetServicesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetServicesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetServicesOK creates a GetServicesOK with default headers values
func NewGetServicesOK() *GetServicesOK {
	return &GetServicesOK{}
}

/*
GetServicesOK describes a response with status code 200, with default header values.

A JSON array of services
*/
type GetServicesOK struct {
	Payload *GetServicesOKBody
}

// IsSuccess returns true when this get services o k response has a 2xx status code
func (o *GetServicesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get services o k response has a 3xx status code
func (o *GetServicesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get services o k response has a 4xx status code
func (o *GetServicesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get services o k response has a 5xx status code
func (o *GetServicesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get services o k response a status code equal to that given
func (o *GetServicesOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetServicesOK) Error() string {
	return fmt.Sprintf("[GET /services][%d] getServicesOK  %+v", 200, o.Payload)
}

func (o *GetServicesOK) String() string {
	return fmt.Sprintf("[GET /services][%d] getServicesOK  %+v", 200, o.Payload)
}

func (o *GetServicesOK) GetPayload() *GetServicesOKBody {
	return o.Payload
}

func (o *GetServicesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetServicesOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetServicesDefault creates a GetServicesDefault with default headers values
func NewGetServicesDefault(code int) *GetServicesDefault {
	return &GetServicesDefault{
		_statusCode: code,
	}
}

/*
GetServicesDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type GetServicesDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get services default response
func (o *GetServicesDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this get services default response has a 2xx status code
func (o *GetServicesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get services default response has a 3xx status code
func (o *GetServicesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get services default response has a 4xx status code
func (o *GetServicesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get services default response has a 5xx status code
func (o *GetServicesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get services default response a status code equal to that given
func (o *GetServicesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *GetServicesDefault) Error() string {
	return fmt.Sprintf("[GET /services][%d] GetServices default  %+v", o._statusCode, o.Payload)
}

func (o *GetServicesDefault) String() string {
	return fmt.Sprintf("[GET /services][%d] GetServices default  %+v", o._statusCode, o.Payload)
}

func (o *GetServicesDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetServicesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetServicesOKBody get services o k body
swagger:model GetServicesOKBody
*/
type GetServicesOKBody struct {

	// services
	Services []*models.Service `json:"services"`
}

// Validate validates this get services o k body
func (o *GetServicesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateServices(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetServicesOKBody) validateServices(formats strfmt.Registry) error {
	if swag.IsZero(o.Services) { // not required
		return nil
	}

	for i := 0; i < len(o.Services); i++ {
		if swag.IsZero(o.Services[i]) { // not required
			continue
		}

		if o.Services[i] != nil {
			if err := o.Services[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get services o k body based on the context it is used
func (o *GetServicesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateServices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetServicesOKBody) contextValidateServices(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Services); i++ {

		if o.Services[i] != nil {
			if err := o.Services[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getServicesOK" + "." + "services" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetServicesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetServicesOKBody) UnmarshalBinary(b []byte) error {
	var res GetServicesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
