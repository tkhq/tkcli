package display

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Formatter int64

const (
	FormatterUnknown Formatter = iota
	FormatterPretty
	FormatterJSON
)

func ParseFormatter(cCtx *cli.Context) (Formatter, error) {
	input := cCtx.String("formatter")

	switch input {
	case "json":
		return FormatterJSON, nil
	case "pretty":
		return FormatterPretty, nil
	default:
		return FormatterUnknown, errors.Errorf("Unknown formatter %s", input)
	}
}

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
	case FormatterPretty:
		bytes, err := yaml.MarshalWithOptions(input, yaml.Indent(4))
		if err != nil {
			return "", errors.Wrap(err, "unable to serialize output to YAML")
		}
		return string(bytes), nil
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
	case FormatterPretty:
		output, err := JSONToYAML(input)
		if err != nil { // Not a JSON
			return string(input), nil
		}

		return string(output), nil
	default:
		return "", errors.Errorf("Unknown formatter %v", formatter)
	}
}

// https://github.com/goccy/go-yaml/blob/894a764b31ce8c62a845a1e626cd43c6bb475a7a/yaml.go#L240 with custom formatting options
func JSONToYAML(bytes []byte) ([]byte, error) {
	var v interface{}
	if err := yaml.UnmarshalWithOptions(bytes, &v, yaml.UseOrderedMap()); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal from json bytes")
	}
	out, err := yaml.MarshalWithOptions(v, yaml.Indent(4))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal")
	}
	return out, nil
}
