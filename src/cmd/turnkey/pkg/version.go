package pkg

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tkhq/tkcli/src/internal/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display build and version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s", version.Version)
		fmt.Printf("Commit:  %s", version.Commit)
	},
}
