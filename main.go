package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkhq/tkcli/internal/apikey"
	"github.com/tkhq/tkcli/internal/clifs"
	"github.com/tkhq/tkcli/internal/display"
	"github.com/tkhq/tkcli/internal/flags"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tk",
		Usage: "The Turnkey CLI",
		Commands: []*cli.Command{
			{
				Name:    "generate-api-key",
				Aliases: []string{"gen"},
				Usage:   "generate a new Turnkey API key",
				Flags: []cli.Flag{
					flags.CreateKeyName(),
					flags.KeysFolder(),
				},
				Action: func(cCtx *cli.Context) error {
					apiKeyName := cCtx.String("name")
					folder := cCtx.String("keys-folder")
					formatter := display.FormatterJSON

					apiKey, err := apikey.NewTkApiKey()
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Could not create new key pair", 1)
					}

					if apiKeyName == "-" {
						displayMessage, err := display.FormatStruct(map[string]interface{}{
							"publicKey":  apiKey.TkPublicKey,
							"privateKey": apiKey.TkPrivateKey,
						}, formatter)
						if err != nil {
							log.Fatalf("Unable to format output: %v", err)
						}
						fmt.Println(displayMessage)
						return nil
					} else {
						tkDirPath, err := clifs.GetKeyDirPath(folder)
						if err != nil {
							log.Fatalln(err)
							return cli.Exit("Could not create determine key directory location", 1)

						}

						err = os.MkdirAll(tkDirPath, os.ModePerm)
						if err != nil {
							log.Fatalln(err)
							return cli.Exit(fmt.Sprintf("Could not create directory %s", tkDirPath), 1)
						}

						publicKeyFile := fmt.Sprintf("%s/%s.public", tkDirPath, apiKeyName)
						privateKeyFile := fmt.Sprintf("%s/%s.private", tkDirPath, apiKeyName)
						clifs.CreateFile(publicKeyFile, apiKey.TkPublicKey, 0755)
						clifs.CreateFile(privateKeyFile, apiKey.TkPrivateKey, 0700)

						displayMessage, err := display.FormatStruct(map[string]interface{}{
							"publicKey":      apiKey.TkPublicKey,
							"publicKeyFile":  publicKeyFile,
							"privateKeyFile": privateKeyFile,
						}, formatter)
						if err != nil {
							log.Fatalf("Unable to format output: %v", err)
						}
						fmt.Println(displayMessage)
					}

					return nil
				},
			},
			{
				Name:      "request",
				Aliases:   []string{"r"},
				Usage:     "make a request",
				UsageText: "generate an approval and make an HTTP(s) request",
				Flags: []cli.Flag{
					flags.Host(),
					flags.Path(),
					flags.Body(),
					flags.Key(),
				},
				Action: func(cCtx *cli.Context) error {
					host := cCtx.String("host")
					path := cCtx.String("path")
					body := cCtx.String("body")
					formatter := display.FormatterJSON
					protocol := "https"

					if pattern := regexp.MustCompile(`^localhost:\d+$`); pattern.MatchString(host) {
						protocol = "http"
					}

					key := cCtx.String("key")
					apiKey, err := clifs.GetApiKey(key)
					if err != nil {
						log.Fatalf("Unable to retrieve API key: %v", err)
					}

					stamp, err := apikey.Stamp(body, apiKey)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Failed to produce a valid API stamp", 1)
					}

					response, err := post(protocol, host, path, body, stamp)
					if err != nil {
						log.Fatalln(err)
					}

					displayResponse, err := display.FormatResponse(response, formatter)
					if err != nil {
						log.Fatalf("unable to display response: %v", err)
					}

					fmt.Println(displayResponse)
					return nil
				},
			},
			{
				Name:      "approve-request",
				Aliases:   []string{"approve"},
				Usage:     "approve a request",
				UsageText: "generate an approval over an HTTP request",
				Flags: []cli.Flag{
					flags.Key(),
					flags.Host(),
					flags.Path(),
					flags.Body(),
				},
				Action: func(cCtx *cli.Context) error {
					host := cCtx.String("host")
					path := cCtx.String("path")
					body := cCtx.String("body")
					formatter := display.FormatterJSON

					key := cCtx.String("key")
					apiKey, err := clifs.GetApiKey(key)
					if err != nil {
						log.Fatalf("Unable to retrieve API key: %v", err)
					}

					stamp, err := apikey.Stamp(body, apiKey)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Failed to produce a valid stamp", 1)
					}

					displayMessage, err := display.FormatStruct(map[string]interface{}{
						"message":     body,
						"stamp":       stamp,
						"curlCommand": generateCurlCommand(host, path, body, stamp),
					}, formatter)
					if err != nil {
						log.Fatalf("Unable to format output: %v", err)
					}
					fmt.Println(displayMessage)

					return nil
				},
			},
			{
				Name:    "stamp",
				Aliases: []string{"s"},
				Usage:   "sign an arbitrary message and produce a valid API Stamp",
				Flags: []cli.Flag{
					flags.Message(),
					flags.Key(),
				},
				Action: func(cCtx *cli.Context) error {
					message := cCtx.String("message")
					formatter := display.FormatterJSON

					key := cCtx.String("key")

					var keyPath string
					if !strings.Contains(key, "/") && !strings.Contains(key, ".") {
						keysDirectory, err := clifs.GetKeyDirPath("")
						if err != nil {
							log.Fatalln(err)
							return cli.Exit("Could not load keys directory path", 1)
						}
						keyPath = fmt.Sprintf("%s/%s.private", keysDirectory, key)
					} else {
						// We have a full file path. Try loading it directly
						keyPath = key
					}
					bytes, err := clifs.GetFileContent(keyPath)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Could load private key", 1)
					}

					apiKey, err := apikey.FromTkPrivateKey(string(bytes))
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Could recover API key from private key file content", 1)
					}
					stamp, err := apikey.Stamp(message, apiKey)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Failed to produce a valid stamp", 1)
					}

					displayMessage, err := display.FormatStruct(map[string]interface{}{
						"message": fmt.Sprintf("%q", message),
						"stamp":   stamp,
					}, formatter)
					if err != nil {
						log.Fatalf("Unable to format output: %v", err)
					}

					fmt.Println(displayMessage)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generateCurlCommand(host, path, body, stamp string) string {
	return fmt.Sprintf("curl -X POST -d'%s' -H'%s' -v 'https://%s%s'", body, stampHeader(stamp), host, path)
}

func stampHeader(stamp string) string {
	return fmt.Sprintf("X-Stamp: %s", stamp)
}

func post(protocol string, host string, path string, body string, stamp string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s%s", protocol, host, path)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error while creating HTTP POST request")
	}

	req.Header.Set("X-Stamp", stamp)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while sending HTTP POST request")
	}
	return response, nil
}
