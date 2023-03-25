package main

import (
	"encoding/json"
	"os"

	"github.com/tkhq/tkcli/cmd"
)

func main() {
   if err := cmd.Execute(); err != nil {
      enc := json.NewEncoder(os.Stderr)
      enc.SetIndent("", "   ")

      enc.Encode(err)
   }
}
