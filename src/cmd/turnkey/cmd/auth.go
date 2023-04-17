package cmd

import (
	"github.com/tkhq/tkcli/src/internal/apikey"
	"github.com/tkhq/tkcli/src/internal/clifs"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

// APIKeypair is the loaded API Keypair
var APIKeypair *apikey.ApiKey

// LoadKeypair require-loads the keypair referenced by the given name or as referenced form the global KeyName variable, if name is empty.
func LoadKeypair(name string) {
	if name == "" {
		name = KeyName
	}

	apiKey, err := clifs.LoadKeypair(name)
	if err != nil {
		OutputError(err)
	}

	if apiKey == nil {
		OutputError(errors.New("API key not loaded"))
	}

	APIKeypair = apiKey

	// If we haven't had the organization explicitly setm try to load it from key metadata.
	if Organization == "" {
		// Add the first (and for now only) org in the key metadata
		for _, o := range APIKeypair.Organizations {
			Organization = o

			break
		}
	}

	// If org is _still_ empty, the API key is not usable.
	if Organization == "" {
		OutputError(errors.New("failed to associate the API key with an organization.  Please manually specify the organization ID."))
	}
}

// Authenticator provides a runtime.ClientAuthInfoWriter for use with the swagger API client.
type Authenticator struct {
	// Key optionally overrides the globally-parsed APIKeypair with a custom key.
	Key *apikey.ApiKey
}

// AuthenticateRequest implements runtime.ClientAuthInfoWriter.
// It adds the X-Stamp header to the request based by generating the Stamp with the request body and API key.
func (auth *Authenticator) AuthenticateRequest(req runtime.ClientRequest, reg strfmt.Registry) (err error) {
	authKey := APIKeypair

	if auth.Key != nil {
		authKey = auth.Key
	}

	stamp, err := apikey.Stamp(req.GetBody(), authKey)
	if err != nil {
		return errors.Wrap(err, "failed to generate API stamp")
	}

	return req.SetHeaderParam("X-Stamp", stamp)
}
