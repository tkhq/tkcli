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

// V1APIKey v1 Api key
//
// swagger:model v1ApiKey
type V1APIKey struct {

	// Unique identifier for a given API Key.
	// Required: true
	APIKeyID *string `json:"apiKeyId"`

	// Human-readable name for an API Key.
	// Required: true
	APIKeyName *string `json:"apiKeyName"`

	// created at
	// Required: true
	CreatedAt *V1Timestamp `json:"createdAt"`

	// credential
	// Required: true
	Credential *V1Credential `json:"credential"`

	// updated at
	// Required: true
	UpdatedAt *V1Timestamp `json:"updatedAt"`
}

// Validate validates this v1 Api key
func (m *V1APIKey) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAPIKeyID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAPIKeyName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCredential(formats); err != nil {
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

func (m *V1APIKey) validateAPIKeyID(formats strfmt.Registry) error {

	if err := validate.Required("apiKeyId", "body", m.APIKeyID); err != nil {
		return err
	}

	return nil
}

func (m *V1APIKey) validateAPIKeyName(formats strfmt.Registry) error {

	if err := validate.Required("apiKeyName", "body", m.APIKeyName); err != nil {
		return err
	}

	return nil
}

func (m *V1APIKey) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("createdAt", "body", m.CreatedAt); err != nil {
		return err
	}

	if m.CreatedAt != nil {
		if err := m.CreatedAt.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("createdAt")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("createdAt")
			}
			return err
		}
	}

	return nil
}

func (m *V1APIKey) validateCredential(formats strfmt.Registry) error {

	if err := validate.Required("credential", "body", m.Credential); err != nil {
		return err
	}

	if m.Credential != nil {
		if err := m.Credential.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("credential")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("credential")
			}
			return err
		}
	}

	return nil
}

func (m *V1APIKey) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", m.UpdatedAt); err != nil {
		return err
	}

	if m.UpdatedAt != nil {
		if err := m.UpdatedAt.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("updatedAt")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("updatedAt")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1 Api key based on the context it is used
func (m *V1APIKey) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCredential(ctx, formats); err != nil {
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

func (m *V1APIKey) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if m.CreatedAt != nil {
		if err := m.CreatedAt.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("createdAt")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("createdAt")
			}
			return err
		}
	}

	return nil
}

func (m *V1APIKey) contextValidateCredential(ctx context.Context, formats strfmt.Registry) error {

	if m.Credential != nil {
		if err := m.Credential.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("credential")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("credential")
			}
			return err
		}
	}

	return nil
}

func (m *V1APIKey) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if m.UpdatedAt != nil {
		if err := m.UpdatedAt.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("updatedAt")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("updatedAt")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1APIKey) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1APIKey) UnmarshalBinary(b []byte) error {
	var res V1APIKey
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}