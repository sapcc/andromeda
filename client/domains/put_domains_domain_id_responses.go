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

// PutDomainsDomainIDReader is a Reader for the PutDomainsDomainID structure.
type PutDomainsDomainIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutDomainsDomainIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 202:
		result := NewPutDomainsDomainIDAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPutDomainsDomainIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPutDomainsDomainIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPutDomainsDomainIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPutDomainsDomainIDAccepted creates a PutDomainsDomainIDAccepted with default headers values
func NewPutDomainsDomainIDAccepted() *PutDomainsDomainIDAccepted {
	return &PutDomainsDomainIDAccepted{}
}

/* PutDomainsDomainIDAccepted describes a response with status code 202, with default header values.

Updated domain.
*/
type PutDomainsDomainIDAccepted struct {
	Payload *PutDomainsDomainIDAcceptedBody
}

func (o *PutDomainsDomainIDAccepted) Error() string {
	return fmt.Sprintf("[PUT /domains/{domain_id}][%d] putDomainsDomainIdAccepted  %+v", 202, o.Payload)
}
func (o *PutDomainsDomainIDAccepted) GetPayload() *PutDomainsDomainIDAcceptedBody {
	return o.Payload
}

func (o *PutDomainsDomainIDAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PutDomainsDomainIDAcceptedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutDomainsDomainIDBadRequest creates a PutDomainsDomainIDBadRequest with default headers values
func NewPutDomainsDomainIDBadRequest() *PutDomainsDomainIDBadRequest {
	return &PutDomainsDomainIDBadRequest{}
}

/* PutDomainsDomainIDBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type PutDomainsDomainIDBadRequest struct {
	Payload *models.Error
}

func (o *PutDomainsDomainIDBadRequest) Error() string {
	return fmt.Sprintf("[PUT /domains/{domain_id}][%d] putDomainsDomainIdBadRequest  %+v", 400, o.Payload)
}
func (o *PutDomainsDomainIDBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *PutDomainsDomainIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutDomainsDomainIDNotFound creates a PutDomainsDomainIDNotFound with default headers values
func NewPutDomainsDomainIDNotFound() *PutDomainsDomainIDNotFound {
	return &PutDomainsDomainIDNotFound{}
}

/* PutDomainsDomainIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type PutDomainsDomainIDNotFound struct {
	Payload *models.Error
}

func (o *PutDomainsDomainIDNotFound) Error() string {
	return fmt.Sprintf("[PUT /domains/{domain_id}][%d] putDomainsDomainIdNotFound  %+v", 404, o.Payload)
}
func (o *PutDomainsDomainIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *PutDomainsDomainIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutDomainsDomainIDDefault creates a PutDomainsDomainIDDefault with default headers values
func NewPutDomainsDomainIDDefault(code int) *PutDomainsDomainIDDefault {
	return &PutDomainsDomainIDDefault{
		_statusCode: code,
	}
}

/* PutDomainsDomainIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type PutDomainsDomainIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the put domains domain ID default response
func (o *PutDomainsDomainIDDefault) Code() int {
	return o._statusCode
}

func (o *PutDomainsDomainIDDefault) Error() string {
	return fmt.Sprintf("[PUT /domains/{domain_id}][%d] PutDomainsDomainID default  %+v", o._statusCode, o.Payload)
}
func (o *PutDomainsDomainIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PutDomainsDomainIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*PutDomainsDomainIDAcceptedBody put domains domain ID accepted body
swagger:model PutDomainsDomainIDAcceptedBody
*/
type PutDomainsDomainIDAcceptedBody struct {

	// domain
	Domain *models.Domain `json:"domain,omitempty"`
}

// Validate validates this put domains domain ID accepted body
func (o *PutDomainsDomainIDAcceptedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomain(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutDomainsDomainIDAcceptedBody) validateDomain(formats strfmt.Registry) error {
	if swag.IsZero(o.Domain) { // not required
		return nil
	}

	if o.Domain != nil {
		if err := o.Domain.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putDomainsDomainIdAccepted" + "." + "domain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putDomainsDomainIdAccepted" + "." + "domain")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this put domains domain ID accepted body based on the context it is used
func (o *PutDomainsDomainIDAcceptedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDomain(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutDomainsDomainIDAcceptedBody) contextValidateDomain(ctx context.Context, formats strfmt.Registry) error {

	if o.Domain != nil {
		if err := o.Domain.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("putDomainsDomainIdAccepted" + "." + "domain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("putDomainsDomainIdAccepted" + "." + "domain")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutDomainsDomainIDAcceptedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutDomainsDomainIDAcceptedBody) UnmarshalBinary(b []byte) error {
	var res PutDomainsDomainIDAcceptedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*PutDomainsDomainIDBody put domains domain ID body
swagger:model PutDomainsDomainIDBody
*/
type PutDomainsDomainIDBody struct {

	// domain
	// Required: true
	Domain *models.Domain `json:"domain"`
}

// Validate validates this put domains domain ID body
func (o *PutDomainsDomainIDBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomain(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutDomainsDomainIDBody) validateDomain(formats strfmt.Registry) error {

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

// ContextValidate validate this put domains domain ID body based on the context it is used
func (o *PutDomainsDomainIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDomain(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutDomainsDomainIDBody) contextValidateDomain(ctx context.Context, formats strfmt.Registry) error {

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
func (o *PutDomainsDomainIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutDomainsDomainIDBody) UnmarshalBinary(b []byte) error {
	var res PutDomainsDomainIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
