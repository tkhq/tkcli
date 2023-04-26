package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/models"
)

func init() {
	transactionTypesCmd.AddCommand(transactionTypesListCmd)

	rootCmd.AddCommand(transactionTypesCmd)
}

var transactionTypesCmd = &cobra.Command{
	Use:   "transaction-types",
	Short: "transaction-types interacts with the available transaction types",
}

var transactionTypesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list returns the available transaction types",
	Run: func(cmd *cobra.Command, args []string) {
		Output(models.TransactionTypes())
	},
}
