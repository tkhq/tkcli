package main_test

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	tmpDir, err := ioutil.TempDir("/tmp", "keys")
	defer os.RemoveAll(tmpDir)
	assert.Nil(t, err)

	out, err := RunCliWithArgs(t, []string{"gen", "--keys-folder", tmpDir, "--name", "mykey"})
	assert.Nil(t, err)

	assert.FileExists(t, tmpDir+"/mykey.public")
	assert.FileExists(t, tmpDir+"/mykey.private")

	var parsedOut map[string]string
	err = json.Unmarshal([]byte(out), &parsedOut)
	assert.Nil(t, err)
	assert.Equal(t, parsedOut["publicKeyFile"], tmpDir+"/mykey.public")
	assert.Equal(t, parsedOut["privateKeyFile"], tmpDir+"/mykey.private")
}

func TestSign(t *testing.T) {

	out, err := RunCliWithArgs(t, []string{"sign", "--key", "fixtures/testkey.private", "--message", "hello!"})
	assert.Nil(t, err)

	var parsedOut map[string]string
	err = json.Unmarshal([]byte(out), &parsedOut)
	assert.Nil(t, err)
	signature := parsedOut["signature"]

	// All signatures start with 30....
	assert.True(t, strings.HasPrefix(signature, "30"))

	_, err = hex.DecodeString(signature)
	// Ensure there is no issue decoding the signature as a hexadecimal string
	assert.Nil(t, err)
}
