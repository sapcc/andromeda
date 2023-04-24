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

// PostDomainsReader is a Reader for the PostDomains structure.
type PostDomainsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostDomainsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostDomainsCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostDomainsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPostDomainsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostDomainsCreated creates a PostDomainsCreated with default headers values
func NewPostDomainsCreated() *PostDomainsCreated {
	return &PostDomainsCreated{}
}

/*
PostDomainsCreated describes a response with status code 201, with default header values.

Created domain.
*/
type PostDomainsCreated struct {
	Payload *PostDomainsCreatedBody
}

// IsSuccess returns true when this post domains created response has a 2xx status code
func (o *PostDomainsCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this post domains created response has a 3xx status code
func (o *PostDomainsCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post domains created response has a 4xx status code
func (o *PostDomainsCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this post domains created response has a 5xx status code
func (o *PostDomainsCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this post domains created response a status code equal to that given
func (o *PostDomainsCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the post domains created response
func (o *PostDomainsCreated) Code() int {
	return 201
}

func (o *PostDomainsCreated) Error() string {
	return fmt.Sprintf("[POST /domains][%d] postDomainsCreated  %+v", 201, o.Payload)
}

func (o *PostDomainsCreated) String() string {
	return fmt.Sprintf("[POST /domains][%d] postDomainsCreated  %+v", 201, o.Payload)
}

func (o *PostDomainsCreated) GetPayload() *PostDomainsCreatedBody {
	return o.Payload
}

func (o *PostDomainsCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PostDomainsCreatedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostDomainsBadRequest creates a PostDomainsBadRequest with default headers values
func NewPostDomainsBadRequest() *PostDomainsBadRequest {
	return &PostDomainsBadRequest{}
}

/*
PostDomainsBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type PostDomainsBadRequest struct {
	Payload *models.Error
}

// IsSuccess returns true when this post domains bad request response has a 2xx status code
func (o *PostDomainsBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post domains bad request response has a 3xx status code
func (o *PostDomainsBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post domains bad request response has a 4xx status code
func (o *PostDomainsBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this post domains bad request response has a 5xx status code
func (o *PostDomainsBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this post domains bad request response a status code equal to that given
func (o *PostDomainsBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the post domains bad request response
func (o *PostDomainsBadRequest) Code() int {
	return 400
}

func (o *PostDomainsBadRequest) Error() string {
	return fmt.Sprintf("[POST /domains][%d] postDomainsBadRequest  %+v", 400, o.Payload)
}

func (o *PostDomainsBadRequest) String() string {
	return fmt.Sprintf("[POST /domains][%d] postDomainsBadRequest  %+v", 400, o.Payload)
}

func (o *PostDomainsBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostDomainsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostDomainsDefault creates a PostDomainsDefault with default headers values
func NewPostDomainsDefault(code int) *PostDomainsDefault {
	return &PostDomainsDefault{
		_statusCode: code,
	}
}

/*
PostDomainsDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type PostDomainsDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this post domains default response has a 2xx status code
func (o *PostDomainsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this post domains default response has a 3xx status code
func (o *PostDomainsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this post domains default response has a 4xx status code
func (o *PostDomainsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this post domains default response has a 5xx status code
func (o *PostDomainsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this post domains default response a status code equal to that given
func (o *PostDomainsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the post domains default response
func (o *PostDomainsDefault) Code() int {
	return o._statusCode
}

func (o *PostDomainsDefault) Error() string {
	return fmt.Sprintf("[POST /domains][%d] PostDomains default  %+v", o._statusCode, o.Payload)
}

func (o *PostDomainsDefault) String() string {
	return fmt.Sprintf("[POST /domains][%d] PostDomains default  %+v", o._statusCode, o.Payload)
}

func (o *PostDomainsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostDomainsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
PostDomainsBody post domains body
swagger:model PostDomainsBody
*/
type PostDomainsBody struct {

	// domain
	// Required: true
	Domain *models.Domain `json:"domain"`
}

// Validate validates this post domains body
func (o *PostDomainsBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomain(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDomainsBody) validateDomain(formats strfmt.Registry) error {

	if err := validate.Required("domain"+"."+"domain", "body", o.Domain); err != nil {
		return err
	}

	if o.Domain != nil {
		if err := o.Domain.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("domain" + "." + "domain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("domain" + "." + "domain")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post domains body based on the context it is used
func (o *PostDomainsBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDomain(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDomainsBody) contextValidateDomain(ctx context.Context, formats strfmt.Registry) error {

	if o.Domain != nil {
		if err := o.Domain.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("domain" + "." + "domain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("domain" + "." + "domain")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostDomainsBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDomainsBody) UnmarshalBinary(b []byte) error {
	var res PostDomainsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
PostDomainsCreatedBody post domains created body
swagger:model PostDomainsCreatedBody
*/
type PostDomainsCreatedBody struct {

	// domain
	Domain *models.Domain `json:"domain,omitempty"`
}

// Validate validates this post domains created body
func (o *PostDomainsCreatedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomain(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDomainsCreatedBody) validateDomain(formats strfmt.Registry) error {
	if swag.IsZero(o.Domain) { // not required
		return nil
	}

	if o.Domain != nil {
		if err := o.Domain.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postDomainsCreated" + "." + "domain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postDomainsCreated" + "." + "domain")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post domains created body based on the context it is used
func (o *PostDomainsCreatedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDomain(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostDomainsCreatedBody) contextValidateDomain(ctx context.Context, formats strfmt.Registry) error {

	if o.Domain != nil {
		if err := o.Domain.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postDomainsCreated" + "." + "domain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postDomainsCreated" + "." + "domain")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostDomainsCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDomainsCreatedBody) UnmarshalBinary(b []byte) error {
	var res PostDomainsCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
