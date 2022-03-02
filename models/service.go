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

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Service service
//
// swagger:model service
type Service struct {

	// The UTC date and timestamp when had the last heartbeat.
	// Example: 2020-05-11 17:21:34
	Heartbeat string `json:"heartbeat,omitempty"`

	// Hostname of the computer the service is running.
	// Example: example.host
	// Format: hostname
	Host strfmt.Hostname `json:"host,omitempty"`

	// ID of the RPC service.
	// Example: andromeda-agent-fbb49979-03f5-4a97-a334-1fd2c9f61e7e
	ID string `json:"id,omitempty"`

	// metadata
	Metadata interface{} `json:"metadata,omitempty"`

	// Provider this service supports.
	// Example: akamai
	Provider string `json:"provider,omitempty"`

	// RPC Endpoint Address.
	// Example: _INBOX.VEfFxcAzZQ9iM9vwGH49It
	RPCAddress string `json:"rpc_address,omitempty"`

	// Type of service.
	// Example: healthcheck
	Type string `json:"type,omitempty"`

	// Version of the service.
	// Example: 1.2.3
	Version string `json:"version,omitempty"`
}

// Validate validates this service
func (m *Service) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHost(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Service) validateHost(formats strfmt.Registry) error {
	if swag.IsZero(m.Host) { // not required
		return nil
	}

	if err := validate.FormatOf("host", "body", "hostname", m.Host.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this service based on context it is used
func (m *Service) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Service) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Service) UnmarshalBinary(b []byte) error {
	var res Service
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
