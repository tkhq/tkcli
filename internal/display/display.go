package display

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Formatter int64

const (
	FormatterJSON Formatter = iota
)

func FormatResponse(response *http.Response, formatter Formatter) (string, error) {
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "cannot read response body")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		displayMessage, err := FormatStruct(map[string]interface{}{
			"responseCode": response.StatusCode,
			"responseBody": string(responseBytes),
		}, formatter)
		if err != nil {
			return "", err
		}

		return displayMessage, nil
	}

	displayMessage, err := formatJSONBytes(responseBytes, formatter)
	if err != nil {
		return "", err
	}

	return displayMessage, nil
}

func FormatStruct(input map[string]interface{}, formatter Formatter) (string, error) {
	switch formatter {
	case FormatterJSON:
		jsonBytes, err := json.MarshalIndent(input, "", "    ")
		if err != nil {
			return "", errors.Wrap(err, "unable to serialize output to JSON")
		}

		return string(jsonBytes), nil
	default:
		return "", errors.Errorf("Unknown formatter %v", formatter)
	}
}

func formatJSONBytes(input []byte, formatter Formatter) (string, error) {
	switch formatter {
	case FormatterJSON:
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, input, "", "    ")
		if err != nil { // Not a JSON
			return string(input), nil
		}

		return prettyJSON.String(), nil
	default:
		return "", errors.Errorf("Unknown formatter %v", formatter)
	}
}
