package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkhq/tkcli/cmd/tk/internal/display"
	"github.com/tkhq/tkcli/internal/apikey"
	"github.com/urfave/cli/v2"
)

const (
	TK_FOLDER_NAME = ".tk"
	DEFAULT_HOST   = "coordinator.tkhq.xyz"
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
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Name of the API key. Will be used to create <folder>/<name>.public and <folder>/<name>.private. If you do not want to write files, use --name -",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "folder",
						Usage:    "Folder in which to put the API key. Defaults to `~/.tk`.",
						Required: false,
					},
				},
				Action: func(cCtx *cli.Context) error {
					apiKeyName := cCtx.String("name")
					folder := cCtx.String("folder")

					apiKey, err := apikey.NewTkApiKey()
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Could not create new key pair", 1)
					}

					if apiKeyName == "-" {
						jsonBytes, err := json.MarshalIndent(map[string]interface{}{
							"publicKey":  apiKey.TkPublicKey,
							"privateKey": apiKey.TkPrivateKey,
						}, "", "    ")
						if err != nil {
							log.Fatalf("Unable to serialize output to JSON: %v", err)
						}
						fmt.Println(string(jsonBytes))
						return nil
					} else {
						tkDirPath := getTkDirPath(folder)
						err := os.MkdirAll(tkDirPath, os.ModePerm)
						if err != nil {
							log.Fatalln(err)
							return cli.Exit("Could not create .tk directory in your home directory", 1)
						}

						publicKeyFile := fmt.Sprintf("%s/%s.public", tkDirPath, apiKeyName)
						privateKeyFile := fmt.Sprintf("%s/%s.private", tkDirPath, apiKeyName)
						createFile(publicKeyFile, apiKey.TkPublicKey, 0755)
						createFile(privateKeyFile, apiKey.TkPrivateKey, 0700)

						jsonBytes, err := json.MarshalIndent(map[string]interface{}{
							"publicKeyFile":  publicKeyFile,
							"privateKeyFile": privateKeyFile,
						}, "", "    ")
						if err != nil {
							log.Fatalf("Unable to serialize output to JSON: %v", err)
						}
						fmt.Println(string(jsonBytes))
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
					&cli.StringFlag{
						Name:     "host",
						Usage:    "HTTP host _without_ protocol. For example: api.turnkey.io",
						Required: false,
						Value:    DEFAULT_HOST,
					},
					&cli.StringFlag{
						Name:     "method",
						Usage:    "HTTP Method. Should be \"GET\" or \"POST\"",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "path",
						Usage:    "Path, including the leading \"/\" and query string if any. For example: /api/v1/keys?curve=ed25519",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "body",
						Usage:    "HTTP body, only relevant for POST requests. For example: {\"message\": \"hello from TKHQ\"}",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "Private key to sign with. Provide a name to lookup the private key in your ~/.tk directory (e.g. \"my_api_key\" will use \"~/.tk/my_api_key.private\"), or a full path to a valid private key (e.g. \"/path/to/key.private\")",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					method := cCtx.String("method")
					host := cCtx.String("host")
					path := cCtx.String("path")
					body := cCtx.String("body")
					protocol := "https"

					if pattern := regexp.MustCompile(`^localhost:\d+$`); pattern.MatchString(host) {
						protocol = "http"
					}

					signaturePayload := apikey.SerializeRequest(method, host, path, body)

					key := cCtx.String("key")
					apiKey, err := getApiKey(key)
					if err != nil {
						log.Fatalf("Unable to retrieve API key: %v", err)
					}

					signature, err := apikey.Sign(signaturePayload, apiKey)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Failed to produce a valid signature", 1)
					}

					var response *http.Response
					if method == "GET" {
						response, err = get(apiKey, protocol, host, path, signature)
						if err != nil {
							log.Fatalln(err)
						}
					} else if method == "POST" {
						response, err = post(apiKey, protocol, host, path, body, signature)
						if err != nil {
							log.Fatalln(err)
						}
					} else {
						return cli.Exit("Invalid method", 1)
					}

					displayResponse, err := display.DisplayResponse(response)
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
					&cli.StringFlag{
						Name:     "host",
						Usage:    "HTTP host _without_ protocol. For example: api.turnkey.io",
						Required: false,
						Value:    DEFAULT_HOST,
					},
					&cli.StringFlag{
						Name:     "method",
						Usage:    "HTTP Method. Should be \"GET\" or \"POST\"",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "path",
						Usage:    "Path, including the leading \"/\" and query string if any. For example: /api/v1/keys?curve=ed25519",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "body",
						Usage:    "HTTP body, only relevant for POST requests. For example: {\"message\": \"hello from TKHQ\"}",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "Private key to sign with. Provide a name to lookup the private key in your ~/.tk directory (e.g. \"my_api_key\" will use \"~/.tk/my_api_key.private\"), or a full path to a valid private key (e.g. \"/path/to/key.private\")",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					method := cCtx.String("method")
					host := cCtx.String("host")
					path := cCtx.String("path")
					body := cCtx.String("body")

					signaturePayload := apikey.SerializeRequest(method, host, path, body)

					key := cCtx.String("key")
					apiKey, err := getApiKey(key)
					if err != nil {
						log.Fatalf("Unable to retrieve API key: %v", err)
					}

					signature, err := apikey.Sign(signaturePayload, apiKey)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Failed to produce a valid signature", 1)
					}

					jsonBytes, err := json.MarshalIndent(map[string]interface{}{
						"message":        fmt.Sprintf("%q", signaturePayload),
						"signature":      signature,
						"approvalHeader": approvalHeader(apiKey, signature),
						"curlCommand":    generateCurlCommand(apiKey, method, host, path, body, signature),
					}, "", "    ")
					if err != nil {
						log.Fatalf("Unable to serialize output to JSON: %v", err)
					}
					fmt.Println(string(jsonBytes))

					return nil
				},
			},
			{
				Name:    "sign",
				Aliases: []string{"s"},
				Usage:   "sign an arbitrary message",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "message",
						Usage:    "Message to sign",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "Private key to sign with. Provide a name to lookup the private key in your ~/.tk directory (e.g. \"my_api_key\" will use \"~/.tk/my_api_key.private\"), or a full path to a valid private key (e.g. \"/path/to/key.private\")",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					message := cCtx.String("message")

					key := cCtx.String("key")
					var keyPath string
					if !strings.Contains(key, "/") && !strings.Contains(key, ".") {
						keyPath = fmt.Sprintf("%s/%s.private", getTkDirPath(""), key)
					} else {
						// We have a full file path. Try loading it directly
						keyPath = key
					}
					bytes, err := getFileContent(keyPath)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Could load private key", 1)
					}

					apiKey, err := apikey.FromTkPrivateKey(string(bytes))
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Could recover API key from private key file content", 1)
					}
					signature, err := apikey.Sign(message, apiKey)
					if err != nil {
						log.Fatalln(err)
						return cli.Exit("Failed to produce a valid signature", 1)
					}

					jsonBytes, err := json.MarshalIndent(map[string]interface{}{
						"message":   fmt.Sprintf("%q", message),
						"signature": signature,
					}, "", "    ")
					if err != nil {
						log.Fatalf("Unable to serialize output to JSON: %v", err)
					}

					fmt.Println(string(jsonBytes))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getTkDirPath(folder string) string {
	if folder == "" {
		usr, _ := user.Current()
		folder = usr.HomeDir + "/" + TK_FOLDER_NAME
		return folder
	} else {
		// Check that the directory exists
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			return folder
		} else {
			log.Fatalf("Cannot put key files in %s: %v", folder, err)
			return ""
		}
	}
}

