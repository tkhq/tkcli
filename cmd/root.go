package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/internal/clifs"
)

var (
	rootKeysDirectory string

	// KeyName is the name of the key with which we are operating.
	KeyName string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootKeysDirectory, "keys-folder", "d", clifs.DefaultKeysDir(), "directory in which to locate keys")
	rootCmd.PersistentFlags().StringVarP(&KeyName, "key-name", "k", "default", "name of key with which to operate")

}

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "turnkey interacts with the Turnkey API",
	Short: "turnkey is the Turnkey CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := clifs.SetKeysDirectory(rootKeysDirectory); err != nil {
			return errors.Wrap(err, "failed to obtain key storage location")
		}

		return nil
	},
}
