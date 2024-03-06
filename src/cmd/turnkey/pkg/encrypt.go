package pkg

import (
	"encoding/hex"
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
)

var (
	// User is the user ID to import wallets and private keys with.
	User string

	// Filepath to write the import bundle to.
	ImportBundlePath string

	// Filepath to read the encrypted bundle from.
	EncryptedBundlePath string

	// Filepath to read the plaintext from that will be encrypted.
	PlaintextPath string
)

func init() {
	encryptCmd.Flags().StringVar(&ImportBundlePath, "import-bundle-path", "/import_bundle.txt", "filepath to write the import bundle to.")
	encryptCmd.Flags().StringVar(&EncryptedBundlePath, "encrypted-bundle-path", "/encrypted_bundle.txt", "filepath to read the encrypted bundle from.")
	encryptCmd.Flags().StringVar(&PlaintextPath, "plaintext-path", "", "filepath to read the plaintext from that will be encrypted.")

	rootCmd.AddCommand(encryptCmd)
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a plaintext",
	Long:  `Encrypt a plaintext into a bundle to be imported to a Turnkey secure enclave.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if PlaintextPath == "" {
			OutputError(eris.New("Filepath for plaintext must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// read from import bundle path
		importBundle, err := readFile(ImportBundlePath)
		if err != nil {
			OutputError(err)
		}

		var serverTargetMsg enclave_encrypt.ServerTargetMsg
		err = json.Unmarshal([]byte(importBundle), &serverTargetMsg)
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

		// encrypt plaintext
		plaintext, err := readFile(PlaintextPath)
		if err != nil {
			OutputError(err)
		}

		var plaintextBytes []byte
		plaintextBytes, err = hex.DecodeString(plaintext)
		if err != nil {
			plaintextBytes = []byte(plaintext)
		}
		clientSendMsg, err := encryptClient.Encrypt(plaintextBytes, serverTargetMsg)
		if err != nil {
			OutputError(err)
		}

		encryptedBundleBytes, err := json.Marshal(clientSendMsg)
		if err != nil {
			OutputError(err)
		}

		// write to encrypted bundle path
		err = writeFile(string(encryptedBundleBytes), EncryptedBundlePath)
		if err != nil {
			OutputError(err)
		}
	},
}
