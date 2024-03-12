package pkg

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
)

var (
	// user is the user ID to import wallets and private keys with.
	user string

	// Filepath to write the import bundle to.
	importBundlePath string

	// Filepath to read the encrypted bundle from.
	encryptedBundlePath string

	// Filepath to read the plaintext from that will be encrypted.
	plaintextPath string
)

func init() {
	encryptCmd.Flags().StringVar(&importBundlePath, "import-bundle-input", "", "filepath to write the import bundle to.")
	encryptCmd.Flags().StringVar(&encryptedBundlePath, "encrypted-bundle-output", "", "filepath to read the encrypted bundle from.")
	encryptCmd.Flags().StringVar(&plaintextPath, "plaintext-input", "", "filepath to read the plaintext from that will be encrypted.")

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
		plaintext, err := readFile(plaintextPath)
		if err != nil {
			OutputError(err)
		}

		clientSendMsg, err := encryptClient.Encrypt([]byte(plaintext), serverTargetMsg)
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
