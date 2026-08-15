package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/apikey"
)

var (
	requestPath, requestBody string
	requestNoPost            bool
)

func init() {
	makeRequest.Flags().StringVar(&requestBody, "body", "-", "body of the request, which can be '-' to indicate stdin or be prefixed with '@' to indicate a source filename")
	makeRequest.Flags().BoolVar(&requestNoPost, "no-post", false, `generates the stamp and displays
		the cURL command to use in order to perform this action,
		but does NOT post the request to the API server`)
	makeRequest.Flags().StringVar(&requestPath, "path", "", "path for the API request")

	rootCmd.AddCommand(makeRequest)
}

var makeRequest = &cobra.Command{
	Use: "request",
	Short: `Given a request body, generate a stamp for it, and send it to the Turnkey API server.
			See options for alternate behavior, such as previewing but not sending the request.`,
	Aliases: []string{"req", "r"},
	PreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
	},
	Run: func(cmd *cobra.Command, args []string) {
		protocol := "https"
		if pattern := regexp.MustCompile(`^localhost:\d+$`); pattern.MatchString(apiHost) {
			protocol = "http"
		}

		bodyReader, err := ParameterToReader(requestBody)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to process request body"))
		}

		body, err := io.ReadAll(bodyReader)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to read message body"))
		}

		stamp, err := apikey.Stamp(body, APIKeypair)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to produce a valid API stamp"))
		}

		if requestNoPost {
			Output(map[string]string{
				"message":     string(body),
				"stamp":       stamp,
				"curlCommand": generateCurlCommand(apiHost, requestPath, body, stamp),
			})
		}

		response, err := post(cmd.Context(), protocol, apiHost, requestPath, body, stamp)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to post request"))
		}

		defer response.Body.Close() //nolint: errcheck

		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			OutputError(&ResponseError{
				Code: response.StatusCode,
				Text: response.Status,
			})
		}

		if response.StatusCode != http.StatusOK {
			OutputError(&ResponseError{
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

func post(ctx context.Context, protocol string, host string, path string, body []byte, stamp string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s%s", protocol, host, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, eris.Wrap(err, "error while creating HTTP POST request")
	}

	req.Header.Set("X-Stamp", stamp)

	client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return eris.New("stopped after 10 redirects")
		}

		origin := via[0].URL
		if !strings.EqualFold(req.URL.Scheme, origin.Scheme) || !strings.EqualFold(req.URL.Host, origin.Host) {
			return eris.Errorf("refusing to follow redirect from %s://%s to %s", origin.Scheme, origin.Host, req.URL.Redacted())
		}

		return nil
	}}

	response, err := client.Do(req)
	if err != nil {
		return nil, eris.Wrap(err, "error while sending HTTP POST request")
	}

	return response, nil
}
