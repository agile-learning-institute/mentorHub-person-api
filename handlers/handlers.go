package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"institute-person-api/config"
	"institute-person-api/models"

	"github.com/gorilla/mux"
)

type Handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{config: config}
}

func (h *Handler) AddPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Read the request body
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new person document
	newPerson := models.PostPerson(body, h.config)

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(newPerson)
}

func (h *Handler) GetPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Get the Person ID from the path
	id := mux.Vars(request)["id"]

	// Get the Person from the database
	person := models.GetPerson(id, h.config)

	// Return the Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(person)
}

func (h *Handler) UpdatePerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Get the Person ID from the path
	id := mux.Vars(request)["id"]

	// Get the Request Body
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the person
	updatedPerson := models.PatchPerson(id, body, h.config)

	// Return the updated Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(updatedPerson)
}

func (h *Handler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(h.config.ToJSONStruct())
}
