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

package members

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

	"github.com/sapcc/andromeda/models"
)

// GetPoolsPoolIDMembersMemberIDReader is a Reader for the GetPoolsPoolIDMembersMemberID structure.
type GetPoolsPoolIDMembersMemberIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPoolsPoolIDMembersMemberIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPoolsPoolIDMembersMemberIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetPoolsPoolIDMembersMemberIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetPoolsPoolIDMembersMemberIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetPoolsPoolIDMembersMemberIDOK creates a GetPoolsPoolIDMembersMemberIDOK with default headers values
func NewGetPoolsPoolIDMembersMemberIDOK() *GetPoolsPoolIDMembersMemberIDOK {
	return &GetPoolsPoolIDMembersMemberIDOK{}
}

/* GetPoolsPoolIDMembersMemberIDOK describes a response with status code 200, with default header values.

Shows the details of a specific member.
*/
type GetPoolsPoolIDMembersMemberIDOK struct {
	Payload *GetPoolsPoolIDMembersMemberIDOKBody
}

func (o *GetPoolsPoolIDMembersMemberIDOK) Error() string {
	return fmt.Sprintf("[GET /pools/{pool_id}/members/{member_id}][%d] getPoolsPoolIdMembersMemberIdOK  %+v", 200, o.Payload)
}
func (o *GetPoolsPoolIDMembersMemberIDOK) GetPayload() *GetPoolsPoolIDMembersMemberIDOKBody {
	return o.Payload
}

func (o *GetPoolsPoolIDMembersMemberIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetPoolsPoolIDMembersMemberIDOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPoolsPoolIDMembersMemberIDNotFound creates a GetPoolsPoolIDMembersMemberIDNotFound with default headers values
func NewGetPoolsPoolIDMembersMemberIDNotFound() *GetPoolsPoolIDMembersMemberIDNotFound {
	return &GetPoolsPoolIDMembersMemberIDNotFound{}
}

/* GetPoolsPoolIDMembersMemberIDNotFound describes a response with status code 404, with default header values.

Not Found
*/
type GetPoolsPoolIDMembersMemberIDNotFound struct {
	Payload *models.Error
}

func (o *GetPoolsPoolIDMembersMemberIDNotFound) Error() string {
	return fmt.Sprintf("[GET /pools/{pool_id}/members/{member_id}][%d] getPoolsPoolIdMembersMemberIdNotFound  %+v", 404, o.Payload)
}
func (o *GetPoolsPoolIDMembersMemberIDNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetPoolsPoolIDMembersMemberIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPoolsPoolIDMembersMemberIDDefault creates a GetPoolsPoolIDMembersMemberIDDefault with default headers values
func NewGetPoolsPoolIDMembersMemberIDDefault(code int) *GetPoolsPoolIDMembersMemberIDDefault {
	return &GetPoolsPoolIDMembersMemberIDDefault{
		_statusCode: code,
	}
}

/* GetPoolsPoolIDMembersMemberIDDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type GetPoolsPoolIDMembersMemberIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get pools pool ID members member ID default response
func (o *GetPoolsPoolIDMembersMemberIDDefault) Code() int {
	return o._statusCode
}

func (o *GetPoolsPoolIDMembersMemberIDDefault) Error() string {
	return fmt.Sprintf("[GET /pools/{pool_id}/members/{member_id}][%d] GetPoolsPoolIDMembersMemberID default  %+v", o._statusCode, o.Payload)
}
func (o *GetPoolsPoolIDMembersMemberIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetPoolsPoolIDMembersMemberIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetPoolsPoolIDMembersMemberIDOKBody get pools pool ID members member ID o k body
swagger:model GetPoolsPoolIDMembersMemberIDOKBody
*/
type GetPoolsPoolIDMembersMemberIDOKBody struct {

	// member
	Member *models.Member `json:"member,omitempty"`
}

// Validate validates this get pools pool ID members member ID o k body
func (o *GetPoolsPoolIDMembersMemberIDOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMember(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetPoolsPoolIDMembersMemberIDOKBody) validateMember(formats strfmt.Registry) error {
	if swag.IsZero(o.Member) { // not required
		return nil
	}

	if o.Member != nil {
		if err := o.Member.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getPoolsPoolIdMembersMemberIdOK" + "." + "member")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getPoolsPoolIdMembersMemberIdOK" + "." + "member")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get pools pool ID members member ID o k body based on the context it is used
func (o *GetPoolsPoolIDMembersMemberIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateMember(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetPoolsPoolIDMembersMemberIDOKBody) contextValidateMember(ctx context.Context, formats strfmt.Registry) error {

	if o.Member != nil {
		if err := o.Member.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getPoolsPoolIdMembersMemberIdOK" + "." + "member")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getPoolsPoolIdMembersMemberIdOK" + "." + "member")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetPoolsPoolIDMembersMemberIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPoolsPoolIDMembersMemberIDOKBody) UnmarshalBinary(b []byte) error {
	var res GetPoolsPoolIDMembersMemberIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
