package main_test

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tkhq/go-sdk/pkg/apikey"
)

var TurnkeyBinaryName = "turnkey.linux-x86_64"

// TempDir is the directory in which temporary files for the tests will be stored.
var TempDir = "/tmp"

func init() {
	if os.Getenv("RUNNER_TEMP") != "" {
		TempDir = os.Getenv("RUNNER_TEMP")
	}
	var arch string
	switch runtime.GOARCH {
	case "arm64":
		arch = "aarch64"
	case "amd64":
		arch = "x86_64"
	}
	TurnkeyBinaryName = fmt.Sprintf("turnkey.%s-%s", runtime.GOOS, arch)
}

func RunCliWithArgs(t *testing.T, args []string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(path.Join(currentDir, "../../../out/", TurnkeyBinaryName), args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func TestHelpText(t *testing.T) {
	out, err := RunCliWithArgs(t, []string{})
	assert.Nil(t, err)
	assert.Contains(t, out, "the Turnkey CLI")
	assert.Contains(t, out, "Usage:")
	assert.Contains(t, out, "Available Commands:")
}

func TestAPIKeygenInTmpFolder(t *testing.T) {
	orgID := uuid.New()

	tmpDir, err := os.MkdirTemp(TempDir, "keys")
	assert.Nil(t, err)

	defer func() { assert.Nil(t, os.RemoveAll(tmpDir)) }()

	out, err := RunCliWithArgs(t, []string{"generate", "api-key", "--keys-folder", tmpDir, "--key-name", "mykey", "--organization", orgID.String()})
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/mykey.public")
	assert.FileExists(t, tmpDir+"/mykey.private")

	publicKeyData, err := os.ReadFile(tmpDir + "/mykey.public")
	assert.Nil(t, err)

	var parsedOut map[string]string

	assert.Nil(t, json.Unmarshal([]byte(out), &parsedOut))

	assert.Equal(t, parsedOut["publicKey"], string(publicKeyData))
	assert.Equal(t, parsedOut["publicKeyFile"], tmpDir+"/mykey.public")
	assert.Equal(t, parsedOut["privateKeyFile"], tmpDir+"/mykey.private")
}

func TestEncryptionKeygenInTmpFolder(t *testing.T) {
	orgID := uuid.New()
	userID := uuid.New()

	tmpDir, err := os.MkdirTemp(TempDir, "encryption-keys")
	assert.Nil(t, err)

	defer func() { assert.Nil(t, os.RemoveAll(tmpDir)) }()

	out, err := RunCliWithArgs(t, []string{"generate", "encryption-key", "--encryption-keys-folder", tmpDir, "--encryption-key-name", "mykey", "--organization", orgID.String(), "--user", userID.String()})
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/mykey.public")
	assert.FileExists(t, tmpDir+"/mykey.private")

	publicKeyData, err := os.ReadFile(tmpDir + "/mykey.public")
	assert.Nil(t, err)

	var parsedOut map[string]string

	assert.Nil(t, json.Unmarshal([]byte(out), &parsedOut))

	assert.Equal(t, parsedOut["publicKey"], string(publicKeyData))
	assert.Equal(t, parsedOut["publicKeyFile"], tmpDir+"/mykey.public")
	assert.Equal(t, parsedOut["privateKeyFile"], tmpDir+"/mykey.private")
}

func TestAPIKeygenDetectExistingKey(t *testing.T) {
	orgID := uuid.New()

	tmpDir, err := os.MkdirTemp(TempDir, "keys")
	defer func() { assert.Nil(t, os.RemoveAll(tmpDir)) }()

	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.public", []byte("mykey.public"), 0o755)
	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.private", []byte("mykey.private"), 0o755)
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/myexistingkey.public")
	assert.FileExists(t, tmpDir+"/myexistingkey.private")

	_, err = RunCliWithArgs(t, []string{"generate", "api-key", "--organization", orgID.String(), "--keys-folder", tmpDir, "--key-name", "myexistingkey"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit status 1")
}

func TestEncryptionKeygenDetectExistingKey(t *testing.T) {
	orgID := uuid.New()
	userID := uuid.New()

	tmpDir, err := os.MkdirTemp(TempDir, "encryption-keys")
	defer func() { assert.Nil(t, os.RemoveAll(tmpDir)) }()

	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.public", []byte("mykey.public"), 0o755)
	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.private", []byte("mykey.private"), 0o755)
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/myexistingkey.public")
	assert.FileExists(t, tmpDir+"/myexistingkey.private")

	_, err = RunCliWithArgs(t, []string{"generate", "encryption-key", "--organization", orgID.String(), "--user", userID.String(), "--encryption-keys-folder", tmpDir, "--encryption-key-name", "myexistingkey"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit status 1")
}

func TestStamp(t *testing.T) {
	orgID := uuid.New()

	out, err := RunCliWithArgs(t, []string{"request", "--no-post", "--keys-folder", ".", "--organization", orgID.String(), "--key-name", "fixtures/testkey.private", "--body", "hello!"})
	assert.Nil(t, err)

	var parsedOut map[string]string

	assert.Nil(t, json.Unmarshal([]byte(out), &parsedOut))

	stamp := parsedOut["stamp"]

	pubkeyBytes, err := os.ReadFile("fixtures/testkey.public")
	assert.Nil(t, err)

	ensureValidStamp(t, stamp, string(pubkeyBytes))
}

func TestApproveRequest(t *testing.T) {
	orgID := uuid.New()

	out, err := RunCliWithArgs(t, []string{"request", "--no-post", "--keys-folder", ".", "--organization", orgID.String(), "--key-name", "fixtures/testkey.private", "--body", "{\"some\": \"field\"}", "--path", "/some/endpoint"})
	assert.Nil(t, err)

	var parsedOut map[string]string
	err = json.Unmarshal([]byte(out), &parsedOut)
	assert.Nil(t, err)

	stamp := parsedOut["stamp"]
	pubkeyBytes, err := os.ReadFile("fixtures/testkey.public")
	assert.Nil(t, err)
	ensureValidStamp(t, stamp, string(pubkeyBytes))

	assert.Equal(t, "{\"some\": \"field\"}", parsedOut["message"])

	assert.Contains(t, parsedOut["curlCommand"], "curl -X POST -d'{\"some\": \"field\"}'")
	assert.Contains(t, parsedOut["curlCommand"], fmt.Sprintf("-H'X-Stamp: %s'", stamp))
	assert.Contains(t, parsedOut["curlCommand"], "https://api.turnkey.com/some/endpoint")
}

func ensureValidStamp(t *testing.T, stamp string, expectedPublicKey string) {
	stampBytes, err := base64.RawURLEncoding.DecodeString(stamp)
	assert.Nil(t, err)

	var parsedStamp *apikey.APIStamp

	assert.Nil(t, json.Unmarshal(stampBytes, &parsedStamp))

	assert.Equal(t, expectedPublicKey, parsedStamp.PublicKey)

	// All signatures start with 30....
	assert.True(t, strings.HasPrefix(parsedStamp.Signature, "30"))

	_, err = hex.DecodeString(parsedStamp.Signature)

	// Ensure there is no issue decoding the signature as a hexadecimal string
	assert.Nil(t, err)
}
