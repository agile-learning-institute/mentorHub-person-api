package tests

import (
	"institute-person-api/src/main/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMongoHandler(t *testing.T) {
	config := config.NewConfig()
	assert.NotNil(t, config)
}

func TestGetAll(t *testing.T) {
	config := config.NewConfig()
	assert.NotNil(t, config)
}
func TestGetOne(t *testing.T) {
	config := config.NewConfig()
	assert.NotNil(t, config)
}
