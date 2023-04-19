package cmd

import (
	"github.com/tkhq/go-sdk"
	"github.com/tkhq/go-sdk/pkg/apikey"

	"github.com/pkg/errors"
)

// APIKeypair is the loaded API Keypair
var APIKeypair *apikey.Key

// APIClient is the API Client
var APIClient *sdk.Client

// LoadKeypair require-loads the keypair referenced by the given name or as referenced form the global KeyName variable, if name is empty.
func LoadKeypair(name string) {
	if name == "" {
		name = KeyName
	}

	apiKey, err := keyStore.Load(name)
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

func LoadClient() {
	var err error

	APIClient, err = sdk.New(APIKeypair.Name)

	if err != nil {
		OutputError(errors.Wrap(err, "failed to create API client"))
	}
}
