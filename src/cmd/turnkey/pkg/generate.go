package pkg

import (
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/apikey"
	"github.com/tkhq/go-sdk/pkg/encryptionkey"
	"github.com/tkhq/go-sdk/pkg/store/local"
)

var (
	curveType string
)

func init() {
	apiKeyCmd.Flags().StringVar(&curveType, "curve", "p256", "curve type for API key; p256, secp256k1, and ed25519 currently supported")
	generateCmd.AddCommand(apiKeyCmd)

	encryptionKeyCmd.Flags().StringVar(&User, "user", "", "ID of user to generating the encryption key")
	generateCmd.AddCommand(encryptionKeyCmd)

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

		curveType, err := cmd.Flags().GetString("curve")
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read curve type"))
		}

		var apiKey *apikey.Key

		switch curveType {
		default:
			OutputError(fmt.Errorf("invalid curve type: %s; supported types are p256, secp256k1, and ed25519", curveType))
		case "p256":
			apiKey, err = apikey.New(Organization, apikey.SchemeP256)
			if err != nil {
				OutputError(eris.Wrap(err, "failed to create API keypair"))
			}
		case "secp256k1":
			apiKey, err = apikey.New(Organization, apikey.SchemeSECP256K1)
			if err != nil {
				OutputError(eris.Wrap(err, "failed to create API keypair"))
			}
		case "ed25519":
			apiKey, err = apikey.New(Organization, apikey.SchemeED25519)
			if err != nil {
				OutputError(eris.Wrap(err, "failed to create API keypair"))
			}
		}

		if name == "-" {
			Output(map[string]string{
				"publicKey":  apiKey.TkPublicKey,
				"privateKey": apiKey.TkPrivateKey,
			})
		}

		apiKey.Metadata.Name = name
		apiKey.Metadata.PublicKey = apiKey.TkPublicKey
		apiKey.Metadata.Scheme = apiKey.Scheme
		apiKey.Metadata.Organizations = []string{Organization}

		if err = apiKeyStore.Store(name, apiKey); err != nil {
			OutputError(eris.Wrap(err, "failed to store new API keypair"))
		}

		localStore, ok := apiKeyStore.(*local.Store[*apikey.Key, apikey.Metadata])
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

// Represents the command to generate an encryption key
var encryptionKeyCmd = &cobra.Command{
	Use:   "encryption-key",
	Short: "Generate a Turnkey encryption key",
	Long:  `Generate a new encryption key that can be used for encrypting text sent from Turnkey secure enclaves.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Organization == "" {
			OutputError(eris.New("--organization must be specified"))
		}

		if User == "" {
			OutputError(eris.New("--user must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("encryption-key-name")
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read encryption key name"))
		}

		encryptionKey, err := encryptionkey.New(User, Organization)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to create encryption keypair"))
		}

		if name == "-" {
			Output(map[string]string{
				"publicKey":  encryptionKey.TkPublicKey,
				"privateKey": encryptionKey.TkPrivateKey,
			})
		}

		encryptionKey.Metadata.Name = name
		encryptionKey.Metadata.PublicKey = encryptionKey.TkPublicKey
		encryptionKey.Metadata.Organization = Organization
		encryptionKey.Metadata.User = User

		if err = encryptionKeyStore.Store(name, encryptionKey); err != nil {
			OutputError(eris.Wrap(err, "failed to store new encryption keypair"))
		}

		localStore, ok := encryptionKeyStore.(*local.Store[*encryptionkey.Key, encryptionkey.Metadata])
		if !ok {
			OutputError(eris.Wrap(err, "unhandled keystore type: expected *local.Store"))
		}

		Output(map[string]string{
			"publicKey":      encryptionKey.TkPublicKey,
			"publicKeyFile":  localStore.PublicKeyFile(name),
			"privateKeyFile": localStore.PrivateKeyFile(name),
		})
	},
}
