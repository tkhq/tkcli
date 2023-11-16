package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/signers"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/util"
)

var (
	rawSigner              string
	rawSignPayloadEncoding string
	rawSignHashFunction    string
	rawSignPayload         string
)

func init() {
	rawSignCmd.Flags().StringVarP(&rawSigner, "signer", "s", "", "wallet account address, private key address, or private key ID")
	rawSignCmd.Flags().StringVar(&rawSignPayloadEncoding, "payload-encoding",
		string(models.PayloadEncodingTextUTF8), "encoding of payload")
	rawSignCmd.Flags().StringVar(&rawSignHashFunction, "hash-function",
		string(models.HashFunctionSha256), "hash function")
	rawSignCmd.Flags().StringVar(&rawSignPayload, "payload",
		"", "payload to be signed")

	rawCmd.AddCommand(rawSignCmd)

	rootCmd.AddCommand(rawCmd)
}

var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Send low-level (raw) requests to the Turnkey API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
	},
}

var rawSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a raw payload",
	Run: func(cmd *cobra.Command, args []string) {
		payloadEncoding := models.PayloadEncoding(rawSignPayloadEncoding)
		hashFunction := models.HashFunction(rawSignHashFunction)

		payload, err := ParameterToString(rawSignPayload)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read payload"))
		}

		activityType := string(models.ActivityTypeSignRawPayloadV2)

		params := signers.NewSignRawPayloadParams().WithBody(
			&models.SignRawPayloadRequest{
				OrganizationID: &Organization,
				Parameters: &models.SignRawPayloadIntentV2{
					SignWith:     &rawSigner,
					Encoding:     &payloadEncoding,
					HashFunction: &hashFunction,
					Payload:      &payload,
				},
				TimestampMs: util.RequestTimestamp(),
				Type:        &activityType,
			},
		)

		resp, err := APIClient.V0().Signers.SignRawPayload(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to sign raw payload: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
