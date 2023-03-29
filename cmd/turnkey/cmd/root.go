package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/internal/clifs"
)

var (
	rootKeysDirectory string

	// KeyName is the name of the key with which we are operating.
	KeyName string

	// Organization is the organization ID to interact with.
	Organization string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootKeysDirectory, "keys-folder", "d", clifs.DefaultKeysDir(), "directory in which to locate keys")
	rootCmd.PersistentFlags().StringVarP(&KeyName, "key-name", "k", "default", "name of key with which to operate")

	rootCmd.PersistentFlags().StringVar(&Organization, "organization", "", "organization ID to be used")
}

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "turnkey interacts with the Turnkey API",
	Short: "turnkey is the Turnkey CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// No non-JSON-formatted output should flow over stdin; thus change
		// output for usage messages to stderr.
		cmd.SetOut(os.Stderr)

		if err := clifs.SetKeysDirectory(rootKeysDirectory); err != nil {
			OutputError(errors.Wrap(err, "failed to obtain key storage location"))
		}
	},
}

// RequestTimestamp returns a timestamp formatted for inclusion in a request.
func RequestTimestamp() *string {
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())
	return &ts
}
