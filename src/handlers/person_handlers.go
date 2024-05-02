/********************************************************************************
** Person Handlers
**    This class proivdes the implementation for all /person endpoints
********************************************************************************/
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
	PersonStore stores.PersonStoreInterface
}

func NewPersonHandler(personStore stores.PersonStoreInterface) *PersonHandler {
	handler := &PersonHandler{}
	handler.PersonStore = personStore
	return handler
}

/********************************************************************************
***** POST - Add A Person
********************************************************************************/
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
	crumb, err := models.NewBreadCrumb(request.RemoteAddr, "000000000000000000000000", correltionId.String())
	if err != nil {
		log.Printf("ERROR CID: %s Breadcrumb Constructuor failed %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new person document
	id, err := handler.PersonStore.Insert(body, crumb)
	if err != nil {
		log.Printf("ERROR CID: %s Insert Person Failed %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the Inserted Document from the database
	newPerson, err := handler.PersonStore.FindId(id)
	if err != nil {
		log.Printf("ERROR CID: %s Get Inserted Document Failed %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(newPerson)
}

/********************************************************************************
***** PATCH - Update A Person
********************************************************************************/
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
	crumb, err := models.NewBreadCrumb(request.RemoteAddr, "000000000000000000000000", correltionId.String())
	if err != nil {
		log.Printf("ERROR CID: %s Breadcrumb Constructuor failed %s", correltionId, err.Error())
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the person
	updatedPerson, err := handler.PersonStore.UpdateId(id, body, crumb)
	if err != nil {
		log.Printf("ERROR CID: %s Bad PatchPerson %s", correltionId, err.Error())
		log.Printf("Request Body: %s ", body)
		responseWriter.Header().Add("CorrelationId", correltionId.String())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the updated Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(updatedPerson)
}

/********************************************************************************
***** GET - Get A Person by ID
********************************************************************************/
func (handler *PersonHandler) GetPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get Person by ID", correltionId)
	defer log.Printf("End CID: %s Get Person by ID", correltionId)

	// Get the Document ID from the path
	id := mux.Vars(request)["id"]

	// Get the Document from the database
	results, err := handler.PersonStore.FindId(id)
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

/********************************************************************************
***** GET - Get a list of Name, PersonID
********************************************************************************/
func (handler *PersonHandler) GetPeople(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := uuid.NewRandom()
	log.Printf("Begin CID: %s Get All People", correltionId)
	defer log.Printf("End CID: %s Get All People", correltionId)

	// Get all the people
	results, err := handler.PersonStore.FindNames()
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
