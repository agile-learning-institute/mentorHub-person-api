package handlers

import (
	"encoding/json"
	"io"
	"log"
	"mentorhub-person-api/src/models"
	"mentorhub-person-api/src/stores"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PersonHandler struct {
	PersonStore *stores.PersonStore
}

func NewPersonHandler(personStore *stores.PersonStore) *PersonHandler {
	handler := &PersonHandler{}
	handler.PersonStore = personStore
	return handler
}

func (handler *PersonHandler) AddPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Add Person", correltionId)
	defer log.Printf("End CID: %s Add Person", correltionId)

	// Read the request body
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("ERROR CID: %s Bad Body Read %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Build the breadcrumb
	crumb := models.NewBreadCrumb(request.RemoteAddr, "SOME-USER-ID", correltionId.String())

	// Insert the new person document
	newPerson, err := handler.PersonStore.Insert(body, crumb)
	if err != nil {
		log.Printf("ERROR CID: %s PostPerson %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(newPerson)
}

func (handler *PersonHandler) UpdatePerson(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Update Person", correltionId)
	defer log.Printf("End CID: %s Update Person", correltionId)

	// Get the Person ID from the path
	id := mux.Vars(request)["id"]

	// Get the Request Body
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("ERROR CID: %s Bad Body Read %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Build the breadcrumb
	crumb := models.NewBreadCrumb(request.RemoteAddr, "SOME-USER-ID", correltionId.String())

	// Update the person
	updatedPerson, err := handler.PersonStore.FindOneAndUpdate(id, body, crumb)
	if err != nil {
		log.Printf("ERROR CID: %s Bad PatchPerson %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the updated Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(updatedPerson)
}
