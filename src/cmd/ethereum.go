package cmd

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/private_keys"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/util"
)

var ethTxPayload string

func init() {
	ethCmd.PersistentFlags().StringVarP(&signingKeyID, "private-key", "s", "", "name or ID of the private signing key")

	rootCmd.AddCommand(ethCmd)

	ethTxCmd.Flags().StringVar(&ethTxPayload, "payload", "", "payload of the transaction")

	ethCmd.AddCommand(ethTxCmd)
}

var ethCmd = &cobra.Command{
	Use:     "ethereum",
	Short:   "Perform actions related to Ethereum",
	Aliases: []string{"eth"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadSigningKey("")
		LoadClient()
	},
}

var ethTxCmd = &cobra.Command{
	Use:     "transaction",
	Short:   "Perform signing and other actions for a transaction",
	Aliases: []string{"tx"},
	Run: func(cmd *cobra.Command, args []string) {
		transactionType := models.Immutableactivityv1TransactionTypeTRANSACTIONTYPEETHEREUM
		activityType := string(models.V1ActivityTypeACTIVITYTYPESIGNTRANSACTION)

		payload, err := ParameterToString(ethTxPayload)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read payload"))
		}

		// NB: eventually, we should add ways of creating transaction payloads, to be more helpful.
		// Until then, this is an error.
		if payload == "" {
			OutputError(eris.New("payload cannot be empty"))
		}

		params := private_keys.NewPublicAPIServiceSignTransactionParams().WithBody(
			&models.V1SignTransactionRequest{
				OrganizationID: &Organization,
				Parameters: &models.V1SignTransactionIntent{
					PrivateKeyID:        &signingKeyID,
					Type:                &transactionType,
					UnsignedTransaction: &payload,
				},
				TimestampMs: util.RequestTimestamp(),
				Type:        &activityType,
			},
		)

		resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceSignTransaction(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
