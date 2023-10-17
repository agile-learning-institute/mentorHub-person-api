package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"institute-person-api/models"

	"github.com/google/uuid"
)

type EnumHandler struct {
	enumerations []*models.Enumerator
}

func NewEnumHandler(enumeratorStore *models.EnumeratorStore) *EnumHandler {
	this := &EnumHandler{}
	this.enumerations = enumeratorStore.Enumerators
	return this
}

func (this *EnumHandler) GetEnums(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get ENums", correltionId)
	defer log.Printf("End CID: %s Get Enums", correltionId)

	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(this.enumerations)
}
