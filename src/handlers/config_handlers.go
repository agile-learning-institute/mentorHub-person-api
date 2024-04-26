package handlers

import (
	"encoding/json"
	"log"
	"mentorhub-person-api/src/config"
	"net/http"

	"github.com/google/uuid"
)

type ConfigHandler struct {
	config  *config.Config
	mongoIO config.MongoIOInterface
}

func NewConfigHandler(theConfig *config.Config, theMongoIO config.MongoIOInterface) *ConfigHandler {
	this := &ConfigHandler{}
	this.config = theConfig
	this.mongoIO = theMongoIO
	return this
}

/********************************************************************************
***** GET - Get the Configuration Information (includes Mentor and Partner lists)
********************************************************************************/
func (handler *ConfigHandler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get Config", correltionId)
	defer log.Printf("End CID: %s Get Config", correltionId)

	// Refresh Mentors list
	err := handler.mongoIO.FetchMentors()
	if err != nil {
		log.Printf("ERROR CID: %s Fetch Mentors Error: %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Refresh Partnes list
	err = handler.mongoIO.FetchPartners()
	if err != nil {
		log.Printf("ERROR CID: %s Fetch Partners Error: %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the Config object
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(handler.config)
}
