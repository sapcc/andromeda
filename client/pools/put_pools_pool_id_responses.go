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
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/sapcc/andromeda/models"
)

// PutPoolsPoolIDReader is a Reader for the PutPoolsPoolID structure.
type PutPoolsPoolIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutPoolsPoolIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 202:
		result := NewPutPoolsPoolIDAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewPutPoolsPoolIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPutPoolsPoolIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPutPoolsPoolIDAccepted creates a PutPoolsPoolIDAccepted with default headers values
func NewPutPoolsPoolIDAccepted() *PutPoolsPoolIDAccepted {
	return &PutPoolsPoolIDAccepted{}
}

/*
PutPoolsPoolIDAccepted describes a response with status code 202, with default header values.

Updated pool.
*/
type PutPoolsPoolIDAccepted struct {
	Payload *PutPoolsPoolIDAcceptedBody
}

// IsSuccess returns true when this put pools pool Id accepted response has a 2xx status code
func (o *PutPoolsPoolIDAccepted) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this put pools pool Id accepted response has a 3xx status code
func (o *PutPoolsPoolIDAccepted) IsRedirect() bool {
	return false
}

// IsClientError returns true when this put pools pool Id accepted response has a 4xx status code
func (o *PutPoolsPoolIDAccepted) IsClientError() bool {
	return false
}

// IsServerError returns true when this put pools pool Id accepted response has a 5xx status code
func (o *PutPoolsPoolIDAccepted) IsServerError() bool {
	return false
}

// IsCode returns true when this put pools pool Id accepted response a status code equal to that given
func (o *PutPoolsPoolIDAccepted) IsCode(code int) bool {
	return code == 202
}

func (o *PutPoolsPoolIDAccepted) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}][%d] putPoolsPoolIdAccepted  %+v", 202, o.Payload)
}

func (o *PutPoolsPoolIDAccepted) String() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}][%d] putPoolsPoolIdAccepted  %+v", 202, o.Payload)
}

func (o *PutPoolsPoolIDAccepted) GetPayload() *PutPoolsPoolIDAcceptedBody {
	return o.Payload
}

func (o *PutPoolsPoolIDAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PutPoolsPoolIDAcceptedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutPoolsPoolIDNotFound creates a PutPoolsPoolIDNotFound with default headers values
func NewPutPoolsPoolIDNotFound() *PutPoolsPoolIDNotFound {
	return &PutPoolsPoolIDNotFound{}
}

/*
PutPoolsPoolIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type PutPoolsPoolIDNotFound struct {
	Payload *models.Error
}

// IsSuccess returns true when this put pools pool Id not found response has a 2xx status code
func (o *PutPoolsPoolIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this put pools pool Id not found response has a 3xx status code
func (o *PutPoolsPoolIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this put pools pool Id not found response has a 4xx status code
func (o *PutPoolsPoolIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this put pools pool Id not found response has a 5xx status code
func (o *PutPoolsPoolIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this put pools pool Id not found response a status code equal to that given
func (o *PutPoolsPoolIDNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *PutPoolsPoolIDNotFound) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}][%d] putPoolsPoolIdNotFound  %+v", 404, o.Payload)
}

func (o *PutPoolsPoolIDNotFound) String() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}][%d] putPoolsPoolIdNotFound  %+v", 404, o.Payload)
}

func (o *PutPoolsPoolIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *PutPoolsPoolIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutPoolsPoolIDDefault creates a PutPoolsPoolIDDefault with default headers values
func NewPutPoolsPoolIDDefault(code int) *PutPoolsPoolIDDefault {
	return &PutPoolsPoolIDDefault{
		_statusCode: code,
	}
}

/*
PutPoolsPoolIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type PutPoolsPoolIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the put pools pool ID default response
func (o *PutPoolsPoolIDDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this put pools pool ID default response has a 2xx status code
func (o *PutPoolsPoolIDDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this put pools pool ID default response has a 3xx status code
func (o *PutPoolsPoolIDDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this put pools pool ID default response has a 4xx status code
func (o *PutPoolsPoolIDDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this put pools pool ID default response has a 5xx status code
func (o *PutPoolsPoolIDDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this put pools pool ID default response a status code equal to that given
func (o *PutPoolsPoolIDDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *PutPoolsPoolIDDefault) Error() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}][%d] PutPoolsPoolID default  %+v", o._statusCode, o.Payload)
}

func (o *PutPoolsPoolIDDefault) String() string {
	return fmt.Sprintf("[PUT /pools/{pool_id}][%d] PutPoolsPoolID default  %+v", o._statusCode, o.Payload)
}

func (o *PutPoolsPoolIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PutPoolsPoolIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
PutPoolsPoolIDAcceptedBody put pools pool ID accepted body
swagger:model PutPoolsPoolIDAcceptedBody
*/
type PutPoolsPoolIDAcceptedBody struct {

	// pool
	Pool *models.Pool `json:"pool,omitempty"`
}

// Validate validates this put pools pool ID accepted body
func (o *PutPoolsPoolIDAcceptedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePool(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutPoolsPoolIDAcceptedBody) validatePool(formats strfmt.Registry) error {
	if swag.IsZero(o.Pool) { // not required
		return nil
	}

	if o.Pool != nil {
		if err := o.Pool.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putPoolsPoolIdAccepted" + "." + "pool")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putPoolsPoolIdAccepted" + "." + "pool")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put pools pool ID accepted body based on the context it is used
func (o *PutPoolsPoolIDAcceptedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidatePool(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutPoolsPoolIDAcceptedBody) contextValidatePool(ctx context.Context, formats strfmt.Registry) error {

	if o.Pool != nil {
		if err := o.Pool.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putPoolsPoolIdAccepted" + "." + "pool")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putPoolsPoolIdAccepted" + "." + "pool")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutPoolsPoolIDAcceptedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutPoolsPoolIDAcceptedBody) UnmarshalBinary(b []byte) error {
	var res PutPoolsPoolIDAcceptedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
PutPoolsPoolIDBody put pools pool ID body
swagger:model PutPoolsPoolIDBody
*/
type PutPoolsPoolIDBody struct {

	// pool
	// Required: true
	Pool *models.Pool `json:"pool"`
}

// Validate validates this put pools pool ID body
func (o *PutPoolsPoolIDBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePool(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutPoolsPoolIDBody) validatePool(formats strfmt.Registry) error {

	if err := validate.Required("pool"+"."+"pool", "body", o.Pool); err != nil {
		return err
	}

	if o.Pool != nil {
		if err := o.Pool.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pool" + "." + "pool")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("pool" + "." + "pool")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put pools pool ID body based on the context it is used
func (o *PutPoolsPoolIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidatePool(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutPoolsPoolIDBody) contextValidatePool(ctx context.Context, formats strfmt.Registry) error {

	if o.Pool != nil {
		if err := o.Pool.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pool" + "." + "pool")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("pool" + "." + "pool")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutPoolsPoolIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutPoolsPoolIDBody) UnmarshalBinary(b []byte) error {
	var res PutPoolsPoolIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
