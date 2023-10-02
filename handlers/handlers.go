package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"institute-person-api/config"
	"institute-person-api/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{config: config}
}

func (h *Handler) AddPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Decode the request body to get the new Person details
	var newPerson models.Person
	err := json.NewDecoder(request.Body).Decode(&newPerson)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the Person to the database
	people := h.config.GetPeopleCollection()
	ctx, cancel := h.config.GetTimeoutContext()
	defer cancel()
	result, err := people.InsertOne(ctx, newPerson)

	// Get the new document
	query := bson.M{"_id": result.InsertedID}
	ctx, cancel = h.config.GetTimeoutContext()
	defer cancel()
	err = people.FindOne(ctx, query).Decode(&newPerson)

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(newPerson)
}

func (h *Handler) GetPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Get the Person from the database
	id := mux.Vars(request)["id"]
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	people := h.config.GetPeopleCollection()
	var getPerson models.Person
	ctx, cancel := h.config.GetTimeoutContext()
	defer cancel()
	people.FindOne(ctx, query).Decode(&getPerson)

	// Return the Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(getPerson)
}

func (h *Handler) UpdatePerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Decode the request body to get the updated Person details
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	var updateValues bson.M
	json.Unmarshal(body, &updateValues)

	// Update the Person with provided values
	people := h.config.GetPeopleCollection()
	id := mux.Vars(request)["id"]
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": updateValues}

	ctx, cancel := h.config.GetTimeoutContext()
	defer cancel()
	var updatedPerson models.Person
	people.FindOneAndUpdate(ctx, query, update, options).Decode(&updatedPerson)
	log.Println("Updated Document: ", updatedPerson)

	// Return the updated Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(updatedPerson)
}

func (h *Handler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(h.config.ToJSONStruct())
}
