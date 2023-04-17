// Package to encapsulate CLI filesystem operations
package clifs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/tkhq/tkcli/src/internal/apikey"

	"github.com/pkg/errors"
)

var keysDirectory string

const (
	// DefaultKeyName is the name of the default API key
	DefaultKeyName = "default"

	turnkeyDirectoryName = "turnkey"
	keysDirectoryName    = "keys"
	publicKeyExtension   = "public"
	privateKeyExtension  = "private"
	metadataExtension    = "meta"
)

func init() {
	keysDirectory = DefaultKeysDir()
}

func createKeyFile(path string, content string, mode fs.FileMode) error {
	return os.WriteFile(path, []byte(content), mode)
}

func createMetadataFile(path string, key *apikey.ApiKey, mode fs.FileMode) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create metadata file")
	}

	return json.NewEncoder(f).Encode(key)
}

// checkFileExists checks that the given file exists and has a non-zero size.
func checkFileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if stat.Size() < 1 {
		return false, fmt.Errorf("file %q is empty", path)
	}

	return true, nil
}

// PublicKeyFile returns the filename for the public key of the given name.
func PublicKeyFile(name string) string {
	if name == "" {
		name = DefaultKeyName
	}

	return path.Join(keysDirectory, fmt.Sprintf("%s.%s", name, publicKeyExtension))
}

// PrivateKeyFile returns the filename for the private key of the given name.
func PrivateKeyFile(name string) string {
	if name == "" {
		name = DefaultKeyName
	}

	return path.Join(keysDirectory, fmt.Sprintf("%s.%s", name, privateKeyExtension))
}

// MetadataFile returns the filename for the metadata of the given key name.
func MetadataFile(name string) string {
	if name == "" {
		name = DefaultKeyName
	}

	return path.Join(keysDirectory, fmt.Sprintf("%s.%s", name, metadataExtension))
}

// DefaultKeysDir returns the default directory for key storage for the user's system.
func DefaultKeysDir() string {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			cfgDir = "."
		} else {
			cfgDir = path.Join(homeDir, ".config")
		}
	}

	return path.Join(cfgDir, turnkeyDirectoryName, keysDirectoryName)
}

// SetKeysDirectory sets the clifs root directory, ensuring its existence and writability.
func SetKeysDirectory(keysPath string) (err error) {
	if keysPath == "" || keysPath == DefaultKeysDir() {
		keysPath = DefaultKeysDir()

		// NB: we only attempt to create the default directory; never a user-supplied one.
		if err = os.MkdirAll(keysPath, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create key store location %q", keysPath)
		}
	}

	stat, err := os.Stat(keysPath)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.Errorf("keys directory %q is not a directory", keysPath)
	}

	keysDirectory = keysPath

	return nil
}

// StoreKeypair saves the given keypair to the key directory with the given name.
func StoreKeypair(name string, keypair *apikey.ApiKey) error {
	pubExists, err := checkFileExists(PublicKeyFile(name))
	if err != nil {
		return errors.Wrap(err, "failed to check for existence of public key")
	}

	privExists, err := checkFileExists(PrivateKeyFile(name))
	if err != nil {
		return errors.Wrap(err, "failed to check for existence of private key")
	}

	if pubExists || privExists {
		return errors.Errorf("a keypair named %q already exists! Exiting...", name)
	}

	if err = createKeyFile(PublicKeyFile(name), keypair.TkPublicKey, 0o0644); err != nil {
		return errors.Wrap(err, "failed to store public key to file")
	}

	if err = createKeyFile(PrivateKeyFile(name), keypair.TkPrivateKey, 0o0600); err != nil {
		return errors.Wrap(err, "failed to store private key to file")
	}

	if err = createMetadataFile(MetadataFile(name), keypair, 0o0600); err != nil {
		return errors.Wrap(err, "failed to store api key metadata")
	}

	return nil
}

// LoadKeypair reads a keypair from the keys directory.
func LoadKeypair(keyname string) (*apikey.ApiKey, error) {
	keyPath := PrivateKeyFile(keyname)

	// If we are given an explicit path, try to use it directly, rather than as the key name.
	if strings.Contains(keyname, "/") {
		keyPath = strings.TrimSuffix(keyname, "."+privateKeyExtension)

		exists, _ := checkFileExists(keyPath)
		if !exists {
			keyPath = keyPath + "." + privateKeyExtension

			exists, _ = checkFileExists(keyPath)
			if !exists {
				return nil, errors.Errorf("failed to load key %q", keyname)
			}
		}
	}

	bytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read from %q", keyPath)
	}

	apiKey, err := apikey.FromTkPrivateKey(string(bytes))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to recover API key from private key file %q", keyPath)
	}

	if ok, _ := checkFileExists(MetadataFile(keyname)); ok {
		metadata, err := loadMetadata(MetadataFile(keyname))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load key metadata from metadata file %q", MetadataFile(keyname))
		}

		if err := apiKey.MergeMetadata(metadata); err != nil {
			return nil, errors.Wrap(err, "failed to merge key metadata with key")
		}
	}

	return apiKey, nil
}

func loadMetadata(fn string) (*apikey.Metadata, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open metadata file")
	}

	md := new(apikey.Metadata)

	if err := json.NewDecoder(f).Decode(md); err != nil {
		return nil, errors.Wrap(err, "failed to decode metadata file")
	}

	return md, nil
}
