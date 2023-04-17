package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var outputFormat string

type encoder interface {
	Encode(data any) error
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format")
}

func getEncoder() encoder {
	switch outputFormat {
	case "yaml":
		enc := yaml.NewEncoder(os.Stdout)
		enc.SetIndent(2)

		return enc
	default: // JSON is the default
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "   ")

		return enc
	}
}

// OutputError prints an error to the console and exits.
func OutputError(err error) {
	if err := getEncoder().Encode(map[string]string{
		"error": err.Error(),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write error to output encoder: %s", err)
	}

	os.Exit(1)
}

// Output prints to the console and exits.
func Output(payload any) {
	if err := getEncoder().Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode output: %s", err)
	}

	os.Exit(0)
}
