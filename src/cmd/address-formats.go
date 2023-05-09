package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/models"
)

func init() {
	addressFormatsCmd.AddCommand(addressFormatsListCmd)

	rootCmd.AddCommand(addressFormatsCmd)
}

var addressFormatsCmd = &cobra.Command{
	Use:   "address-formats",
	Short: "Interact with the available address formats",
}

var addressFormatsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return the available key formats",
	Run: func(cmd *cobra.Command, args []string) {
		Output(models.AddressFormats())
	},
}
