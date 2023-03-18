package clifs_test

import (
	"io/ioutil"
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

	dir, err := clifs.GetKeyDirPath("")
	assert.Nil(t, err)
	assert.Equal(t, dir, "/home/dir/.config/turnkey/keys")
}

// On UNIX, we expect XDG_CONFIG_HOME to be set
// If it's not set, we're back to a MacOSX-like system
func TestGetKeyDirPathUnix(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "/special/dir")
	defer os.Unsetenv("XDG_CONFIG_HOME")

	os.Setenv("HOME", "/home/dir")
	defer os.Unsetenv("HOME")

	dir, err := clifs.GetKeyDirPath("")
	assert.Nil(t, err)
	assert.Equal(t, dir, "/special/dir/turnkey/keys")
}

// In the case where we don't have a $HOME defined, bail!
func TestGetKeyDirPathDysfunctionalOS(t *testing.T) {
	originalValue := os.Getenv("HOME")
	os.Unsetenv("HOME")
	defer os.Setenv("HOME", originalValue)

	dir, err := clifs.GetKeyDirPath("")
	assert.Equal(t, dir, "")
	assert.Equal(t, "error while reading user home directory: $HOME is not defined", err.Error())
}

// If calling with a path, we should get this back if the path exists
// If not we should get an error
func TestGetKeyDirPathOverride(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "keys")
	defer os.RemoveAll(tmpDir)
	assert.Nil(t, err)

	dir, err := clifs.GetKeyDirPath("/does/not/exist")
	assert.Equal(t, "Cannot put key files in /does/not/exist: stat /does/not/exist: no such file or directory", err.Error())
	assert.Equal(t, "", dir)

	dir, err = clifs.GetKeyDirPath(tmpDir)
	assert.Nil(t, err)
	assert.Equal(t, tmpDir, dir)
}
