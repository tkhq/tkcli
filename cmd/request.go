package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/internal/apikey"
	"github.com/tkhq/tkcli/internal/clifs"
	"github.com/tkhq/tkcli/internal/display"
)

var (
	requestHost, requestPath, requestBody string
	requestNoPost, requestShowCurlCommand bool
)

func init() {
	makeRequest.Flags().StringVar(&requestHost, "host", "coordinator-beta.turnkey.io", "hostname of the API server")
	makeRequest.Flags().StringVar(&requestPath, "path", "", "path for the API request")
	makeRequest.Flags().StringVar(&requestPath, "body", "", "body of the request, which can be '-' to indicate stdin or be prefixed with '@' to indicate a source filename")
	makeRequest.Flags().BoolVar(&requestNoPost, "no-post", false, "only provide the signature, do not post the request to the API server")
	makeRequest.Flags().BoolVar(&requestShowCurlCommand, "show-curl", false, "only provide the signature, do not post the request to the API server")
}

var makeRequest = &cobra.Command{
	Use:     "request takes a request body, generates a stamp for the given request and optionally sends it to the Turnkey API server",
	Short:   "request makes a basic API request",
	Aliases: []string{"req", "r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		protocol := "https"
		if pattern := regexp.MustCompile(`^localhost:\d+$`); pattern.MatchString(requestHost) {
			protocol = "http"
		}

		keyName, err := cmd.PersistentFlags().GetString("key")
		if err != nil {
			return errors.Wrap(err, "failed to read key name parameter")
		}

		apiKey, err := clifs.LoadKeypair(keyName)
		if err != nil {
			return errors.Wrap(err, "failed to get API key")
		}

		bodyReader, err := processRequestBody(requestBody)
		if err != nil {
			return errors.Wrap(err, "failed to process request body")
		}

		body, err := io.ReadAll(bodyReader)
		if err != nil {
			return errors.Wrap(err, "failed to read message body")
		}

		stamp, err := apikey.Signature(body, apiKey)
		if err != nil {
			return errors.Wrap(err, "failed to produce a valid API stamp")
		}

		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "   ")

		if requestNoPost {
			ret := map[string][]byte{
				"message": body,
				"stamp":   stamp,
			}

			if requestShowCurlCommand {
				ret["curlCommand"] = []byte(generateCurlCommand(requestHost, requestPath, body, stamp))
			}

			return enc.Encode(ret)
		}

		response, err := post(protocol, requestHost, requestPath, body, stamp)
		if err != nil {
			return errors.Wrap(err, "failed to post request")
		}

		defer response.Body.Close()

		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return &display.ErrorResponse{
				Code: response.StatusCode,
				Text: response.Status,
			}
		}

		if response.StatusCode != http.StatusOK {
			return &display.ErrorResponse{
				Code: response.StatusCode,
				Text: string(responseBodyBytes),
			}
		}

		return enc.Encode(responseBodyBytes)
	},
}

func processRequestBody(bodyParam string) (io.Reader, error) {
	if bodyParam == "-" {
		return os.Stdin, nil
	}

	if strings.HasPrefix(bodyParam, "@") {
		return os.Open(strings.TrimPrefix(bodyParam, "@"))
	}

	buf := new(bytes.Buffer)

	if _, err := buf.WriteString(bodyParam); err != nil {
		return nil, errors.Wrap(err, "failed to read from body parameter")
	}

	return buf, nil
}

func generateCurlCommand(host, path string, body, stamp []byte) string {
	return fmt.Sprintf("curl -X POST -d'%s' -H'%s' -v 'https://%s%s'", body, stampHeader(stamp), host, path)
}

func stampHeader(stamp []byte) string {
	return fmt.Sprintf("X-Stamp: %s", string(stamp))
}

func post(protocol string, host string, path string, body []byte, stamp []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s%s", protocol, host, path)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error while creating HTTP POST request")
	}

	req.Header.Set("X-Stamp", string(stamp))

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while sending HTTP POST request")
	}

	return response, nil
}
