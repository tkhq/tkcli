package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/internal/apikey"
	"github.com/tkhq/tkcli/internal/clifs"
	"github.com/tkhq/tkcli/internal/display"
)

var (
	requestHost, requestPath, requestBody string
	requestNoPost                         bool
)

func init() {
	makeRequest.Flags().StringVar(&requestHost, "host", "coordinator-beta.turnkey.io", "hostname of the API server")
	makeRequest.Flags().StringVar(&requestPath, "path", "", "path for the API request")
	makeRequest.Flags().StringVar(&requestBody, "body", "-", "body of the request, which can be '-' to indicate stdin or be prefixed with '@' to indicate a source filename")
	makeRequest.Flags().BoolVar(&requestNoPost, "no-post", false, "generates the stamp and displays the cURL command to use in order to perform this action, but does NOT post the request to the API server")

	rootCmd.AddCommand(makeRequest)
}

var makeRequest = &cobra.Command{
	Use:     "request takes a request body, generates a stamp for the given request, and sends it to the Turnkey API server.  See options for alternate behaviour, such as not sending the request.",
	Short:   "request makes a basic API request",
	Aliases: []string{"req", "r"},
	Run: func(cmd *cobra.Command, args []string) {
		protocol := "https"
		if pattern := regexp.MustCompile(`^localhost:\d+$`); pattern.MatchString(requestHost) {
			protocol = "http"
		}

		apiKey, err := clifs.LoadKeypair(KeyName)
		if err != nil {
			OutputError(errors.Wrap(err, "failed to get API key"))
		}

		bodyReader, err := ParameterToReader(requestBody)
		if err != nil {
			OutputError(errors.Wrap(err, "failed to process request body"))
		}

		body, err := io.ReadAll(bodyReader)
		if err != nil {
			OutputError(errors.Wrap(err, "failed to read message body"))
		}

		stamp, err := apikey.Stamp(body, apiKey)
		if err != nil {
			OutputError(errors.Wrap(err, "failed to produce a valid API stamp"))
		}

		if requestNoPost {
			Output(map[string]string{
				"message":     string(body),
				"stamp":       stamp,
				"curlCommand": generateCurlCommand(requestHost, requestPath, body, stamp),
			})
		}

		response, err := post(protocol, requestHost, requestPath, body, stamp)
		if err != nil {
			OutputError(errors.Wrap(err, "failed to post request"))
		}

		defer response.Body.Close()

		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			OutputError(&display.ErrorResponse{
				Code: response.StatusCode,
				Text: response.Status,
			})
		}

		if response.StatusCode != http.StatusOK {
			OutputError(&display.ErrorResponse{
				Code: response.StatusCode,
				Text: string(responseBodyBytes),
			})
		}

		Output(responseBodyBytes)
	},
}

func generateCurlCommand(host, path string, body []byte, stamp string) string {
	return fmt.Sprintf("curl -X POST -d'%s' -H'%s' -v 'https://%s%s'", body, stampHeader(stamp), host, path)
}

func stampHeader(stamp string) string {
	return fmt.Sprintf("X-Stamp: %s", stamp)
}

func post(protocol string, host string, path string, body []byte, stamp string) (*http.Response, error) {
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
