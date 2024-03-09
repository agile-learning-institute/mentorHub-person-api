package handlers

import (
	"encoding/json"
	"log"
	"mentorhub-person-api/src/config"
	"net/http"

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

func (handler *ConfigHandler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get Config", correltionId)
	defer log.Printf("End CID: %s Get Config", correltionId)

	// Return the Config object as JSON
	err := handler.config.LoadLists()
	if err != nil {
		log.Printf("ERROR CID: %s ERROR %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(handler.config)
}
