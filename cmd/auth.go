package cmd

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/tkhq/tkcli/internal/apikey"
	"github.com/tkhq/tkcli/internal/clifs"
)

// APIKeypair is the loaded API Keypair
var APIKeypair *apikey.ApiKey

// LoadKeypair loads the keypair referenced by the given name or as referenced form the global KeyName variable, if name is empty.
func LoadKeypair(name string) error {
	if name == "" {
		name = KeyName
	}

	apiKey, err := clifs.LoadKeypair(name)
	if err != nil {
		return err
	}

	APIKeypair = apiKey

	return nil
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
