package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
	"github.com/tkhq/go-sdk/pkg/encryption_key"
)

var (
	// Filepath to write the export bundle to.
	exportBundlePath string

	// EncryptionKeypair is the loaded Encryption Keypair.
	EncryptionKeypair *encryption_key.Key
)

func init() {
	decryptCmd.Flags().StringVar(&exportBundlePath, "export-bundle-input", "", "filepath to read the export bundle from.")
	decryptCmd.Flags().StringVar(&plaintextPath, "plaintext-output", "", "optional filepath to write the plaintext from that will be decrypted.")

	rootCmd.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a ciphertext",
	Long:  `Decrypt a ciphertext from a bundle exported from a Turnkey secure enclave.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadEncryptionKeypair("")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if exportBundlePath == "" {
			OutputError(eris.New("--export-bundle-input must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// read from export bundle path
		exportBundle, err := readFile(exportBundlePath)
		if err != nil {
			OutputError(err)
		}

		// get encryption key
		tkPrivateKey := EncryptionKeypair.GetPrivateKey()
		kemPrivateKey, err := encryption_key.DecodeTurnkeyPrivateKey(tkPrivateKey)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to decode encryption private key"))
		}

		// set up enclave encrypt client
		signerPublic, err := hexToPublicKey(signerPublicKey)
		if err != nil {
			OutputError(err)
		}

		encryptClient, err := enclave_encrypt.NewEnclaveEncryptClientFromTargetKey(signerPublic, *kemPrivateKey)
		if err != nil {
			OutputError(err)
		}

		// decrypt ciphertext
		plaintextBytes, err := encryptClient.Decrypt([]byte(exportBundle), Organization)
		if err != nil {
			OutputError(err)
		}

		plaintext := string(plaintextBytes)

		// output the plaintext if no filepath is passed
		if plaintextPath == "" {
			Output(plaintext)
			return
		}

		err = writeFile(plaintext, plaintextPath)
		if err != nil {
			OutputError(err)
		}
	},
}

// LoadEncryptionKeypair require-loads the keypair referenced by the given name or as referenced form the global KeyName variable, if name is empty.
func LoadEncryptionKeypair(name string) {
	if name == "" {
		name = EncryptionKeyName
	}

	if encryptionKeyStore == nil {
		OutputError(eris.New("encryption keystore not loaded"))
	}

	encryptionKey, err := encryptionKeyStore.Load(name)
	if err != nil {
		OutputError(err)
	}

	if encryptionKey == nil {
		OutputError(eris.New("Encryption key not loaded"))
	}

	EncryptionKeypair = encryptionKey

	// If we haven't had the organization explicitly set try to load it from key metadata.
	if Organization == "" {
		Organization = encryptionKey.Organization
	}

	// If org is _still_ empty, the API key is not usable.
	if Organization == "" {
		OutputError(eris.New("failed to associate the encryption key with an organization; please manually specify the organization ID"))
	}
}
