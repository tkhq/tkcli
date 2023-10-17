package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionString string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display build and version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionString)
	},
}
