package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"institute-person-api/src/main/stores"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MongoHandler struct {
	MongoStore *stores.MongoStore
}

func NewMongoHandler(mongoStore *stores.MongoStore) *MongoHandler {
	this := &MongoHandler{}
	this.MongoStore = mongoStore
	return this
}

func (this *MongoHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get All from %s", correltionId, this.MongoStore.CollectionName)
	defer log.Printf("End CID: %s Get All from %s", correltionId, this.MongoStore.CollectionName)

	// Get all the people
	results, err := this.MongoStore.FindDocuments()
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

func (this *MongoHandler) GetOne(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get %s by ID", correltionId, this.MongoStore.CollectionName)
	defer log.Printf("End CID: %s Get %s by ID", correltionId, this.MongoStore.CollectionName)

	// Get the Document ID from the path
	id := mux.Vars(request)["id"]

	// Get the Document from the database
	results, err := this.MongoStore.FindId(id)
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
