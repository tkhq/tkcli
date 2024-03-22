package pkg

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/apikey"
	"github.com/tkhq/go-sdk/pkg/encryption_key"
	"github.com/tkhq/go-sdk/pkg/store"
	"github.com/tkhq/go-sdk/pkg/store/local"
)

var (
	rootKeysDirectory string

	encryptionKeysDirectory string

	apiKeyStore store.Store[apikey.Key]

	encryptionKeyStore store.Store[encryption_key.Key]

	// KeyName is the name of the key with which we are operating.
	KeyName string

	apiHost string

	// Organization is the organization ID to interact with.
	Organization string
)

// Turnkey Signer enclave's public key.
const signerPublicKey = "04ca7c0d624c75de6f34af342e87a21e0d8c83efd1bd5b5da0c0177c147f744fba6f01f9f37356f9c617659aafa55f6e0af8d169a8f054d153ab3201901fb63ecb04cf288fe433cc4e1aa0ce1632feac4ea26bf2f5a09dcfe5a42c398e06898710330f0572882f4dbdf0f5304b8fc8703acd69adca9a4bbf7f5d00d20a5e364b2569"

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootKeysDirectory, "keys-folder", "d", local.DefaultKeysDir(), "directory in which to locate keys")
	rootCmd.PersistentFlags().StringVarP(&encryptionKeysDirectory, "encryption-keys-folder", "d", local.DefaultKeysDir(), "directory in which to locate keys")
	rootCmd.PersistentFlags().StringVarP(&KeyName, "key-name", "k", "default", "name of API key with which to interact with the Turnkey API service")
	rootCmd.PersistentFlags().StringVar(&apiHost, "host", "api.turnkey.com", "hostname of the API server")

	rootCmd.PersistentFlags().StringVar(&Organization, "organization", "", "organization ID to be used")
}

func basicSetup(cmd *cobra.Command) {
	// No non-JSON-formatted output should flow over stdin; thus change
	// output for usage messages to stderr.
	cmd.SetOut(os.Stderr)

	err := detectAndMoveDeprecatedDefaultKeysDirOnMacOs()
	if err != nil {
		OutputError(err)
	}

	if apiKeyStore == nil {
		localKeyStore := local.New[apikey.Key]()

		if err := localKeyStore.SetKeysDirectory(rootKeysDirectory); err != nil {
			OutputError(eris.Wrap(err, "failed to obtain key storage location"))
		}

		apiKeyStore = localKeyStore
	}

	if encryptionKeyStore == nil {
		localEncryptionKeyStore := local.New[encryption_key.Key]()

		if err := localEncryptionKeyStore.SetKeysDirectory(encryptionKeysDirectory); err != nil {
			OutputError(eris.Wrap(err, "failed to obtain key storage location"))
		}

		encryptionKeyStore = localEncryptionKeyStore
	}
}

// Execute runs the cobra command for the Turnkey CLI.
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "turnkey",
	Short: "turnkey is the Turnkey CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
	},
}

// ParameterToReader converts a commandline parameter to an io.Reader based on its syntax.
// Values of "-" return stdin.
// Values beginning with "@" return the file with name following the "@".
// Other values are interpreted literally.
func ParameterToReader(param string) (io.Reader, error) {
	if param == "-" {
		return os.Stdin, nil
	}

	if strings.HasPrefix(param, "@") {
		return os.Open(strings.TrimPrefix(param, "@"))
	}

	return bytes.NewReader([]byte(param)), nil
}

// ParameterToString processes a commandline parameter with ParameterToReader, reads it in, and then returns a string with its contents.
// See ParameterToReader for conversion details.
func ParameterToString(param string) (string, error) {
	payloadReader, err := ParameterToReader(param)
	if err != nil {
		return "", eris.Wrap(err, "failed to process payload")
	}

	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(payloadReader); err != nil {
		return "", eris.Wrap(err, "failed to read payload data")
	}

	return buf.String(), nil
}

func detectAndMoveDeprecatedDefaultKeysDirOnMacOs() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	deprecatedDir := local.DeprecatedDefaultKeysDir()
	if deprecatedDir == "" {
		return nil
	}

	newDir := local.DefaultKeysDir()
	fmt.Printf("Legacy keys directory detected; will migrate keys to new location\n- Legacy: %s\n- New: %s\n\n", deprecatedDir, newDir)

	err := filepath.WalkDir(deprecatedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		relativeFilePath, err := filepath.Rel(deprecatedDir, path)
		if err != nil {
			return err
		}

		destFilePath := filepath.Join(newDir, relativeFilePath)

		err = safeRename(path, destFilePath)

		if err != nil {
			return err
		}

		fmt.Printf("Moved `%s` to %s\n", relativeFilePath, destFilePath)

		return nil
	})
	if err != nil {
		return err
	}

	err = os.RemoveAll(deprecatedDir)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("Successfully migrated legacy keys directory.")
	fmt.Println("")

	return nil
}

// Like `os.Rename(...)`, but does not allow overwriting
func safeRename(oldPath string, newPath string) error {
	exists, err := checkExists(newPath)
	if err != nil {
		return err
	}

	if exists {
		return eris.Errorf("target path already exists: %s", newPath)
	}

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	return nil
}

// Reads the content from the given file path and trims whitespace.
func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", eris.Wrap(err, "error reading file")
	}

	return strings.TrimSpace(string(content)), nil
}

// Writes the given content to a file at the specified path.
func writeFile(content string, path string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return eris.Wrap(err, "error writing file")
	}
	return nil
}

func checkExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Convert a hex-encoded string to an ECDSA P-256 public key.
// This key is used in encryption and decryption of data transferred to
// and from Turnkey secure enclaves.
func hexToPublicKey(hexString string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	if len(publicKeyBytes) != 130 {
		return nil, eris.New("invalid public key length")
	}
	encryptionPublicKeyBytes := publicKeyBytes[65:130]

	// init curve instance
	curve := elliptic.P256()

	// curve's bitsize converted to length in bytes
	byteLen := (curve.Params().BitSize + 7) / 8

	// ensure the public key bytes have the correct length
	if len(encryptionPublicKeyBytes) != 1+2*byteLen {
		return nil, eris.New("invalid encryption public key length")
	}

	// extract X and Y coordinates from the public key bytes
	// ignore first byte (prefix)
	x := new(big.Int).SetBytes(encryptionPublicKeyBytes[1 : 1+byteLen])
	y := new(big.Int).SetBytes(encryptionPublicKeyBytes[1+byteLen:])

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}
