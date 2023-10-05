package main

import (
	"log"
	"net/http"

	"institute-person-api/config"
	"institute-person-api/handlers"
	"institute-person-api/models"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Setup the ConfigHandler
	config := config.NewConfig()
	configHandler := handlers.NewConfigHandler()

	// Setup the PersonHandler and Store
	personStore := models.NewPersonStore()
	defer personStore.Disconnect()
	person := models.NewPerson(&personStore)
	personHandler := handlers.NewPersonHandler(person)

	// Setup the HttpServer Router
	gorillaRouter := mux.NewRouter()
	// gorillaRouter.Use(loggingMiddleware)

	// Configure cors filters
	// originsOk := gorillaHandlers.AllowedOrigins([]string{"*"})
	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := gorillaHandlers.AllowedOrigins([]string{"http://localhost:8080"}) // Your frontend's origin
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Define the Routes
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/", personHandler.GetPeople).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")

	// Start the server with Cors handler

	log.Printf("INFO: Server Version %s", config.Version)
	log.Printf("INFO: Server Listening at %s", config.Port)
	http.ListenAndServe(":8081", gorillaHandlers.CORS(originsOk, headersOk, methodsOk)(gorillaRouter))
}
