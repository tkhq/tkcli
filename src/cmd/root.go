package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/store"
	"github.com/tkhq/go-sdk/pkg/store/local"
)

var (
	rootKeysDirectory string

	keyStore store.Store

	// KeyName is the name of the key with which we are operating.
	KeyName string

	// Organization is the organization ID to interact with.
	Organization string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootKeysDirectory, "keys-folder", "d", local.DefaultKeysDir(), "directory in which to locate keys")
	rootCmd.PersistentFlags().StringVarP(&KeyName, "key-name", "k", "default", "name of API key with which to interact with the Turnkey API service")

	rootCmd.PersistentFlags().StringVar(&Organization, "organization", "", "organization ID to be used")
}

func basicSetup(cmd *cobra.Command) {
	// No non-JSON-formatted output should flow over stdin; thus change
	// output for usage messages to stderr.
	cmd.SetOut(os.Stderr)

	if keyStore == nil {
		localKeyStore := local.New()

		if err := localKeyStore.SetKeysDirectory(rootKeysDirectory); err != nil {
			OutputError(errors.Wrap(err, "failed to obtain key storage location"))
		}

		keyStore = localKeyStore
	}
}

// Execute runs the cobra command for the Turnkey CLI.
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "turnkey interacts with the Turnkey API",
	Short: "turnkey is the Turnkey CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
	},
}

// RequestTimestamp returns a timestamp formatted for inclusion in a request.
func RequestTimestamp() *string {
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())

	return &ts
}

// ParameterToReader converts a commandline parameter to an io.Reader based on its syntax.
// Values of "-" return stdin.
// Values beginning with "@" return the file with name following the "@".
// Other values are interpreted literally.
func ParameterToReader(param string) (io.Reader, error) {
	if param == "-" {
		return os.Stdin, nil
	}

	if strings.HasPrefix(param, "@") {
		return os.Open(strings.TrimPrefix(param, "@"))
	}

	return bytes.NewReader([]byte(param)), nil
}

// ParameterToString processes a commandline parameter with ParameterToReader, reads it in, and then returns a string with its contents.
// See ParameterToReader for conversion details.
func ParameterToString(param string) (string, error) {
	payloadReader, err := ParameterToReader(param)
	if err != nil {
		return "", errors.Wrap(err, "failed to process payload")
	}

	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(payloadReader); err != nil {
		return "", errors.Wrap(err, "failed to read payload data")
	}

	return buf.String(), nil
}
