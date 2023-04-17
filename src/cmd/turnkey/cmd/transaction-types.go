package cmd

import (
	"github.com/tkhq/tkcli/src/api/models"

	"github.com/spf13/cobra"
)

func init() {
	transactionTypesCmd.AddCommand(transactionTypesListCmd)

	rootCmd.AddCommand(transactionTypesCmd)
}

var transactionTypesCmd = &cobra.Command{
	Use:   "transaction-types interacts with the available transaction types",
	Short: "transaction-types interacts with the available transaction types",
}

var transactionTypesListCmd = &cobra.Command{
	Use:   "list returns the available transaction types",
	Short: "list returns the available transaction types",
	Run: func(cmd *cobra.Command, args []string) {
		Output(models.TransactionTypes())
	},
}
