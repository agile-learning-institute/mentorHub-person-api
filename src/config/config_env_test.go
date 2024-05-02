package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test to verify port configuration
func TestEnvPort(t *testing.T) {
	os.Setenv("PORT", "foo")
	defer os.Unsetenv("PORT")

	cfg := NewConfig()
	assert.Equal(t, "foo", cfg.GetPort())
	assertConfigItem(t, cfg, "PORT", "environment", "foo")
}

// Test to verify patch level configuration
func TestEnvPatch(t *testing.T) {
	os.Setenv("PATCH_LEVEL", "foo")
	defer os.Unsetenv("PATCH_LEVEL")

	cfg := NewConfig()
	assert.Equal(t, "foo", cfg.GetPatch())
	assertConfigItem(t, cfg, "PATCH_LEVEL", "environment", "foo")
}

// Test to verify config folder retrieval
func TestEnvConfigFolder(t *testing.T) {
	os.Setenv("CONFIG_FOLDER", "foo")
	defer os.Unsetenv("CONFIG_FOLDER")

	cfg := NewConfig()
	assert.Equal(t, "foo", cfg.GetConfigFolder())
	assertConfigItem(t, cfg, "CONFIG_FOLDER", "environment", "foo")
}

// Test to verify connection string configuration
func TestEnvConnectionString(t *testing.T) {
	os.Setenv("CONNECTION_STRING", "foo")
	defer os.Unsetenv("CONNECTION_STRING")

	cfg := NewConfig()
	assert.Equal(t, "foo", cfg.GetConnectionString())
	assertConfigItem(t, cfg, "CONNECTION_STRING", "environment", "Secret")
}

// Test to verify database name configuration
func TestEnvDatabaseName(t *testing.T) {
	os.Setenv("DATABASE_NAME", "foo")
	defer os.Unsetenv("DATABASE_NAME")

	cfg := NewConfig()
	assert.Equal(t, "foo", cfg.GetDatabaseName())
	assertConfigItem(t, cfg, "DATABASE_NAME", "environment", "foo")
}

// Test to verify database timeout configuration
func TestEnvDatabaseTimeout(t *testing.T) {
	os.Setenv("CONNECTION_TIMEOUT", "99")
	defer os.Unsetenv("CONNECTION_TIMEOUT")

	cfg := NewConfig()
	assert.Equal(t, 99, cfg.GetDatabaseTimeout())
	assertConfigItem(t, cfg, "CONNECTION_TIMEOUT", "environment", "99")
}
