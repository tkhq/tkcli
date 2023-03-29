package main_test

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tkhq/tkcli/internal/apikey"
)

const TurnkeyBinaryName = "turnkey"

func RunCliWithArgs(t *testing.T, args []string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(path.Join(currentDir, "..", "..", "build", TurnkeyBinaryName), args...)
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

func TestKeygenInTmpFolder(t *testing.T) {
	orgID := uuid.Must(uuid.NewV4())

	tmpDir, err := os.MkdirTemp("/tmp", "keys")
	defer os.RemoveAll(tmpDir)

	assert.Nil(t, err)

	out, err := RunCliWithArgs(t, []string{"gen", "--keys-folder", tmpDir, "--key-name", "mykey", "--organization", orgID.String()})
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

func TestKeygenDetectExistingKey(t *testing.T) {
	orgID := uuid.Must(uuid.NewV4())

	tmpDir, err := os.MkdirTemp("/tmp", "keys")
	defer os.RemoveAll(tmpDir)

	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.public", []byte("mykey.public"), 0755)
	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.private", []byte("mykey.private"), 0755)
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/myexistingkey.public")
	assert.FileExists(t, tmpDir+"/myexistingkey.private")

	_, err = RunCliWithArgs(t, []string{"gen", "--organization", orgID.String(), "--keys-folder", tmpDir, "--key-name", "myexistingkey"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit status 1")
}

func TestStamp(t *testing.T) {
	out, err := RunCliWithArgs(t, []string{"request", "--no-post", "--key-name", "fixtures/testkey.private", "--body", "hello!"})
	assert.Nil(t, err)

	var parsedOut map[string]string
	err = json.Unmarshal([]byte(out), &parsedOut)
	assert.Nil(t, err)
	stamp := parsedOut["stamp"]

	pubkeyBytes, err := os.ReadFile("fixtures/testkey.public")
	assert.Nil(t, err)
	ensureValidStamp(t, stamp, string(pubkeyBytes))
}

func TestApproveRequest(t *testing.T) {
	out, err := RunCliWithArgs(t, []string{"request", "--no-post", "--host", "api.turnkey.io", "--key-name", "fixtures/testkey.private", "--body", "{\"some\": \"field\"}", "--path", "/some/endpoint"})
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
	assert.Contains(t, parsedOut["curlCommand"], "https://api.turnkey.io/some/endpoint")
}

func ensureValidStamp(t *testing.T, stamp string, expectedPublicKey string) {
	stampBytes, err := base64.RawURLEncoding.DecodeString(stamp)
	assert.Nil(t, err)

	var parsedStamp *apikey.ApiStamp

	assert.Nil(t, json.Unmarshal(stampBytes, &parsedStamp))

	assert.Equal(t, expectedPublicKey, parsedStamp.PublicKey)

	// All signatures start with 30....
	assert.True(t, strings.HasPrefix(string(parsedStamp.Signature), "30"))

	_, err = hex.DecodeString(parsedStamp.Signature)

	// Ensure there is no issue decoding the signature as a hexadecimal string
	assert.Nil(t, err)
}
