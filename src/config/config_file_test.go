package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test to verify config folder retrieval
func TestFileConfigFolder(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "../../test/config_files/")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, "../../test/config_files/", cfg.GetConfigFolder())
	assertConfigItem(t, cfg, "CONFIG_FOLDER", "environment", "../../test/config_files/")
}

// Test to verify port configuration
func TestFilePort(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "../../test/config_files/")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, "FooBar", cfg.GetPort())
	assertConfigItem(t, cfg, "PORT", "file", "FooBar")
}

// Test to verify patch level configuration
func TestFilePatch(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "../../test/config_files/")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, "LocalDev", cfg.GetPatch())
	assertConfigItem(t, cfg, "PATCH_LEVEL", "default", "LocalDev")
}

// Test to verify connection string configuration
func TestFileConnectionString(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "../../test/config_files/")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, "FooBar", cfg.GetConnectionString())
	assertConfigItem(t, cfg, "CONNECTION_STRING", "file", "Secret")
}

// Test to verify database name configuration
func TestFileDatabaseName(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "../../test/config_files/")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, "FooBar", cfg.GetDatabaseName())
	assertConfigItem(t, cfg, "DATABASE_NAME", "file", "FooBar")
}

// Test to verify database timeout configuration
func TestFileDatabaseTimeout(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "../../test/config_files/")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, 555, cfg.GetDatabaseTimeout())
	assertConfigItem(t, cfg, "CONNECTION_TIMEOUT", "file", "555")
}
