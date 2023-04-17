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

// Externaldatav1AuthenticatorTransport externaldatav1 authenticator transport
//
// swagger:model externaldatav1AuthenticatorTransport
type Externaldatav1AuthenticatorTransport string

func NewExternaldatav1AuthenticatorTransport(value Externaldatav1AuthenticatorTransport) *Externaldatav1AuthenticatorTransport {
	return &value
}

// Pointer returns a pointer to a freshly-allocated Externaldatav1AuthenticatorTransport.
func (m Externaldatav1AuthenticatorTransport) Pointer() *Externaldatav1AuthenticatorTransport {
	return &m
}

const (

	// Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTBLE captures enum value "AUTHENTICATOR_TRANSPORT_BLE"
	Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTBLE Externaldatav1AuthenticatorTransport = "AUTHENTICATOR_TRANSPORT_BLE"

	// Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTINTERNAL captures enum value "AUTHENTICATOR_TRANSPORT_INTERNAL"
	Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTINTERNAL Externaldatav1AuthenticatorTransport = "AUTHENTICATOR_TRANSPORT_INTERNAL"

	// Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTNFC captures enum value "AUTHENTICATOR_TRANSPORT_NFC"
	Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTNFC Externaldatav1AuthenticatorTransport = "AUTHENTICATOR_TRANSPORT_NFC"

	// Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTUSB captures enum value "AUTHENTICATOR_TRANSPORT_USB"
	Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTUSB Externaldatav1AuthenticatorTransport = "AUTHENTICATOR_TRANSPORT_USB"

	// Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTHYBRID captures enum value "AUTHENTICATOR_TRANSPORT_HYBRID"
	Externaldatav1AuthenticatorTransportAUTHENTICATORTRANSPORTHYBRID Externaldatav1AuthenticatorTransport = "AUTHENTICATOR_TRANSPORT_HYBRID"
)

// for schema
var externaldatav1AuthenticatorTransportEnum []interface{}

func init() {
	var res []Externaldatav1AuthenticatorTransport
	if err := json.Unmarshal([]byte(`["AUTHENTICATOR_TRANSPORT_BLE","AUTHENTICATOR_TRANSPORT_INTERNAL","AUTHENTICATOR_TRANSPORT_NFC","AUTHENTICATOR_TRANSPORT_USB","AUTHENTICATOR_TRANSPORT_HYBRID"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		externaldatav1AuthenticatorTransportEnum = append(externaldatav1AuthenticatorTransportEnum, v)
	}
}

func (m Externaldatav1AuthenticatorTransport) validateExternaldatav1AuthenticatorTransportEnum(path, location string, value Externaldatav1AuthenticatorTransport) error {
	if err := validate.EnumCase(path, location, value, externaldatav1AuthenticatorTransportEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this externaldatav1 authenticator transport
func (m Externaldatav1AuthenticatorTransport) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateExternaldatav1AuthenticatorTransportEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this externaldatav1 authenticator transport based on context it is used
func (m Externaldatav1AuthenticatorTransport) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}