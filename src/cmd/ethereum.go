package cmd

import (
	"github.com/tkhq/go-sdk/pkg/api/client/private_keys"
	"github.com/tkhq/go-sdk/pkg/api/models"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ethTxPayload string

func init() {
	ethCmd.PersistentFlags().StringVarP(&signingKeyID, "signing-key", "s", "", "name or ID of the signing key")

	rootCmd.AddCommand(ethCmd)

	ethTxCmd.Flags().StringVar(&ethTxPayload, "payload", "", "payload of the transaction")

	ethCmd.AddCommand(ethTxCmd)
}

var ethCmd = &cobra.Command{
	Use:     "ethereum performs actions related to Ethereum",
	Short:   "ethereum performs actions related to Ethereum",
	Aliases: []string{"eth"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		LoadKeypair("")
		LoadSigningKey("")
		LoadClient()
	},
}

var ethTxCmd = &cobra.Command{
	Use:     "transaction provides signing and other actions for a transaction",
	Short:   "transaction provides signing and other actions for a transaction",
	Aliases: []string{"tx"},
	Run: func(cmd *cobra.Command, args []string) {
		transactionType := models.Immutableactivityv1TransactionTypeTRANSACTIONTYPEETHEREUM
		activityType := string(models.V1ActivityTypeACTIVITYTYPESIGNTRANSACTION)

		payload, err := ParameterToString(ethTxPayload)
		if err != nil {
			OutputError(errors.Wrap(err, "failed to read payload"))
		}

		// NB: eventually, we should add ways of creating transaction payloads, to be more helpful.
		// Until then, this is an error.
		if payload == "" {
			OutputError(errors.New("payload cannot be empty"))
		}

		params := private_keys.NewPublicAPIServiceSignTransactionParams().WithBody(
			&models.V1SignTransactionRequest{
				OrganizationID: &privateKeysOrgID,
				Parameters: &models.V1SignTransactionIntent{
					PrivateKeyID:        &signingKeyID,
					Type:                &transactionType,
					UnsignedTransaction: &payload,
				},
				TimestampMs: RequestTimestamp(),
				Type:        &activityType,
			},
		)

		resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceSignTransaction(params, APIClient.Authenticator)
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
