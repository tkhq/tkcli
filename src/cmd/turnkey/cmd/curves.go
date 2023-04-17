package cmd

import (
	"github.com/tkhq/tkcli/src/api/models"

	"github.com/spf13/cobra"
)

func init() {
	curvesCmd.AddCommand(curvesListCmd)

	rootCmd.AddCommand(curvesCmd)
}

var curvesCmd = &cobra.Command{
	Use:   "curves interacts with the available curves",
	Short: "curves interacts with the available curves",
}

var curvesListCmd = &cobra.Command{
	Use:   "list returns the available curve types",
	Short: "list returns the available curve types",
	Run: func(cmd *cobra.Command, args []string) {
		Output(models.Curves())
	},
}
