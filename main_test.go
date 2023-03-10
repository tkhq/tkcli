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

	"github.com/stretchr/testify/assert"
	"github.com/tkhq/tkcli/internal/apikey"
)

const TURNKEY_BINARY_NAME = "turnkey"

func RunCliWithArgs(t *testing.T, args []string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(path.Join(currentDir, "build", TURNKEY_BINARY_NAME), args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func TestHelpText(t *testing.T) {
	out, err := RunCliWithArgs(t, []string{})
	assert.Nil(t, err)
	assert.Contains(t, out, "The Turnkey CLI")
	assert.Contains(t, out, "USAGE")
	assert.Contains(t, out, "COMMANDS")
}

func TestKeygenArgValidation(t *testing.T) {
	out, err := RunCliWithArgs(t, []string{"gen"})
	assert.Equal(t, err.Error(), "exit status 1")

	assert.Contains(t, out, "Required flag \"name\" not set")
}

func TestKeygenInTmpFolder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("/tmp", "keys")
	defer os.RemoveAll(tmpDir)

	assert.Nil(t, err)

	out, err := RunCliWithArgs(t, []string{"gen", "--keys-folder", tmpDir, "--name", "mykey"})
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/mykey.public")
	assert.FileExists(t, tmpDir+"/mykey.private")

	publicKeyData, err := os.ReadFile(tmpDir + "/mykey.public")
	assert.Nil(t, err)

	var parsedOut map[string]string
	err = json.Unmarshal([]byte(out), &parsedOut)
	assert.Nil(t, err)
	assert.Equal(t, parsedOut["publicKey"], string(publicKeyData))
	assert.Equal(t, parsedOut["publicKeyFile"], tmpDir+"/mykey.public")
	assert.Equal(t, parsedOut["privateKeyFile"], tmpDir+"/mykey.private")
}

func TestKeygenDetectExistingKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("/tmp", "keys")
	defer os.RemoveAll(tmpDir)

	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.public", []byte("mykey.public"), 0755)
	assert.Nil(t, err)

	err = os.WriteFile(tmpDir+"/myexistingkey.private", []byte("mykey.private"), 0755)
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/myexistingkey.public")
	assert.FileExists(t, tmpDir+"/myexistingkey.private")

	_, err = RunCliWithArgs(t, []string{"gen", "--keys-folder", tmpDir, "--name", "myexistingkey"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit status 1")
}

func TestStamp(t *testing.T) {
	out, err := RunCliWithArgs(t, []string{"stamp", "--key", "fixtures/testkey.private", "--message", "hello!"})
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
	out, err := RunCliWithArgs(t, []string{"approve-request", "--host", "api.turnkey.io", "--key", "fixtures/testkey.private", "--body", "{\"some\": \"field\"}", "--path", "/some/endpoint"})
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
	json.Unmarshal(stampBytes, &parsedStamp)

	assert.Equal(t, expectedPublicKey, parsedStamp.PublicKey)

	// All signatures start with 30....
	assert.True(t, strings.HasPrefix(parsedStamp.Signature, "30"))

	_, err = hex.DecodeString(parsedStamp.Signature)
	// Ensure there is no issue decoding the signature as a hexadecimal string
	assert.Nil(t, err)
}
