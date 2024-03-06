package pkg

import (
	"encoding/hex"
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/tkhq/go-sdk/pkg/enclave_encrypt"
)

// Filepath to write the export bundle to.
var ExportBundlePath string

func init() {
	decryptCmd.Flags().StringVar(&ExportBundlePath, "export-bundle-path", "/export_bundle.txt", "filepath to write the export bundle to.")

	rootCmd.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a ciphertext",
	Long:  `Decrypt a ciphertext from a bundle exported from a Turnkey secure enclave.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if PlaintextPath == "" {
			OutputError(eris.New("Filepath for plaintext must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// read from export bundle path
		exportBundle, err := readFile(ExportBundlePath)
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
		plaintext := hex.EncodeToString(plaintextBytes)

		// output the hex-encoded plaintext if no filepath is passed
		if PlaintextPath == "" {
			Output(plaintext)
			return
		}

		err = writeFile(plaintext, PlaintextPath)
		if err != nil {
			OutputError(err)
		}
	},
}
