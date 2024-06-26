package pkg

import (
	"regexp"

	"github.com/rotisserie/eris"

	"github.com/tkhq/go-sdk"
	"github.com/tkhq/go-sdk/pkg/api/client"
	"github.com/tkhq/go-sdk/pkg/apikey"
)

// APIKeypair is the loaded API Keypair.
var APIKeypair *apikey.Key

// APIClient is the API Client.
var APIClient *sdk.Client

// LoadKeypair require-loads the keypair referenced by the given name or as referenced form the global KeyName variable, if name is empty.
func LoadKeypair(name string) {
	if name == "" {
		name = ApiKeyName
	}

	if apiKeyStore == nil {
		OutputError(eris.New("keystore not loaded"))
	}

	apiKey, err := apiKeyStore.Load(name)
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
	scheme := "https"
	if pattern := regexp.MustCompile(`^localhost:\d+$`); pattern.MatchString(apiHost) {
		scheme = "http"
	}
	transportConfig := client.DefaultTransportConfig().WithHost(apiHost).WithSchemes([]string{scheme})

	APIClient = &sdk.Client{
		Client:        client.NewHTTPClientWithConfig(nil, transportConfig),
		Authenticator: &sdk.Authenticator{Key: APIKeypair},
		APIKey:        APIKeypair,
	}
}
