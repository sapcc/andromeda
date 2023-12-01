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

// Member member
//
// swagger:model member
type Member struct {

	// Address to use.
	// Example: 1.2.3.4
	// Required: true
	// Format: ipv4
	Address *strfmt.IPv4 `json:"address"`

	// The administrative state of the resource, which is up (true) or down (false). Default is true.
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The UTC date and timestamp when the resource was created.
	// Example: 2020-05-11T17:21:34
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// Datacenter assigned for this member.
	// Format: uuid
	DatacenterID *strfmt.UUID `json:"datacenter_id,omitempty"`

	// The id of the resource.
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// Human-readable name of the resource.
	// Max Length: 255
	Name *string `json:"name,omitempty"`

	// pool id.
	// Format: uuid
	PoolID *strfmt.UUID `json:"pool_id,omitempty"`

	// Port to use for monitor checks.
	// Example: 80
	// Required: true
	// Maximum: 65535
	// Minimum: 0
	Port *int64 `json:"port"`

	// The ID of the project owning this resource.
	// Example: fa84c217f361441986a220edf9b1e337
	// Max Length: 32
	// Min Length: 32
	ProjectID *string `json:"project_id,omitempty"`

	// provisioning status
	// Read Only: true
	// Enum: [PENDING ACTIVE ERROR]
	ProvisioningStatus string `json:"provisioning_status,omitempty"`

	// status
	// Read Only: true
	// Enum: [ONLINE NO_MONITOR DOWN]
	Status string `json:"status,omitempty"`

	// The UTC date and timestamp when the resource was created.
	// Example: 2020-09-09T14:52:15
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updated_at,omitempty"`
}

// Validate validates this member
func (m *Member) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDatacenterID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePoolID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePort(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProjectID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProvisioningStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Member) validateAddress(formats strfmt.Registry) error {

	if err := validate.Required("address", "body", m.Address); err != nil {
		return err
	}

	if err := validate.FormatOf("address", "body", "ipv4", m.Address.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Member) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Member) validateDatacenterID(formats strfmt.Registry) error {
	if swag.IsZero(m.DatacenterID) { // not required
		return nil
	}

	if err := validate.FormatOf("datacenter_id", "body", "uuid", m.DatacenterID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Member) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Member) validateName(formats strfmt.Registry) error {
	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := validate.MaxLength("name", "body", *m.Name, 255); err != nil {
		return err
	}

	return nil
}

func (m *Member) validatePoolID(formats strfmt.Registry) error {
	if swag.IsZero(m.PoolID) { // not required
		return nil
	}

	if err := validate.FormatOf("pool_id", "body", "uuid", m.PoolID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Member) validatePort(formats strfmt.Registry) error {

	if err := validate.Required("port", "body", m.Port); err != nil {
		return err
	}

	if err := validate.MinimumInt("port", "body", *m.Port, 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("port", "body", *m.Port, 65535, false); err != nil {
		return err
	}

	return nil
}

func (m *Member) validateProjectID(formats strfmt.Registry) error {
	if swag.IsZero(m.ProjectID) { // not required
		return nil
	}

	if err := validate.MinLength("project_id", "body", *m.ProjectID, 32); err != nil {
		return err
	}

	if err := validate.MaxLength("project_id", "body", *m.ProjectID, 32); err != nil {
		return err
	}

	return nil
}

var memberTypeProvisioningStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["PENDING","ACTIVE","ERROR"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		memberTypeProvisioningStatusPropEnum = append(memberTypeProvisioningStatusPropEnum, v)
	}
}

const (

	// MemberProvisioningStatusPENDING captures enum value "PENDING"
	MemberProvisioningStatusPENDING string = "PENDING"

	// MemberProvisioningStatusACTIVE captures enum value "ACTIVE"
	MemberProvisioningStatusACTIVE string = "ACTIVE"

	// MemberProvisioningStatusERROR captures enum value "ERROR"
	MemberProvisioningStatusERROR string = "ERROR"
)

// prop value enum
func (m *Member) validateProvisioningStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, memberTypeProvisioningStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Member) validateProvisioningStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.ProvisioningStatus) { // not required
		return nil
	}

	// value enum
	if err := m.validateProvisioningStatusEnum("provisioning_status", "body", m.ProvisioningStatus); err != nil {
		return err
	}

	return nil
}

var memberTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ONLINE","NO_MONITOR","DOWN"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		memberTypeStatusPropEnum = append(memberTypeStatusPropEnum, v)
	}
}

const (

	// MemberStatusONLINE captures enum value "ONLINE"
	MemberStatusONLINE string = "ONLINE"

	// MemberStatusNOMONITOR captures enum value "NO_MONITOR"
	MemberStatusNOMONITOR string = "NO_MONITOR"

	// MemberStatusDOWN captures enum value "DOWN"
	MemberStatusDOWN string = "DOWN"
)

// prop value enum
func (m *Member) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, memberTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Member) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

func (m *Member) validateUpdatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("updated_at", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this member based on the context it is used
func (m *Member) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProvisioningStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUpdatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Member) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "created_at", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *Member) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", strfmt.UUID(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *Member) contextValidateProvisioningStatus(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "provisioning_status", "body", string(m.ProvisioningStatus)); err != nil {
		return err
	}

	return nil
}

func (m *Member) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "status", "body", string(m.Status)); err != nil {
		return err
	}

	return nil
}

func (m *Member) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "updated_at", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Member) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Member) UnmarshalBinary(b []byte) error {
	var res Member
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
