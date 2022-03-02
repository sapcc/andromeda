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

package datacenters

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

// PostDatacentersReader is a Reader for the PostDatacenters structure.
type PostDatacentersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostDatacentersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostDatacentersCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewPostDatacentersNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPostDatacentersDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostDatacentersCreated creates a PostDatacentersCreated with default headers values
func NewPostDatacentersCreated() *PostDatacentersCreated {
	return &PostDatacentersCreated{}
}

/* PostDatacentersCreated describes a response with status code 201, with default header values.

Created datacenter.
*/
type PostDatacentersCreated struct {
	Payload *PostDatacentersCreatedBody
}

func (o *PostDatacentersCreated) Error() string {
	return fmt.Sprintf("[POST /datacenters][%d] postDatacentersCreated  %+v", 201, o.Payload)
}
func (o *PostDatacentersCreated) GetPayload() *PostDatacentersCreatedBody {
	return o.Payload
}

func (o *PostDatacentersCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PostDatacentersCreatedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostDatacentersNotFound creates a PostDatacentersNotFound with default headers values
func NewPostDatacentersNotFound() *PostDatacentersNotFound {
	return &PostDatacentersNotFound{}
}

/* PostDatacentersNotFound describes a response with status code 404, with default header values.

Not Found
*/
type PostDatacentersNotFound struct {
	Payload *models.Error
}

func (o *PostDatacentersNotFound) Error() string {
	return fmt.Sprintf("[POST /datacenters][%d] postDatacentersNotFound  %+v", 404, o.Payload)
}
func (o *PostDatacentersNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostDatacentersNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostDatacentersDefault creates a PostDatacentersDefault with default headers values
func NewPostDatacentersDefault(code int) *PostDatacentersDefault {
	return &PostDatacentersDefault{
		_statusCode: code,
	}
}

/* PostDatacentersDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type PostDatacentersDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the post datacenters default response
func (o *PostDatacentersDefault) Code() int {
	return o._statusCode
}

func (o *PostDatacentersDefault) Error() string {
	return fmt.Sprintf("[POST /datacenters][%d] PostDatacenters default  %+v", o._statusCode, o.Payload)
}
func (o *PostDatacentersDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostDatacentersDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*PostDatacentersBody post datacenters body
swagger:model PostDatacentersBody
*/
type PostDatacentersBody struct {

	// datacenter
	// Required: true
	Datacenter *models.Datacenter `json:"datacenter"`
}

// Validate validates this post datacenters body
func (o *PostDatacentersBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDatacenter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersBody) validateDatacenter(formats strfmt.Registry) error {

	if err := validate.Required("datacenter"+"."+"datacenter", "body", o.Datacenter); err != nil {
		return err
	}

	if o.Datacenter != nil {
		if err := o.Datacenter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("datacenter" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("datacenter" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post datacenters body based on the context it is used
func (o *PostDatacentersBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDatacenter(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersBody) contextValidateDatacenter(ctx context.Context, formats strfmt.Registry) error {

	if o.Datacenter != nil {
		if err := o.Datacenter.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("datacenter" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("datacenter" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostDatacentersBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDatacentersBody) UnmarshalBinary(b []byte) error {
	var res PostDatacentersBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*PostDatacentersCreatedBody post datacenters created body
swagger:model PostDatacentersCreatedBody
*/
type PostDatacentersCreatedBody struct {

	// datacenter
	Datacenter *models.Datacenter `json:"datacenter,omitempty"`
}

// Validate validates this post datacenters created body
func (o *PostDatacentersCreatedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDatacenter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersCreatedBody) validateDatacenter(formats strfmt.Registry) error {
	if swag.IsZero(o.Datacenter) { // not required
		return nil
	}

	if o.Datacenter != nil {
		if err := o.Datacenter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postDatacentersCreated" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postDatacentersCreated" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post datacenters created body based on the context it is used
func (o *PostDatacentersCreatedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDatacenter(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDatacentersCreatedBody) contextValidateDatacenter(ctx context.Context, formats strfmt.Registry) error {

	if o.Datacenter != nil {
		if err := o.Datacenter.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postDatacentersCreated" + "." + "datacenter")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postDatacentersCreated" + "." + "datacenter")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostDatacentersCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDatacentersCreatedBody) UnmarshalBinary(b []byte) error {
	var res PostDatacentersCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
