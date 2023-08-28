package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/apikey"
	"github.com/tkhq/go-sdk/pkg/store/local"
)

func init() {
	rootCmd.AddCommand(genAPIKeyCmd)
}

var genAPIKeyCmd = &cobra.Command{
	Use:     "generate-api-key",
	Short:   "Generate a Turnkey API key",
	Aliases: []string{"g", "gen"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if Organization == "" {
			OutputError(eris.New("please supply an organization ID (UUID)"))
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
