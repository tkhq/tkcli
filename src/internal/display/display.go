package display

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func DisplayResponse(response *http.Response) (string, error) {
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "cannot read response body")
	}
	if response.StatusCode == 200 {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, responseBytes, "", "    ")
		if err != nil { // Not a JSON
			return string(responseBytes), nil
		}

		return prettyJSON.String(), nil
	}

	jsonBytes, err := json.MarshalIndent(map[string]interface{}{
		"responseCode": response.StatusCode,
		"responseBody": string(responseBytes),
	}, "", "    ")
	if err != nil {
		return "", errors.Wrap(err, "unable to serialize output to JSON")
	}
	return string(jsonBytes), nil
}
