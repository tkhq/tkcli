package pkg

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
	"github.com/tkhq/go-sdk/pkg/encryptionkey"
)

var (
	// Filepath to read the export bundle from.
	exportBundlePath string

	// EncryptionKeypair is the loaded Encryption Keypair.
	EncryptionKeypair *encryptionkey.Key
)

func init() {
	decryptCmd.Flags().StringVar(&exportBundlePath, "export-bundle-input", "", "filepath to read the export bundle from.")
	decryptCmd.Flags().StringVar(&plaintextPath, "plaintext-output", "", "optional filepath to write the plaintext from that will be decrypted.")
	decryptCmd.Flags().StringVar(&signerPublicKeyOverride, "signer-quorum-key", "", "optional override for the signer quorum key. This option should be used for testing only. Leave this value empty for production decryptions.")

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
		kemPrivateKey, err := encryptionkey.DecodeTurnkeyPrivateKey(tkPrivateKey)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to decode encryption private key"))
		}

		var signerKey *ecdsa.PublicKey
		if signerPublicKeyOverride != "" {
			signerKey, err = hexToPublicKey(signerPublicKeyOverride)
		} else {
			signerKey, err = hexToPublicKey(signerProductionPublicKey)
		}
		if err != nil {
			OutputError(err)
		}

		// set up enclave encrypt client
		encryptClient, err := enclave_encrypt.NewEnclaveEncryptClientFromTargetKey(signerKey, *kemPrivateKey)
		if err != nil {
			OutputError(err)
		}

		// decrypt ciphertext
		plaintextBytes, err := encryptClient.Decrypt([]byte(exportBundle), Organization)
		if err != nil {
			OutputError(err)
		}

		plaintext := hex.EncodeToString(plaintextBytes)

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

// LoadEncryptionKeypair require-loads the keypair referenced by the given name or as referenced from the global EncryptionKeyName variable, if name is empty.
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

	// If org is _still_ empty, the encryption key is not usable.
	if Organization == "" {
		OutputError(eris.New("failed to associate the encryption key with an organization; please manually specify the organization ID"))
	}

	// If we haven't had the user explicitly set try to load it from key metadata.
	if User == "" {
		User = encryptionKey.User
	}

	// If user is _still_ empty, the encryption key is still usable in some cases where user ID isn't needed (export)
	// Hence we do not error out here if encryptionKey.User is empty.
}
