package cmd

import (
	"github.com/tkhq/tkcli/src/internal/apikey"
	"github.com/tkhq/tkcli/src/internal/clifs"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genApiKey)
}

var genApiKey = &cobra.Command{
	Use:     "generate-api-key generates a Turnkey API key",
	Short:   "generate-api-key generates a Turnkey API key",
	Aliases: []string{"g", "gen"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if Organization == "" {
			OutputError(errors.New("please supply an organization ID (UUID)"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("key-name")
		if err != nil {
			OutputError(errors.Wrap(err, "failed to read API key name"))
		}

		apiKey, err := apikey.NewTkApiKey()
		if err != nil {
			OutputError(errors.Wrap(err, "failed to create API keypair"))
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

		if err = clifs.StoreKeypair(name, apiKey); err != nil {
			OutputError(errors.Wrap(err, "failed to store new API keypair"))
		}

		Output(map[string]string{
			"publicKey":      string(apiKey.TkPublicKey),
			"publicKeyFile":  clifs.PublicKeyFile(name),
			"privateKeyFile": clifs.PrivateKeyFile(name),
		})
	},
}
