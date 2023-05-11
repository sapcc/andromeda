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

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// QuotaUsage quota usage
//
// swagger:model quota_usage
type QuotaUsage struct {

	// The current quota usage of datacenter.
	// Example: 5
	InUseDatacenter int64 `json:"in_use_datacenter"`

	// The current quota usage of domain.
	// Example: 5
	InUseDomain int64 `json:"in_use_domain"`

	// The current quota usage of member.
	// Example: 5
	InUseMember int64 `json:"in_use_member"`

	// The current quota usage of monitor.
	// Example: 5
	InUseMonitor int64 `json:"in_use_monitor"`

	// The current quota usage of pool.
	// Example: 5
	InUsePool int64 `json:"in_use_pool"`
}

// Validate validates this quota usage
func (m *QuotaUsage) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this quota usage based on context it is used
func (m *QuotaUsage) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *QuotaUsage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *QuotaUsage) UnmarshalBinary(b []byte) error {
	var res QuotaUsage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}