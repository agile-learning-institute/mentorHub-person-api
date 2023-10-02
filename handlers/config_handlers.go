package handlers

import (
	"encoding/json"
	"net/http"

	"institute-person-api/config"
)

type ConfigHandler struct {
	config *config.Config
}

func NewConfigHandler() *ConfigHandler {
	this := &ConfigHandler{}
	this.config = config.NewConfig()
	return &ConfigHandler{}
}

func (h *ConfigHandler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(h.config)
}
