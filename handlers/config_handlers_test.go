package handlers

import (
	"institute-person-api/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigHandler(t *testing.T) {
	// Setup a config
	config := config.NewConfig()
	configHandler := NewConfigHandler(*config)

	// Examine the result
	assert.NotNil(t, configHandler)
}

func TestGetConfig(t *testing.T) {
	// Setup
	config := config.NewConfig()
	configHandler := NewConfigHandler(*config)
	request := httptest.NewRequest("GET", "/config/", nil)
	responseRecorder := httptest.NewRecorder()
	// jsonString, _ := json.Marshal(config)

	// Invoke NewPerson
	configHandler.GetConfig(responseRecorder, request)

	// Examine the result
	assert.NotNil(t, configHandler)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	// assert.Equal(t, config, jsonString, responseRecorder.Body)
}
