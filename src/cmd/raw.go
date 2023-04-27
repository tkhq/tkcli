package cmd

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/private_keys"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/util"
)

var (
	rawSignPayloadEncoding string
	rawSignHashFunction    string
	rawSignPayload         string
)

func init() {
	rawCmd.Flags().StringVar(&signingKeyID, "private-key", "", "name or ID of the private signing key")

	rawSignCmd.Flags().StringVar(&rawSignPayloadEncoding, "payload-encoding",
		string(models.V1PayloadEncodingPAYLOADENCODINGTEXTUTF8), "encoding of payload")
	rawSignCmd.Flags().StringVar(&rawSignHashFunction, "hash-function",
		string(models.V1HashFunctionHASHFUNCTIONSHA256), "hash function")
	rawSignCmd.Flags().StringVar(&rawSignPayload, "payload",
		"", "payload to be signed")

	rawCmd.AddCommand(rawSignCmd)

	rootCmd.AddCommand(rawCmd)
}

var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "raw allows low-level (raw) requests to the Turnkey API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)

		LoadKeypair("")

		LoadSigningKey("")

		LoadClient()
	},
}

var rawSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign signs a raw payload",
	Run: func(cmd *cobra.Command, args []string) {
		payloadEncoding := models.V1PayloadEncoding(rawSignPayloadEncoding)

		hashFunction := models.V1HashFunction(rawSignHashFunction)

		payload, err := ParameterToString(rawSignPayload)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read payload"))
		}

		activityType := string(models.V1ActivityTypeACTIVITYTYPESIGNRAWPAYLOAD)

		params := private_keys.NewPublicAPIServiceSignRawPayloadParams().WithBody(
			&models.V1SignRawPayloadRequest{
				OrganizationID: &privateKeysOrgID,
				Parameters: &models.V1SignRawPayloadIntent{
					Encoding:     &payloadEncoding,
					HashFunction: &hashFunction,
					PrivateKeyID: &signingKeyID,
					Payload:      &payload,
				},
				TimestampMs: util.RequestTimestamp(),
				Type:        &activityType,
			},
		)

		resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceSignRawPayload(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
