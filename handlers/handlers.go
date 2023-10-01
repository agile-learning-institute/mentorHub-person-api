package handlers

import (
	"encoding/json"
	"net/http"

	"institute-person-api/models"
)

func AddPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Decode the request body to get the new Person details
	var newPerson models.Person
	err := json.NewDecoder(request.Body).Decode(&newPerson)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Add the Person to your database or data store
	// newPerson = connection.insertOne(newPerson);

	// Return the new Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(newPerson)
}

func GetPerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Decode the request body to get the new Person details
	var getPerson models.Person

	// TODO: Get the Person from the database or data store
	// id := mux.Vars(request)["id"]
	// query := {_id: new ObjectID(id)}
	// getPerson = connection.findOne(query);

	// Return the Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(getPerson)
}

func UpdatePerson(responseWriter http.ResponseWriter, request *http.Request) {
	// Decode the request body to get the updated Person details
	var updatedPerson models.Person
	err := json.NewDecoder(request.Body).Decode(&updatedPerson)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Update the Person in your database or data store
	// Extract the Person ID from the URL parameters
	// id := mux.Vars(request)["id"]
	// query := {_id:new ObjectId(id)}
	// update := {$set{updatePerson}}
	// options := {returnAfter:true}
	// updatePerson = connection.findOneAndUpdate(query, update, options)

	// Return the updated Person as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(updatedPerson)
}

func GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// Create a Config object
	config := models.Config{APIVersion: "1.0", DataVersion: "1.0"}

	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(config)
}
