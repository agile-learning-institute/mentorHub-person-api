package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test to verify API version retrieval
func TestGetApiVersion(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, VersionMajor+"."+VersionMinor+"."+"LocalDev", cfg.ApiVersion)
}

// Test to verify port configuration
func TestGetPort(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultPort, cfg.GetPort())
	assertConfigItem(t, cfg, "PORT", "default", DefaultPort)
}

// Test to verify patch level configuration
func TestGetPatch(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, "LocalDev", cfg.GetPatch())
	assertConfigItem(t, cfg, "PATCH_LEVEL", "default", "LocalDev")
}

// Test to verify config folder retrieval
func TestGetConfigFolder(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultConfigFolder, cfg.GetConfigFolder())
	assertConfigItem(t, cfg, "CONFIG_FOLDER", "default", DefaultConfigFolder)
}

// Test to verify connection string configuration
func TestGetConnectionString(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultConnectionString, cfg.GetConnectionString())
	assertConfigItem(t, cfg, "CONNECTION_STRING", "default", "Secret")
}

// Test to verify database name configuration
func TestGetDatabaseName(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultDatabaseName, cfg.GetDatabaseName())
	assertConfigItem(t, cfg, "DATABASE_NAME", "default", DefaultDatabaseName)
}

// Test to verify database timeout configuration
func TestGetDatabaseTimeout(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, 10, cfg.GetDatabaseTimeout())
	assertConfigItem(t, cfg, "CONNECTION_TIMEOUT", "default", DefaultTimeout)
}
