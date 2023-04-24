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
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/models"
)

// GetDomainsReader is a Reader for the GetDomains structure.
type GetDomainsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDomainsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDomainsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetDomainsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetDomainsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDomainsOK creates a GetDomainsOK with default headers values
func NewGetDomainsOK() *GetDomainsOK {
	return &GetDomainsOK{}
}

/*
GetDomainsOK describes a response with status code 200, with default header values.

A JSON array of domains
*/
type GetDomainsOK struct {
	Payload *GetDomainsOKBody
}

// IsSuccess returns true when this get domains o k response has a 2xx status code
func (o *GetDomainsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get domains o k response has a 3xx status code
func (o *GetDomainsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get domains o k response has a 4xx status code
func (o *GetDomainsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get domains o k response has a 5xx status code
func (o *GetDomainsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get domains o k response a status code equal to that given
func (o *GetDomainsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get domains o k response
func (o *GetDomainsOK) Code() int {
	return 200
}

func (o *GetDomainsOK) Error() string {
	return fmt.Sprintf("[GET /domains][%d] getDomainsOK  %+v", 200, o.Payload)
}

func (o *GetDomainsOK) String() string {
	return fmt.Sprintf("[GET /domains][%d] getDomainsOK  %+v", 200, o.Payload)
}

func (o *GetDomainsOK) GetPayload() *GetDomainsOKBody {
	return o.Payload
}

func (o *GetDomainsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetDomainsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDomainsBadRequest creates a GetDomainsBadRequest with default headers values
func NewGetDomainsBadRequest() *GetDomainsBadRequest {
	return &GetDomainsBadRequest{}
}

/*
GetDomainsBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type GetDomainsBadRequest struct {
	Payload *models.Error
}

// IsSuccess returns true when this get domains bad request response has a 2xx status code
func (o *GetDomainsBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get domains bad request response has a 3xx status code
func (o *GetDomainsBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get domains bad request response has a 4xx status code
func (o *GetDomainsBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get domains bad request response has a 5xx status code
func (o *GetDomainsBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get domains bad request response a status code equal to that given
func (o *GetDomainsBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get domains bad request response
func (o *GetDomainsBadRequest) Code() int {
	return 400
}

func (o *GetDomainsBadRequest) Error() string {
	return fmt.Sprintf("[GET /domains][%d] getDomainsBadRequest  %+v", 400, o.Payload)
}

func (o *GetDomainsBadRequest) String() string {
	return fmt.Sprintf("[GET /domains][%d] getDomainsBadRequest  %+v", 400, o.Payload)
}

func (o *GetDomainsBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetDomainsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDomainsDefault creates a GetDomainsDefault with default headers values
func NewGetDomainsDefault(code int) *GetDomainsDefault {
	return &GetDomainsDefault{
		_statusCode: code,
	}
}

/*
GetDomainsDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type GetDomainsDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get domains default response has a 2xx status code
func (o *GetDomainsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get domains default response has a 3xx status code
func (o *GetDomainsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get domains default response has a 4xx status code
func (o *GetDomainsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get domains default response has a 5xx status code
func (o *GetDomainsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get domains default response a status code equal to that given
func (o *GetDomainsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get domains default response
func (o *GetDomainsDefault) Code() int {
	return o._statusCode
}

func (o *GetDomainsDefault) Error() string {
	return fmt.Sprintf("[GET /domains][%d] GetDomains default  %+v", o._statusCode, o.Payload)
}

func (o *GetDomainsDefault) String() string {
	return fmt.Sprintf("[GET /domains][%d] GetDomains default  %+v", o._statusCode, o.Payload)
}

func (o *GetDomainsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetDomainsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetDomainsOKBody get domains o k body
swagger:model GetDomainsOKBody
*/
type GetDomainsOKBody struct {

	// domains
	Domains []*models.Domain `json:"domains"`

	// links
	Links []*models.Link `json:"links,omitempty"`
}

// Validate validates this get domains o k body
func (o *GetDomainsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDomains(formats); err != nil {
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

func (o *GetDomainsOKBody) validateDomains(formats strfmt.Registry) error {
	if swag.IsZero(o.Domains) { // not required
		return nil
	}

	for i := 0; i < len(o.Domains); i++ {
		if swag.IsZero(o.Domains[i]) { // not required
			continue
		}

		if o.Domains[i] != nil {
			if err := o.Domains[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetDomainsOKBody) validateLinks(formats strfmt.Registry) error {
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
					return ve.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get domains o k body based on the context it is used
func (o *GetDomainsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDomains(ctx, formats); err != nil {
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

func (o *GetDomainsOKBody) contextValidateDomains(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Domains); i++ {

		if o.Domains[i] != nil {
			if err := o.Domains[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "domains" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetDomainsOKBody) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Links); i++ {

		if o.Links[i] != nil {
			if err := o.Links[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getDomainsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetDomainsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDomainsOKBody) UnmarshalBinary(b []byte) error {
	var res GetDomainsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
