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

// V1SignRawPayloadResult v1 sign raw payload result
//
// swagger:model v1SignRawPayloadResult
type V1SignRawPayloadResult struct {

	// Component of an ECSDA signature.
	// Required: true
	R *string `json:"r"`

	// Component of an ECSDA signature.
	// Required: true
	S *string `json:"s"`

	// Component of an ECSDA signature.
	// Required: true
	V *string `json:"v"`
}

// Validate validates this v1 sign raw payload result
func (m *V1SignRawPayloadResult) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateR(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateS(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateV(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1SignRawPayloadResult) validateR(formats strfmt.Registry) error {

	if err := validate.Required("r", "body", m.R); err != nil {
		return err
	}

	return nil
}

func (m *V1SignRawPayloadResult) validateS(formats strfmt.Registry) error {

	if err := validate.Required("s", "body", m.S); err != nil {
		return err
	}

	return nil
}

func (m *V1SignRawPayloadResult) validateV(formats strfmt.Registry) error {

	if err := validate.Required("v", "body", m.V); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this v1 sign raw payload result based on context it is used
func (m *V1SignRawPayloadResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1SignRawPayloadResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1SignRawPayloadResult) UnmarshalBinary(b []byte) error {
	var res V1SignRawPayloadResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}