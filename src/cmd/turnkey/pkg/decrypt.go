package pkg

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
)

// Filepath to write the export bundle to.
var exportBundlePath string

func init() {
	decryptCmd.Flags().StringVar(&exportBundlePath, "export-bundle-path", "/export_bundle.txt", "filepath to write the export bundle to.")
	decryptCmd.Flags().StringVar(&plaintextPath, "plaintext-path", "", "filepath to write the plaintext from that will be decrypted.")

	rootCmd.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a ciphertext",
	Long:  `Decrypt a ciphertext from a bundle exported from a Turnkey secure enclave.`,
	Run: func(cmd *cobra.Command, args []string) {
		// read from export bundle path
		exportBundle, err := readFile(exportBundlePath)
		if err != nil {
			OutputError(err)
		}

		var serverSendMsg enclave_encrypt.ServerSendMsg
		err = json.Unmarshal([]byte(exportBundle), &serverSendMsg)
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

		// decrypt ciphertext
		plaintextBytes, err := encryptClient.Decrypt(serverSendMsg)
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
