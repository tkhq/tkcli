// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// V1ActivityStatus The current processing status of an Activity.
//
// swagger:model v1ActivityStatus
type V1ActivityStatus string

func NewV1ActivityStatus(value V1ActivityStatus) *V1ActivityStatus {
	return &value
}

// Pointer returns a pointer to a freshly-allocated V1ActivityStatus.
func (m V1ActivityStatus) Pointer() *V1ActivityStatus {
	return &m
}

const (

	// V1ActivityStatusACTIVITYSTATUSCREATED captures enum value "ACTIVITY_STATUS_CREATED"
	V1ActivityStatusACTIVITYSTATUSCREATED V1ActivityStatus = "ACTIVITY_STATUS_CREATED"

	// V1ActivityStatusACTIVITYSTATUSPENDING captures enum value "ACTIVITY_STATUS_PENDING"
	V1ActivityStatusACTIVITYSTATUSPENDING V1ActivityStatus = "ACTIVITY_STATUS_PENDING"

	// V1ActivityStatusACTIVITYSTATUSCOMPLETED captures enum value "ACTIVITY_STATUS_COMPLETED"
	V1ActivityStatusACTIVITYSTATUSCOMPLETED V1ActivityStatus = "ACTIVITY_STATUS_COMPLETED"

	// V1ActivityStatusACTIVITYSTATUSFAILED captures enum value "ACTIVITY_STATUS_FAILED"
	V1ActivityStatusACTIVITYSTATUSFAILED V1ActivityStatus = "ACTIVITY_STATUS_FAILED"

	// V1ActivityStatusACTIVITYSTATUSCONSENSUSNEEDED captures enum value "ACTIVITY_STATUS_CONSENSUS_NEEDED"
	V1ActivityStatusACTIVITYSTATUSCONSENSUSNEEDED V1ActivityStatus = "ACTIVITY_STATUS_CONSENSUS_NEEDED"

	// V1ActivityStatusACTIVITYSTATUSREJECTED captures enum value "ACTIVITY_STATUS_REJECTED"
	V1ActivityStatusACTIVITYSTATUSREJECTED V1ActivityStatus = "ACTIVITY_STATUS_REJECTED"
)

// for schema
var v1ActivityStatusEnum []interface{}

func init() {
	var res []V1ActivityStatus
	if err := json.Unmarshal([]byte(`["ACTIVITY_STATUS_CREATED","ACTIVITY_STATUS_PENDING","ACTIVITY_STATUS_COMPLETED","ACTIVITY_STATUS_FAILED","ACTIVITY_STATUS_CONSENSUS_NEEDED","ACTIVITY_STATUS_REJECTED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		v1ActivityStatusEnum = append(v1ActivityStatusEnum, v)
	}
}

func (m V1ActivityStatus) validateV1ActivityStatusEnum(path, location string, value V1ActivityStatus) error {
	if err := validate.EnumCase(path, location, value, v1ActivityStatusEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this v1 activity status
func (m V1ActivityStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateV1ActivityStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this v1 activity status based on context it is used
func (m V1ActivityStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}