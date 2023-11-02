package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/signers"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/util"
)

var (
	ethTxSigner  string
	ethTxPayload string
)

func init() {
	rootCmd.AddCommand(ethCmd)

	ethTxCmd.Flags().StringVarP(&ethTxSigner, "signer", "s", "", "wallet account address, private key address, or private key ID")
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
		LoadClient()
	},
}

var ethTxCmd = &cobra.Command{
	Use:     "transaction",
	Short:   "Perform signing and other actions for a transaction",
	Aliases: []string{"tx"},
	Run: func(cmd *cobra.Command, args []string) {
		transactionType := models.TransactionTypeEthereum
		activityType := string(models.ActivityTypeSignTransaction)

		payload, err := ParameterToString(ethTxPayload)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read payload"))
		}

		// NB: eventually, we should add ways of creating transaction payloads, to be more helpful.
		// Until then, this is an error.
		if payload == "" {
			OutputError(eris.New("payload cannot be empty"))
		}

		params := signers.NewSignTransactionParams().WithBody(
			&models.SignTransactionRequest{
				OrganizationID: &Organization,
				Parameters: &models.SignTransactionIntentV2{
					SignWith:            &ethTxSigner,
					Type:                &transactionType,
					UnsignedTransaction: &payload,
				},
				TimestampMs: util.RequestTimestamp(),
				Type:        &activityType,
			},
		)

		resp, err := APIClient.V0().Signers.SignTransaction(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
