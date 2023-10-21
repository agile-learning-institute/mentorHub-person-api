package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	 "institute-person-api/src/config"

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
	log.Printf("Begin CID: %s Get Config", correltionId)
	defer log.Printf("End CID: %s Get Config", correltionId)

	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(h.config)
}
