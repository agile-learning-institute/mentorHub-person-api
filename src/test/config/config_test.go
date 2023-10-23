package tests

import (
	"context"
	"institute-person-api/src/main/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := config.NewConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "1.1.LocalDev", cfg.ApiVersion)

	assert.Equal(t, 6, len(cfg.ConfigItems))
	assert.Equal(t, "PATCH_LEVEL", cfg.ConfigItems[0].Name)
	assert.Equal(t, "default", cfg.ConfigItems[0].From)
	assert.Equal(t, "LocalDev", cfg.ConfigItems[0].Value)
	assert.Equal(t, "CONFIG_FOLDER", cfg.ConfigItems[1].Name)
	assert.Equal(t, "default", cfg.ConfigItems[1].From)
	assert.Equal(t, config.DefaultConfigFolder, cfg.ConfigItems[1].Value)
	assert.Equal(t, "CONNECTION_STRING", cfg.ConfigItems[2].Name)
	assert.Equal(t, "default", cfg.ConfigItems[2].From)
	assert.Equal(t, "Secret", cfg.ConfigItems[2].Value)
	assert.Equal(t, "DATABASE_NAME", cfg.ConfigItems[3].Name)
	assert.Equal(t, "default", cfg.ConfigItems[3].From)
	assert.Equal(t, config.DefaultDatabaseName, cfg.ConfigItems[3].Value)
	assert.Equal(t, "CONNECTION_TIMEOUT", cfg.ConfigItems[4].Name)
	assert.Equal(t, "default", cfg.ConfigItems[4].From)
	assert.Equal(t, "10", cfg.ConfigItems[4].Value)
	assert.Equal(t, "PORT", cfg.ConfigItems[5].Name)
	assert.Equal(t, "default", cfg.ConfigItems[5].From)
	assert.Equal(t, ":8080", cfg.ConfigItems[5].Value)
}

func TestDisconnect(t *testing.T) {
	cfg := config.NewConfig()
	assert.Equal(t, "Foo", cfg)
}

func TestGetPort(t *testing.T) {
	cfg := config.NewConfig()
	assert.Equal(t, config.DefaultPort, cfg.GetPort())
}

func TestGetCollection(t *testing.T) {
	cfg := config.NewConfig()
	assert.Equal(t, "Foo", cfg)
}

func TestAddConfigStore(t *testing.T) {
	cfg := config.NewConfig()
	assert.Equal(t, "Foo", cfg)
}

func TestGetTimeoutContext(t *testing.T) {
	cfg := config.NewConfig()
	ctx, cancel := cfg.GetTimeoutContext()

	// Check deadline is set
	deadline, hasDeadline := ctx.Deadline()
	assert.True(t, hasDeadline)
	assert.NotNil(t, deadline)

	// Check cancel function
	cancel()
	assert.Equal(t, context.Canceled, ctx.Err())
}
