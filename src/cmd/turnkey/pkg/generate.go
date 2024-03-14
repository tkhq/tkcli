package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/apikey"
	"github.com/tkhq/go-sdk/pkg/store/local"
)

func init() {
	generateCmd.AddCommand(apiKeyCmd)

	rootCmd.AddCommand(generateCmd)
}

// generateCmd represents the base command for generating different kinds of keys
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate keys",
}

// Represents the command to generate an API key
var apiKeyCmd = &cobra.Command{
	Use:   "api-key",
	Short: "Generate a Turnkey API key",
	Long:  `Generate a new API key that can be used for authenticating with the API.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Organization == "" {
			OutputError(eris.New("--organization must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("key-name")
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read API key name"))
		}

		apiKey, err := apikey.New(Organization)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to create API keypair"))
		}

		if name == "-" {
			Output(map[string]string{
				"publicKey":  apiKey.TkPublicKey,
				"privateKey": apiKey.TkPrivateKey,
			})
		}

		apiKey.Metadata.Name = name
		apiKey.Metadata.PublicKey = apiKey.TkPublicKey
		apiKey.Metadata.Organizations = []string{Organization}

		if err = keyStore.Store(name, apiKey); err != nil {
			OutputError(eris.Wrap(err, "failed to store new API keypair"))
		}

		localStore, ok := keyStore.(*local.Store)
		if !ok {
			OutputError(eris.Wrap(err, "unhandled keystore type: expected *local.Store"))
		}

		Output(map[string]string{
			"publicKey":      apiKey.TkPublicKey,
			"publicKeyFile":  localStore.PublicKeyFile(name),
			"privateKeyFile": localStore.PrivateKeyFile(name),
		})
	},
}
