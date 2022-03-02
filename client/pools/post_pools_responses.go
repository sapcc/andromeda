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

// PostPoolsReader is a Reader for the PostPools structure.
type PostPoolsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostPoolsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostPoolsCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostPoolsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPostPoolsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostPoolsCreated creates a PostPoolsCreated with default headers values
func NewPostPoolsCreated() *PostPoolsCreated {
	return &PostPoolsCreated{}
}

/* PostPoolsCreated describes a response with status code 201, with default header values.

Created pool.
*/
type PostPoolsCreated struct {
	Payload *PostPoolsCreatedBody
}

func (o *PostPoolsCreated) Error() string {
	return fmt.Sprintf("[POST /pools][%d] postPoolsCreated  %+v", 201, o.Payload)
}
func (o *PostPoolsCreated) GetPayload() *PostPoolsCreatedBody {
	return o.Payload
}

func (o *PostPoolsCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PostPoolsCreatedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostPoolsBadRequest creates a PostPoolsBadRequest with default headers values
func NewPostPoolsBadRequest() *PostPoolsBadRequest {
	return &PostPoolsBadRequest{}
}

/* PostPoolsBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type PostPoolsBadRequest struct {
	Payload *models.Error
}

func (o *PostPoolsBadRequest) Error() string {
	return fmt.Sprintf("[POST /pools][%d] postPoolsBadRequest  %+v", 400, o.Payload)
}
func (o *PostPoolsBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostPoolsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostPoolsDefault creates a PostPoolsDefault with default headers values
func NewPostPoolsDefault(code int) *PostPoolsDefault {
	return &PostPoolsDefault{
		_statusCode: code,
	}
}

/* PostPoolsDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type PostPoolsDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the post pools default response
func (o *PostPoolsDefault) Code() int {
	return o._statusCode
}

func (o *PostPoolsDefault) Error() string {
	return fmt.Sprintf("[POST /pools][%d] PostPools default  %+v", o._statusCode, o.Payload)
}
func (o *PostPoolsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostPoolsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*PostPoolsBody post pools body
swagger:model PostPoolsBody
*/
type PostPoolsBody struct {

	// pool
	// Required: true
	Pool *models.Pool `json:"pool"`
}

// Validate validates this post pools body
func (o *PostPoolsBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePool(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostPoolsBody) validatePool(formats strfmt.Registry) error {

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

// ContextValidate validate this post pools body based on the context it is used
func (o *PostPoolsBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidatePool(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostPoolsBody) contextValidatePool(ctx context.Context, formats strfmt.Registry) error {

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
func (o *PostPoolsBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostPoolsBody) UnmarshalBinary(b []byte) error {
	var res PostPoolsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*PostPoolsCreatedBody post pools created body
swagger:model PostPoolsCreatedBody
*/
type PostPoolsCreatedBody struct {

	// pool
	Pool *models.Pool `json:"pool,omitempty"`
}

// Validate validates this post pools created body
func (o *PostPoolsCreatedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePool(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostPoolsCreatedBody) validatePool(formats strfmt.Registry) error {
	if swag.IsZero(o.Pool) { // not required
		return nil
	}

	if o.Pool != nil {
		if err := o.Pool.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postPoolsCreated" + "." + "pool")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postPoolsCreated" + "." + "pool")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post pools created body based on the context it is used
func (o *PostPoolsCreatedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidatePool(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostPoolsCreatedBody) contextValidatePool(ctx context.Context, formats strfmt.Registry) error {

	if o.Pool != nil {
		if err := o.Pool.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postPoolsCreated" + "." + "pool")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postPoolsCreated" + "." + "pool")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostPoolsCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostPoolsCreatedBody) UnmarshalBinary(b []byte) error {
	var res PostPoolsCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
