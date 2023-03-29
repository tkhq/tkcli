package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tkhq/tkcli/cmd/turnkey/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "   ")

		if encErr := enc.Encode(err); encErr != nil {
			fmt.Print(err)
		}

		os.Exit(1)
	}
}
