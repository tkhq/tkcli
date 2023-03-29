// Code generated by go-swagger; DO NOT EDIT.

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

// V1GetActivityRequest v1 get activity request
//
// swagger:model v1GetActivityRequest
type V1GetActivityRequest struct {

	// Unique identifier for a given Activity object.
	// Required: true
	ActivityID *string `json:"activityId"`

	// Unique identifier for a given Organization.
	// Required: true
	OrganizationID *string `json:"organizationId"`
}

// Validate validates this v1 get activity request
func (m *V1GetActivityRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActivityID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOrganizationID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1GetActivityRequest) validateActivityID(formats strfmt.Registry) error {

	if err := validate.Required("activityId", "body", m.ActivityID); err != nil {
		return err
	}

	return nil
}

func (m *V1GetActivityRequest) validateOrganizationID(formats strfmt.Registry) error {

	if err := validate.Required("organizationId", "body", m.OrganizationID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this v1 get activity request based on context it is used
func (m *V1GetActivityRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1GetActivityRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1GetActivityRequest) UnmarshalBinary(b []byte) error {
	var res V1GetActivityRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}