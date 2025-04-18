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

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Link link
//
// swagger:model link
type Link struct {

	// href
	// Format: uri
	Href strfmt.URI `json:"href,omitempty" db:"href,omitempty"`

	// rel
	// Enum: [next previous]
	Rel string `json:"rel,omitempty" db:"rel,omitempty"`
}

// Validate validates this link
func (m *Link) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHref(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRel(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Link) validateHref(formats strfmt.Registry) error {
	if swag.IsZero(m.Href) { // not required
		return nil
	}

	if err := validate.FormatOf("href", "body", "uri", m.Href.String(), formats); err != nil {
		return err
	}

	return nil
}

var linkTypeRelPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["next","previous"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		linkTypeRelPropEnum = append(linkTypeRelPropEnum, v)
	}
}

const (

	// LinkRelNext captures enum value "next"
	LinkRelNext string = "next"

	// LinkRelPrevious captures enum value "previous"
	LinkRelPrevious string = "previous"
)

// prop value enum
func (m *Link) validateRelEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, linkTypeRelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Link) validateRel(formats strfmt.Registry) error {
	if swag.IsZero(m.Rel) { // not required
		return nil
	}

	// value enum
	if err := m.validateRelEnum("rel", "body", m.Rel); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this link based on the context it is used
func (m *Link) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Link) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Link) UnmarshalBinary(b []byte) error {
	var res Link
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