func createFile(path string, content string, mode fs.FileMode) error {
	return os.WriteFile(path, []byte(content), mode)
}

func getFileContent(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getApiKey(key string) (*apikey.ApiKey, error) {
	var keyPath string
	if !strings.Contains(key, "/") && !strings.Contains(key, ".") {
		keyPath = fmt.Sprintf("%s/%s.private", getTkDirPath(""), key)
	} else {
		// We have a full file path. Try loading it directly
		keyPath = key
	}

	bytes, err := getFileContent(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load private key:")
	}

	apiKey, err := apikey.FromTkPrivateKey(string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "could recover API key from private key file content:")
	}
	return apiKey, nil
}

func generateCurlCommand(apiKey *apikey.ApiKey, method, host, path, body, signature string) string {
	if method == "POST" {
		return fmt.Sprintf("curl -X POST -d'%s' -H'%s' -v 'https://%s%s'", body, approvalHeader(apiKey, signature), host, path)
	} else {
		return fmt.Sprintf("curl -H'%s' -v 'https://%s%s'", approvalHeader(apiKey, signature), host, path)
	}
}

func approvalHeader(apiKey *apikey.ApiKey, signature string) string {
	return fmt.Sprintf("X-Approved-By-%s: %s", apiKey.TkPublicKey, signature)
}

func get(key *apikey.ApiKey, protocol string, host string, path string, signature string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s%s", protocol, host, path)
	req, _ := http.NewRequest("GET", url, nil)

	headerName := fmt.Sprintf("X-Approved-By-%s", key.TkPublicKey)
	req.Header.Set(headerName, signature)

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func post(key *apikey.ApiKey, protocol string, host string, path string, body string, signature string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s%s", protocol, host, path)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error while creating HTTP POST request")
	}

	headerName := fmt.Sprintf("X-Approved-By-%s", key.TkPublicKey)
	req.Header.Set(headerName, signature)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while sending HTTP POST request")
	}
	return response, nil
}
