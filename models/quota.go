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

// Quota quota
//
// swagger:model quota
type Quota struct {

	// The configured datacenter quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited.
	// Example: 5
	Datacenter *int64 `json:"datacenter,omitempty"`

	// The configured domain quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited.
	// Example: 5
	Domain *int64 `json:"domain,omitempty"`

	// The configured member quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited.
	// Example: 5
	Member *int64 `json:"member,omitempty"`

	// The configured monitor quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited.
	// Example: 5
	Monitor *int64 `json:"monitor,omitempty"`

	// The configured pool quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited.
	// Example: 5
	Pool *int64 `json:"pool,omitempty"`
}

// Validate validates this quota
func (m *Quota) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this quota based on context it is used
func (m *Quota) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Quota) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Quota) UnmarshalBinary(b []byte) error {
	var res Quota
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}