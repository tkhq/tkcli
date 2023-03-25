package clifs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkhq/tkcli/internal/clifs"
)

// MacOSX has $HOME set by default
func TestGetKeyDirPathMacOSX(t *testing.T) {
	os.Setenv("HOME", "/home/dir")
	defer os.Unsetenv("HOME")

	// Need to unset this explicitly: the test runner has this set by default!
	originalValue := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalValue)

	dir := clifs.DefaultKeysDir()
	assert.Equal(t, dir, "/home/dir/.config/turnkey/keys")
}

// On UNIX, we expect XDG_CONFIG_HOME to be set
// If it's not set, we're back to a MacOSX-like system
func TestGetKeyDirPathUnix(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "/special/dir")
	defer os.Unsetenv("XDG_CONFIG_HOME")

	os.Setenv("HOME", "/home/dir")
	defer os.Unsetenv("HOME")

	dir := clifs.DefaultKeysDir()
	assert.Equal(t, dir, "/special/dir/turnkey/keys")
}

// If calling with a path, we should get this back if the path exists
// If not we should get an error
func TestGetKeyDirPathOverride(t *testing.T) {
   tmpDir := os.TempDir() //nolint:staticcheck

   defer os.RemoveAll(tmpDir) //nolint:staticcheck

	assert.NotNil(t, clifs.SetKeysDirectory("/does/not/exist"))

	assert.Nil(t, clifs.SetKeysDirectory(tmpDir))
}
