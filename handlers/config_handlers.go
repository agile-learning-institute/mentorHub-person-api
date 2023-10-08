package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"institute-person-api/config"

	"github.com/google/uuid"
)

type ConfigHandler struct {
	config *config.Config
}

func NewConfigHandler(theConfig *config.Config) *ConfigHandler {
	this := &ConfigHandler{}
	this.config = theConfig
	return this
}

func (h *ConfigHandler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("TRANSACTION CID: %s Get Config Start", correltionId)
	defer log.Printf("TRANSACTION CID: %s Get Config Complete", correltionId)

	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(h.config)
}
