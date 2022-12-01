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
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/models"
)

// GetPoolsReader is a Reader for the GetPools structure.
type GetPoolsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPoolsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPoolsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetPoolsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetPoolsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetPoolsOK creates a GetPoolsOK with default headers values
func NewGetPoolsOK() *GetPoolsOK {
	return &GetPoolsOK{}
}

/*
GetPoolsOK describes a response with status code 200, with default header values.

A JSON array of pools
*/
type GetPoolsOK struct {
	Payload *GetPoolsOKBody
}

// IsSuccess returns true when this get pools o k response has a 2xx status code
func (o *GetPoolsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get pools o k response has a 3xx status code
func (o *GetPoolsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get pools o k response has a 4xx status code
func (o *GetPoolsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get pools o k response has a 5xx status code
func (o *GetPoolsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get pools o k response a status code equal to that given
func (o *GetPoolsOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetPoolsOK) Error() string {
	return fmt.Sprintf("[GET /pools][%d] getPoolsOK  %+v", 200, o.Payload)
}

func (o *GetPoolsOK) String() string {
	return fmt.Sprintf("[GET /pools][%d] getPoolsOK  %+v", 200, o.Payload)
}

func (o *GetPoolsOK) GetPayload() *GetPoolsOKBody {
	return o.Payload
}

func (o *GetPoolsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetPoolsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPoolsBadRequest creates a GetPoolsBadRequest with default headers values
func NewGetPoolsBadRequest() *GetPoolsBadRequest {
	return &GetPoolsBadRequest{}
}

/*
GetPoolsBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type GetPoolsBadRequest struct {
	Payload *models.Error
}

// IsSuccess returns true when this get pools bad request response has a 2xx status code
func (o *GetPoolsBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get pools bad request response has a 3xx status code
func (o *GetPoolsBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get pools bad request response has a 4xx status code
func (o *GetPoolsBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get pools bad request response has a 5xx status code
func (o *GetPoolsBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get pools bad request response a status code equal to that given
func (o *GetPoolsBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *GetPoolsBadRequest) Error() string {
	return fmt.Sprintf("[GET /pools][%d] getPoolsBadRequest  %+v", 400, o.Payload)
}

func (o *GetPoolsBadRequest) String() string {
	return fmt.Sprintf("[GET /pools][%d] getPoolsBadRequest  %+v", 400, o.Payload)
}

func (o *GetPoolsBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetPoolsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPoolsDefault creates a GetPoolsDefault with default headers values
func NewGetPoolsDefault(code int) *GetPoolsDefault {
	return &GetPoolsDefault{
		_statusCode: code,
	}
}

/*
GetPoolsDefault describes a response with status code -1, with default header values.

Unexpected Error
*/
type GetPoolsDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get pools default response
func (o *GetPoolsDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this get pools default response has a 2xx status code
func (o *GetPoolsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get pools default response has a 3xx status code
func (o *GetPoolsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get pools default response has a 4xx status code
func (o *GetPoolsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get pools default response has a 5xx status code
func (o *GetPoolsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get pools default response a status code equal to that given
func (o *GetPoolsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *GetPoolsDefault) Error() string {
	return fmt.Sprintf("[GET /pools][%d] GetPools default  %+v", o._statusCode, o.Payload)
}

func (o *GetPoolsDefault) String() string {
	return fmt.Sprintf("[GET /pools][%d] GetPools default  %+v", o._statusCode, o.Payload)
}

func (o *GetPoolsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetPoolsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetPoolsOKBody get pools o k body
swagger:model GetPoolsOKBody
*/
type GetPoolsOKBody struct {

	// links
	Links []*models.Link `json:"links,omitempty"`

	// pools
	Pools []*models.Pool `json:"pools"`
}

// Validate validates this get pools o k body
func (o *GetPoolsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePools(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetPoolsOKBody) validateLinks(formats strfmt.Registry) error {
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
					return ve.ValidateName("getPoolsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getPoolsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetPoolsOKBody) validatePools(formats strfmt.Registry) error {
	if swag.IsZero(o.Pools) { // not required
		return nil
	}

	for i := 0; i < len(o.Pools); i++ {
		if swag.IsZero(o.Pools[i]) { // not required
			continue
		}

		if o.Pools[i] != nil {
			if err := o.Pools[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getPoolsOK" + "." + "pools" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getPoolsOK" + "." + "pools" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get pools o k body based on the context it is used
func (o *GetPoolsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidatePools(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetPoolsOKBody) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Links); i++ {

		if o.Links[i] != nil {
			if err := o.Links[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getPoolsOK" + "." + "links" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getPoolsOK" + "." + "links" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetPoolsOKBody) contextValidatePools(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Pools); i++ {

		if o.Pools[i] != nil {
			if err := o.Pools[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getPoolsOK" + "." + "pools" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getPoolsOK" + "." + "pools" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetPoolsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPoolsOKBody) UnmarshalBinary(b []byte) error {
	var res GetPoolsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
