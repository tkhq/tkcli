// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/tkhq/tkcli/src/api/models"
)

// PublicAPIServiceGetUsersReader is a Reader for the PublicAPIServiceGetUsers structure.
type PublicAPIServiceGetUsersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PublicAPIServiceGetUsersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPublicAPIServiceGetUsersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 403:
		result := NewPublicAPIServiceGetUsersForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPublicAPIServiceGetUsersNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPublicAPIServiceGetUsersDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPublicAPIServiceGetUsersOK creates a PublicAPIServiceGetUsersOK with default headers values
func NewPublicAPIServiceGetUsersOK() *PublicAPIServiceGetUsersOK {
	return &PublicAPIServiceGetUsersOK{}
}

/*
PublicAPIServiceGetUsersOK describes a response with status code 200, with default header values.

A successful response.
*/
type PublicAPIServiceGetUsersOK struct {
	Payload *models.V1GetUsersResponse
}

// IsSuccess returns true when this public Api service get users o k response has a 2xx status code
func (o *PublicAPIServiceGetUsersOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this public Api service get users o k response has a 3xx status code
func (o *PublicAPIServiceGetUsersOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this public Api service get users o k response has a 4xx status code
func (o *PublicAPIServiceGetUsersOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this public Api service get users o k response has a 5xx status code
func (o *PublicAPIServiceGetUsersOK) IsServerError() bool {
	return false
}

// IsCode returns true when this public Api service get users o k response a status code equal to that given
func (o *PublicAPIServiceGetUsersOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the public Api service get users o k response
func (o *PublicAPIServiceGetUsersOK) Code() int {
	return 200
}

func (o *PublicAPIServiceGetUsersOK) Error() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] publicApiServiceGetUsersOK  %+v", 200, o.Payload)
}

func (o *PublicAPIServiceGetUsersOK) String() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] publicApiServiceGetUsersOK  %+v", 200, o.Payload)
}

func (o *PublicAPIServiceGetUsersOK) GetPayload() *models.V1GetUsersResponse {
	return o.Payload
}

func (o *PublicAPIServiceGetUsersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1GetUsersResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPublicAPIServiceGetUsersForbidden creates a PublicAPIServiceGetUsersForbidden with default headers values
func NewPublicAPIServiceGetUsersForbidden() *PublicAPIServiceGetUsersForbidden {
	return &PublicAPIServiceGetUsersForbidden{}
}

/*
PublicAPIServiceGetUsersForbidden describes a response with status code 403, with default header values.

Returned when the user does not have permission to access the resource.
*/
type PublicAPIServiceGetUsersForbidden struct {
	Payload interface{}
}

// IsSuccess returns true when this public Api service get users forbidden response has a 2xx status code
func (o *PublicAPIServiceGetUsersForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this public Api service get users forbidden response has a 3xx status code
func (o *PublicAPIServiceGetUsersForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this public Api service get users forbidden response has a 4xx status code
func (o *PublicAPIServiceGetUsersForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this public Api service get users forbidden response has a 5xx status code
func (o *PublicAPIServiceGetUsersForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this public Api service get users forbidden response a status code equal to that given
func (o *PublicAPIServiceGetUsersForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the public Api service get users forbidden response
func (o *PublicAPIServiceGetUsersForbidden) Code() int {
	return 403
}

func (o *PublicAPIServiceGetUsersForbidden) Error() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] publicApiServiceGetUsersForbidden  %+v", 403, o.Payload)
}

func (o *PublicAPIServiceGetUsersForbidden) String() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] publicApiServiceGetUsersForbidden  %+v", 403, o.Payload)
}

func (o *PublicAPIServiceGetUsersForbidden) GetPayload() interface{} {
	return o.Payload
}

func (o *PublicAPIServiceGetUsersForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPublicAPIServiceGetUsersNotFound creates a PublicAPIServiceGetUsersNotFound with default headers values
func NewPublicAPIServiceGetUsersNotFound() *PublicAPIServiceGetUsersNotFound {
	return &PublicAPIServiceGetUsersNotFound{}
}

/*
PublicAPIServiceGetUsersNotFound describes a response with status code 404, with default header values.

Returned when the resource does not exist.
*/
type PublicAPIServiceGetUsersNotFound struct {
	Payload string
}

// IsSuccess returns true when this public Api service get users not found response has a 2xx status code
func (o *PublicAPIServiceGetUsersNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this public Api service get users not found response has a 3xx status code
func (o *PublicAPIServiceGetUsersNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this public Api service get users not found response has a 4xx status code
func (o *PublicAPIServiceGetUsersNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this public Api service get users not found response has a 5xx status code
func (o *PublicAPIServiceGetUsersNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this public Api service get users not found response a status code equal to that given
func (o *PublicAPIServiceGetUsersNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the public Api service get users not found response
func (o *PublicAPIServiceGetUsersNotFound) Code() int {
	return 404
}

func (o *PublicAPIServiceGetUsersNotFound) Error() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] publicApiServiceGetUsersNotFound  %+v", 404, o.Payload)
}

func (o *PublicAPIServiceGetUsersNotFound) String() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] publicApiServiceGetUsersNotFound  %+v", 404, o.Payload)
}

func (o *PublicAPIServiceGetUsersNotFound) GetPayload() string {
	return o.Payload
}

func (o *PublicAPIServiceGetUsersNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPublicAPIServiceGetUsersDefault creates a PublicAPIServiceGetUsersDefault with default headers values
func NewPublicAPIServiceGetUsersDefault(code int) *PublicAPIServiceGetUsersDefault {
	return &PublicAPIServiceGetUsersDefault{
		_statusCode: code,
	}
}

/*
PublicAPIServiceGetUsersDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type PublicAPIServiceGetUsersDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this public Api service get users default response has a 2xx status code
func (o *PublicAPIServiceGetUsersDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this public Api service get users default response has a 3xx status code
func (o *PublicAPIServiceGetUsersDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this public Api service get users default response has a 4xx status code
func (o *PublicAPIServiceGetUsersDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this public Api service get users default response has a 5xx status code
func (o *PublicAPIServiceGetUsersDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this public Api service get users default response a status code equal to that given
func (o *PublicAPIServiceGetUsersDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the public Api service get users default response
func (o *PublicAPIServiceGetUsersDefault) Code() int {
	return o._statusCode
}

func (o *PublicAPIServiceGetUsersDefault) Error() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] PublicApiService_GetUsers default  %+v", o._statusCode, o.Payload)
}

func (o *PublicAPIServiceGetUsersDefault) String() string {
	return fmt.Sprintf("[POST /public/v1/query/list_users][%d] PublicApiService_GetUsers default  %+v", o._statusCode, o.Payload)
}

func (o *PublicAPIServiceGetUsersDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *PublicAPIServiceGetUsersDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}