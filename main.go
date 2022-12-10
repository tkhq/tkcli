package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkhq/tkcli/internal/apikey"
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
						tkDirPath, err := getKeyDirPath(folder)
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
					flags.Host(),
					flags.Method(),
					flags.Path(),
					flags.Body(),
					flags.Key(),
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
					flags.Host(),
					flags.Method(),
					flags.Path(),
					flags.Body(),
					flags.KeysFolder(),
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
					flags.Message(),
					flags.Key(),
				},
				Action: func(cCtx *cli.Context) error {
					message := cCtx.String("message")

					key := cCtx.String("key")

					var keyPath string
					if !strings.Contains(key, "/") && !strings.Contains(key, ".") {
						keysDirectory, err := getKeyDirPath("")
						if err != nil {
							log.Fatalln(err)
							return cli.Exit("Could not load keys directory path", 1)
						}
						keyPath = fmt.Sprintf("%s/%s.private", keysDirectory, key)
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

// Given a user-specified directory, return the path.
// The logic is in the case where users do not specify a folder.
// If the folder isn't specified, we default to $XDG_CONFIG_HOME/.config/turnkey/keys.
// If this env var isn't set, we default to $HOME/.config/turnkey/keys
// If $HOME isn't set, this function returns an error.
func getKeyDirPath(userSpecifiedPath string) (string, error) {
	if userSpecifiedPath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", nil
		}
		return userConfigDir + "turnkey/keys", nil
	} else {
		if _, err := os.Stat(userSpecifiedPath); !os.IsNotExist(err) {
			return userSpecifiedPath, nil
		} else {
			return "", errors.Errorf("Cannot put key files in %s: %v", userSpecifiedPath, err)
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
		keysDirectory, err := getKeyDirPath("")
		if err != nil {
			return nil, errors.Wrap(err, "unable to get keys directory path")
		}
		keyPath = fmt.Sprintf("%s/%s.private", keysDirectory, key)
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
