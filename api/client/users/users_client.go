// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new users API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for users API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	PublicAPIServiceCreateAPIKeys(params *PublicAPIServiceCreateAPIKeysParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceCreateAPIKeysOK, error)

	PublicAPIServiceDeleteAPIKeys(params *PublicAPIServiceDeleteAPIKeysParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceDeleteAPIKeysOK, error)

	PublicAPIServiceGetUser(params *PublicAPIServiceGetUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceGetUserOK, error)

	PublicAPIServiceGetUsers(params *PublicAPIServiceGetUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceGetUsersOK, error)

	PublicAPIServiceGetWhoami(params *PublicAPIServiceGetWhoamiParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceGetWhoamiOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
PublicAPIServiceCreateAPIKeys creates API keys

Add api keys to an existing User
*/
func (a *Client) PublicAPIServiceCreateAPIKeys(params *PublicAPIServiceCreateAPIKeysParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceCreateAPIKeysOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPublicAPIServiceCreateAPIKeysParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PublicApiService_CreateApiKeys",
		Method:             "POST",
		PathPattern:        "/public/v1/submit/create_api_keys",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PublicAPIServiceCreateAPIKeysReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PublicAPIServiceCreateAPIKeysOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PublicAPIServiceCreateAPIKeysDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
PublicAPIServiceDeleteAPIKeys deletes API keys

Remove api keys from a User
*/
func (a *Client) PublicAPIServiceDeleteAPIKeys(params *PublicAPIServiceDeleteAPIKeysParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceDeleteAPIKeysOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPublicAPIServiceDeleteAPIKeysParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PublicApiService_DeleteApiKeys",
		Method:             "POST",
		PathPattern:        "/public/v1/submit/delete_api_keys",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PublicAPIServiceDeleteAPIKeysReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PublicAPIServiceDeleteAPIKeysOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PublicAPIServiceDeleteAPIKeysDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
PublicAPIServiceGetUser gets user

Get details about a User
*/
func (a *Client) PublicAPIServiceGetUser(params *PublicAPIServiceGetUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceGetUserOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPublicAPIServiceGetUserParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PublicApiService_GetUser",
		Method:             "POST",
		PathPattern:        "/public/v1/query/get_user",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PublicAPIServiceGetUserReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PublicAPIServiceGetUserOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PublicAPIServiceGetUserDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
PublicAPIServiceGetUsers lists users

List all Users within an Organization
*/
func (a *Client) PublicAPIServiceGetUsers(params *PublicAPIServiceGetUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceGetUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPublicAPIServiceGetUsersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PublicApiService_GetUsers",
		Method:             "POST",
		PathPattern:        "/public/v1/query/list_users",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PublicAPIServiceGetUsersReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PublicAPIServiceGetUsersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PublicAPIServiceGetUsersDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
PublicAPIServiceGetWhoami whos am i

Get basic information about your current API user and your organization
*/
func (a *Client) PublicAPIServiceGetWhoami(params *PublicAPIServiceGetWhoamiParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PublicAPIServiceGetWhoamiOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPublicAPIServiceGetWhoamiParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PublicApiService_GetWhoami",
		Method:             "POST",
		PathPattern:        "/public/v1/query/whoami",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PublicAPIServiceGetWhoamiReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PublicAPIServiceGetWhoamiOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PublicAPIServiceGetWhoamiDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}