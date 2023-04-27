package cmd

import (
	"github.com/rotisserie/eris"

	"github.com/tkhq/go-sdk"
	"github.com/tkhq/go-sdk/pkg/apikey"
)

// APIKeypair is the loaded API Keypair.
var APIKeypair *apikey.Key

// APIClient is the API Client.
var APIClient *sdk.Client

// LoadKeypair require-loads the keypair referenced by the given name or as referenced form the global KeyName variable, if name is empty.
func LoadKeypair(name string) {
	if name == "" {
		name = KeyName
	}

	if keyStore == nil {
		OutputError(eris.New("keystore not loaded"))
	}

	apiKey, err := keyStore.Load(name)
	if err != nil {
		OutputError(err)
	}

	if apiKey == nil {
		OutputError(eris.New("API key not loaded"))
	}

	APIKeypair = apiKey

	// If we haven't had the organization explicitly set try to load it from key metadata.
	if Organization == "" {
		// Add the first (and for now only) org in the key metadata
		for _, o := range APIKeypair.Organizations {
			Organization = o

			break
		}
	}

	// If org is _still_ empty, the API key is not usable.
	if Organization == "" {
		OutputError(eris.New("failed to associate the API key with an organization; please manually specify the organization ID"))
	}
}

// LoadClient creates an API client from the preloaded API keypair.
func LoadClient() {
	var err error

	APIClient, err = sdk.New(APIKeypair.Name)

	if err != nil {
		OutputError(eris.Wrap(err, "failed to create API client"))
	}
}
