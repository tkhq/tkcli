package cmd

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/api/client"
	"github.com/tkhq/tkcli/api/client/private_keys"
	"github.com/tkhq/tkcli/api/models"
)

var (
	signingKeyID string

	signRawPayloadEncoding string
	signRawHashFunction    string

	signTransactionType string
)

func init() {
	signCmd.Flags().StringVar(&signingKeyID, "signing-key", "", "name or ID of the signing key")

	signRawCmd.Flags().StringVar(&signRawPayloadEncoding, "payload-encoding",
		string(models.V1PayloadEncodingPAYLOADENCODINGTEXTUTF8), "encoding of payload")
	signRawCmd.Flags().StringVar(&signRawHashFunction, "hash-function",
		string(models.V1HashFunctionHASHFUNCTIONSHA256), "hash function")

	signTransactionCmd.Flags().StringVar(&signTransactionType, "type", string(models.Immutableactivityv1TransactionTypeTRANSACTIONTYPEETHEREUM), "type of transaction; for a list of valid transaction types, issue `turnkey transaction-types list`")

	rootCmd.AddCommand(signCmd)
}

var signCmd = &cobra.Command{
	Use:     "sign operations available with Turnkey",
	Short:   "sign accesses signing operations provided by Turnkey",
	Aliases: []string{"s"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := LoadKeypair(""); err != nil {
			OutputError(errors.Wrap(err, "failed to load API key"))
		}

		if _, err := uuid.Parse(signingKeyID); err != nil {
			signingKeyID, err = lookupPrivateKeyByName(signingKeyID)
			if err != nil {
				OutputError(errors.Wrap(err, "provided private key was not a UUID and lookup by name failed"))
			}
		}
	},
}

var signRawCmd = &cobra.Command{
	Use:     "raw signs a raw payload",
	Short:   "raw signs a raw payload",
	Aliases: []string{"tx"},
	Run: func(cmd *cobra.Command, args []string) {
		payloadEncoding := models.V1PayloadEncoding(signRawPayloadEncoding)

		hashFunction := models.V1HashFunction(signRawHashFunction)
		
		payload := "TODO"

		activityType := string(models.V1ActivityTypeACTIVITYTYPESIGNRAWPAYLOAD)

		params := private_keys.NewPublicAPIServiceSignRawPayloadParams()
		params.SetBody(&models.V1SignRawPayloadRequest{
			OrganizationID: &privateKeysOrgID,
			Parameters: &models.V1SignRawPayloadIntent{
				Encoding:     &payloadEncoding,
				HashFunction: &hashFunction,
				PrivateKeyID: &signingKeyID,
				Payload:      &payload,
			},
			TimestampMs: RequestTimestamp(),
			Type:        &activityType,
		})

		resp, err := client.Default.PrivateKeys.PublicAPIServiceSignRawPayload(params, new(Authenticator))
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to create private key: %d: %s", resp.Code(), resp.Error()))
		}

		Output(resp.Payload)
	},
}

var signTransactionCmd = &cobra.Command{
	Use:     "transaction signs a transaction",
	Short:   "transaction provides signing for a transaction",
	Aliases: []string{"tx"},
	Run: func(cmd *cobra.Command, args []string) {
		transactionType := models.Immutableactivityv1TransactionType(signTransactionType)
		activityType := string(models.V1ActivityTypeACTIVITYTYPESIGNTRANSACTION)

		payload := ""

		params := private_keys.NewPublicAPIServiceSignTransactionParams()
		params.SetBody(&models.V1SignTransactionRequest{
			OrganizationID: &privateKeysOrgID,
			Parameters: &models.V1SignTransactionIntent{
				PrivateKeyID:        &signingKeyID,
				Type:                &transactionType,
				UnsignedTransaction: &payload,
			},
			TimestampMs: RequestTimestamp(),
			Type:        &activityType,
		})

		resp, err := client.Default.PrivateKeys.PublicAPIServiceSignTransaction(params, new(Authenticator))
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to create private key: %d: %s", resp.Code(), resp.Error()))
		}

		Output(resp.Payload)
	},
}

func lookupPrivateKeyByName(name string) (string, error) {
	params := private_keys.NewPublicAPIServiceGetPrivateKeysParams()

	params.SetBody(&models.V1GetPrivateKeysRequest{
		OrganizationID: &Organization,
	})

	if err := params.Body.Validate(nil); err != nil {
		return "", errors.Wrap(err, "formulation of a lookup by name request failed")
	}

	resp, err := client.Default.PrivateKeys.PublicAPIServiceGetPrivateKeys(params, new(Authenticator))
	if err != nil {
		return "", errors.Wrap(err, "lookup by name failed")
	}

	if !resp.IsSuccess() {
		return "", errors.Errorf("lookup by name failed: %d: %s", resp.Code(), resp.Error())
	}

	for _, k := range resp.Payload.PrivateKeys {
		if *k.PrivateKeyName == name {
			return *k.PrivateKeyID, nil
		}
	}

	return "", errors.Errorf("private key name %q not found in list of private keys", name)
}
