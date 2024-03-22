package pkg

import (
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcutil/base58"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
)

var (
	// user is the user ID to import wallets and private keys with.
	User string

	// Filepath to write the import bundle to.
	importBundlePath string

	// Filepath to read the encrypted bundle from.
	encryptedBundlePath string

	// Filepath to read the plaintext from that will be encrypted.
	plaintextPath string

	// Format to apply to the plaintext key before it's encrypted: `mnemonic`, `hexadecimal`, `solana`. Defaults to `mnemonic`.
	keyFormat string
)

func init() {
	encryptCmd.Flags().StringVar(&importBundlePath, "import-bundle-input", "", "filepath to write the import bundle to.")
	encryptCmd.Flags().StringVar(&encryptedBundlePath, "encrypted-bundle-output", "", "filepath to read the encrypted bundle from.")
	encryptCmd.Flags().StringVar(&plaintextPath, "plaintext-input", "", "filepath to read the plaintext from that will be encrypted.")
	encryptCmd.Flags().StringVar(&keyFormat, "key-format", "mnemonic", "optional formatting to apply to the plaintext before it is encrypted.")
	encryptCmd.Flags().StringVar(&user, "user", "", "ID of user to encrypting the plaintext.")

	rootCmd.AddCommand(encryptCmd)
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a plaintext",
	Long:  `Encrypt a plaintext into a bundle to be imported to a Turnkey secure enclave.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if importBundlePath == "" {
			OutputError(eris.New("--import-bundle-input must be specified"))
		}

		if encryptedBundlePath == "" {
			OutputError(eris.New("--encrypted-bundle-output must be specified"))
		}

		if plaintextPath == "" {
			OutputError(eris.New("--plaintext-input must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// read from import bundle path
		importBundle, err := readFile(importBundlePath)
		if err != nil {
			OutputError(err)
		}

		// set up enclave encrypt client
		signerPublic, err := hexToPublicKey(signerPublicKey)
		if err != nil {
			OutputError(err)
		}

		encryptClient, err := enclave_encrypt.NewEnclaveEncryptClient(signerPublic)
		if err != nil {
			OutputError(err)
		}

		// format the plaintext key
		plaintext, err := readFile(plaintextPath)
		if err != nil {
			OutputError(err)
		}
		var plaintextBytes []byte
		switch keyFormat {
		case "mnemonic":
			plaintextBytes = []byte(plaintext)
		case "hexadecimal":
			plaintextBytes, err = hex.DecodeString(plaintext)
			if err != nil {
				OutputError(err)
			}
		case "solana":
			decoded := base58.Decode(plaintext)
			if len(decoded) < 32 {
				OutputError(eris.New("invalid plaintext length. must be at least 32 bytes for key-format `solana`"))
			}
			plaintextBytes = decoded[:32]
		default:
			OutputError(eris.New("--key-format is invalid. accepted values: mnemonic, hexadecimal, solana."))
		}

		// encrypt plaintext
		clientSendMsg, err := encryptClient.Encrypt(plaintextBytes, []byte(importBundle), Organization, user)
		if err != nil {
			OutputError(err)
		}

		encryptedBundleBytes, err := json.Marshal(clientSendMsg)
		if err != nil {
			OutputError(err)
		}

		// write to encrypted bundle path
		err = writeFile(string(encryptedBundleBytes), encryptedBundlePath)
		if err != nil {
			OutputError(err)
		}
	},
}
