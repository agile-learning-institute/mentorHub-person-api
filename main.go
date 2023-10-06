// This application implements a simple Person API
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

	// Setup the PersonHandler and Store
	var personStore models.PersonStoreInterface
	personStore = models.NewPersonStore(config)
	configHandler := handlers.NewConfigHandler(config)
	defer personStore.Disconnect()
	var person models.PersonInterface
	person = models.NewPerson(personStore)
	personHandler := handlers.NewPersonHandler(person)

	// Setup the HttpServer Router
	gorillaRouter := mux.NewRouter()

	// Configure cors filters
	originsOk := gorillaHandlers.AllowedOrigins([]string{"*"})
	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// originsOk := gorillaHandlers.AllowedOrigins([]string{"http://localhost:8080"}) // Your frontend's origin
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"})

	// Define the Routes
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/", personHandler.GetPeople).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")

	// Start the server with Cors handler
	port := config.GetPort()
	log.Printf("INFO: Server Version %s", config.Version)
	log.Printf("INFO: Server Listening at %s", port)
	err := http.ListenAndServe(port, gorillaHandlers.CORS(originsOk, headersOk, methodsOk)(gorillaRouter))
	if err != nil {
		log.Println("ERROR: Server Ending with error", err)
	}
}
