package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/api/client"
	"github.com/tkhq/tkcli/api/client/private_keys"
	"github.com/tkhq/tkcli/api/models"
)

var (
)

func init() {
	ethCmd.Flags().StringVar(&signingKeyID, "signing-key", "", "name or ID of the signing key")

	rootCmd.AddCommand(ethCmd)
}

var ethCmd = &cobra.Command{
	Use:     "ethereum performs actions related to Ethereum",
	Short:   "ethereum performs actions related to Ethereum",
	Aliases: []string{"eth"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		LoadKeypair("")
		LoadSigningKey("")
	},
}

var ethTxCmd = &cobra.Command{
	Use:     "transaction signs a transaction",
	Short:   "transaction provides signing for a transaction",
	Aliases: []string{"tx"},
	Run: func(cmd *cobra.Command, args []string) {
		transactionType := models.Immutableactivityv1TransactionTypeTRANSACTIONTYPEETHEREUM
		activityType := string(models.V1ActivityTypeACTIVITYTYPESIGNTRANSACTION)

		// TODO: define the payload!
		payload := ""

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
