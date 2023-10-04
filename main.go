package main

import (
	"log"
	"net/http"

	"institute-person-api/config"
	"institute-person-api/handlers"
	"institute-person-api/models"

	"github.com/gorilla/mux"
)

func main() {
	// Setup the ConfigHandler
	config := config.NewConfig()
	log.Printf("Server Starting at port %s", config.Port)
	log.Printf("Using connection string %s", config.GetConnectionString())
	configHandler := handlers.NewConfigHandler()

	// Setup the PersonHandler and Store
	personStore := models.NewPersonStore()
	defer personStore.Disconnect()
	person := models.NewPerson(&personStore)
	personHandler := handlers.NewPersonHandler(person)

	// Setup the HttpServer Router
	gorillaRouter := mux.NewRouter()

	// Define the Routes
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/", personHandler.GetPeople).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")

	// Start the server
	http.ListenAndServe(config.Port, gorillaRouter)
}
