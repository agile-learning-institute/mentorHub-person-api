package handlers

import (
	"institute-person-api/src/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPersonHandler(t *testing.T) {
	config := config.NewConfig()
	assert.NotNil(t, config)
}

func TestAddPerson(t *testing.T) {
	config := config.NewConfig()
	assert.NotNil(t, config)
}

func TestUpdatePerson(t *testing.T) {
	config := config.NewConfig()
	assert.NotNil(t, config)
}
