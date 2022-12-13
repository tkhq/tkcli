// Package to encapsulate CLI filesystem operations
package clifs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkhq/tkcli/internal/apikey"
)

// Given a user-specified directory, return the path.
// The logic is in the case where users do not specify a folder.
// If the folder isn't specified, we default to $XDG_CONFIG_HOME/.config/turnkey/keys.
// If this env var isn't set, we default to $HOME/.config/turnkey/keys
// If $HOME isn't set, this function returns an error.
func GetKeyDirPath(userSpecifiedPath string) (string, error) {
	if userSpecifiedPath == "" {
		var configHome string
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			configHome = os.Getenv("XDG_CONFIG_HOME")
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", errors.Wrap(err, "error while reading user home directory")
			}
			configHome = homeDir + "/.config"
		}
		return configHome + "/turnkey/keys", nil
	} else {
		if _, err := os.Stat(userSpecifiedPath); !os.IsNotExist(err) {
			return userSpecifiedPath, nil
		} else {
			return "", errors.Errorf("Cannot put key files in %s: %v", userSpecifiedPath, err)
		}
	}
}

func CreateFile(path string, content string, mode fs.FileMode) error {
	return os.WriteFile(path, []byte(content), mode)
}

func GetFileContent(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GetApiKey(key string) (*apikey.ApiKey, error) {
	var keyPath string
	if !strings.Contains(key, "/") && !strings.Contains(key, ".") {
		keysDirectory, err := GetKeyDirPath("")
		if err != nil {
			return nil, errors.Wrap(err, "unable to get keys directory path")
		}
		keyPath = fmt.Sprintf("%s/%s.private", keysDirectory, key)
	} else {
		// We have a full file path. Try loading it directly
		keyPath = key
	}

	bytes, err := GetFileContent(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load private key:")
	}

	apiKey, err := apikey.FromTkPrivateKey(string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "could recover API key from private key file content:")
	}
	return apiKey, nil
}
