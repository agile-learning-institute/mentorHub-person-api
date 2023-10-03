package main

import (
	"fmt"
	"net/http"

	"institute-person-api/handlers"
	"institute-person-api/models"

	"github.com/gorilla/mux"
)

func main() {
	// Setup the PersonHandler and Store
	personStore := models.NewPersonStore()
	defer personStore.Disconnect()
	person := models.NewPerson(personStore)
	personHandler := handlers.NewPersonHandler(person)

	// Setup the ConfigHandler
	configHandler := handlers.NewConfigHandler()

	// Setup the HttpServer Router
	gorillaRouter := mux.NewRouter()

	// Define the Routes
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/", personHandler.GetPeople).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")

	// Start the server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", gorillaRouter)
}
