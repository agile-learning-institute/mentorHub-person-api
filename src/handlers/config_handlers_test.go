/********************************************************************************
** Config Handlers
**    This implementes a handler for the /config path
********************************************************************************/
package handlers

import (
	"mentorhub-person-api/src/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigHandler(t *testing.T) {
	// Setup a config
	cfg := config.NewConfig()
	mockMongoIO := config.NewMockMongoIO(cfg)
	configHandler := NewConfigHandler(cfg, mockMongoIO)

	// Examine the result
	assert.NotNil(t, configHandler)
}

func TestGetConfig(t *testing.T) {
	// Setup
	cfg := config.NewConfig()
	mongoIO := config.NewMockMongoIO(cfg)
	configHandler := NewConfigHandler(cfg, mongoIO)
	request := httptest.NewRequest("GET", "/config/", nil)
	responseRecorder := httptest.NewRecorder()

	// Invoke GetConfig
	configHandler.GetConfig(responseRecorder, request)

	// Examine the result
	assert.NotNil(t, configHandler)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	// assert.Equal(t, config, jsonString, responseRecorder.Body)
}
