package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, DefaultConfigFolder, cfg.configFolder)
	assert.Equal(t, DefaultConnectionString, cfg.connectionString)
	assert.Equal(t, DefaultDatabaseName, cfg.databaseName)
	assert.Equal(t, DefaultPeopleCollectionName, cfg.peopleCollectionName)
	assert.Equal(t, DefaultTimeout, cfg.databaseTimeout)
	assert.Equal(t, DefaultPort, cfg.Port)
	assert.Equal(t, "LocalDev", cfg.patch)
	assert.Equal(t, "1.0.LocalDev", cfg.Version)

	assert.Equal(t, 7, len(cfg.ConfigItems))
	assert.Equal(t, "CONFIG_FOLDER", cfg.ConfigItems[0].Name)
	assert.Equal(t, "default", cfg.ConfigItems[0].From)
	assert.Equal(t, DefaultConfigFolder, cfg.ConfigItems[0].Value)
	assert.Equal(t, "CONNECTION_STRING", cfg.ConfigItems[1].Name)
	assert.Equal(t, "default", cfg.ConfigItems[1].From)
	assert.Equal(t, "Secret", cfg.ConfigItems[1].Value)
	assert.Equal(t, "DATABASE_NAME", cfg.ConfigItems[2].Name)
	assert.Equal(t, "default", cfg.ConfigItems[2].From)
	assert.Equal(t, DefaultDatabaseName, cfg.ConfigItems[2].Value)
	assert.Equal(t, "PEOPLE_COLLECTION_NAME", cfg.ConfigItems[3].Name)
	assert.Equal(t, "default", cfg.ConfigItems[3].From)
	assert.Equal(t, DefaultPeopleCollectionName, cfg.ConfigItems[3].Value)
	assert.Equal(t, "CONNECTION_TIMEOUT", cfg.ConfigItems[4].Name)
	assert.Equal(t, "default", cfg.ConfigItems[4].From)
	assert.Equal(t, "10", cfg.ConfigItems[4].Value)
	assert.Equal(t, "PATCH_LEVEL", cfg.ConfigItems[5].Name)
	assert.Equal(t, "default", cfg.ConfigItems[5].From)
	assert.Equal(t, "LocalDev", cfg.ConfigItems[5].Value)
	assert.Equal(t, "PORT", cfg.ConfigItems[6].Name)
	assert.Equal(t, "default", cfg.ConfigItems[6].From)
	assert.Equal(t, ":8080", cfg.ConfigItems[6].Value)
}

func TestGetConnectionString(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultConnectionString, cfg.GetConnectionString())
}

func TestGetDatabaseName(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultDatabaseName, cfg.GetDatabaseName())
}

func TestGetTimeoutContext(t *testing.T) {
	cfg := NewConfig()
	ctx, cancel := cfg.GetTimeoutContext()

	// Check deadline is set
	deadline, hasDeadline := ctx.Deadline()
	assert.True(t, hasDeadline)
	assert.NotNil(t, deadline)

	// Check cancel function
	cancel()
	assert.Equal(t, context.Canceled, ctx.Err())
}
func TestGetPeopleCollectionName(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, DefaultPeopleCollectionName, cfg.GetPeopleCollectionName())
}
