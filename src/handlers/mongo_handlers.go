package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"mentorhub-person-api/src/stores"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MongoHandler struct {
	MongoStore *stores.MongoStore
}

func NewMongoHandler(mongoReadStore *stores.MongoStore) *MongoHandler {
	handler := &MongoHandler{}
	handler.MongoStore = mongoReadStore
	return handler
}

func (handler *MongoHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	collection := handler.MongoStore.CollectionName
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get All from %s", correltionId, collection)
	defer log.Printf("End CID: %s Get All from %s", correltionId, collection)

	// Get all the people
	results, err := handler.MongoStore.FindDocuments()
	if err != nil {
		log.Printf("ERROR CID: %s GetAllNames %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(results)
}

func (handler *MongoHandler) GetNames(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	collection := handler.MongoStore.CollectionName
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get All from %s", correltionId, collection)
	defer log.Printf("End CID: %s Get All from %s", correltionId, collection)

	// Get all the people
	results, err := handler.MongoStore.FindNames()
	if err != nil {
		log.Printf("ERROR CID: %s GetAllNames %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(results)
}

func (handler *MongoHandler) GetOne(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	collection := handler.MongoStore.CollectionName
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get %s by ID", correltionId, collection)
	defer log.Printf("End CID: %s Get %s by ID", correltionId, collection)

	// Get the Document ID from the path
	id := mux.Vars(request)["id"]

	// Get the Document from the database
	results, err := handler.MongoStore.FindId(id)
	if err != nil {
		log.Printf("ERROR CID: %s ERROR %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(results)
}
