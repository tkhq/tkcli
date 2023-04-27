package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/models"
)

func init() {
	curvesCmd.AddCommand(curvesListCmd)

	rootCmd.AddCommand(curvesCmd)
}

var curvesCmd = &cobra.Command{
	Use:   "curves",
	Short: "curves interacts with the available curves",
}

var curvesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list returns the available curve types",
	Run: func(cmd *cobra.Command, args []string) {
		Output(models.Curves())
	},
}
