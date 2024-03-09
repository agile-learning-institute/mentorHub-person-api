// This application implements a simple Person API
package main

import (
	"log"
	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/handlers"
	"mentorhub-person-api/src/stores"
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"
)

func main() {
	// Setup the Config
	config := config.NewConfig()

	// Connect to the database
	config.Connect()
	defer config.Disconnect()

	// Setup the Stores
	personStore := stores.NewPersonStore(config)

	// Setup the Handlers
	configHandler := handlers.NewConfigHandler(config)
	personHandler := handlers.NewPersonHandler(personStore)

	// Setup the HttpServer Router
	gorillaRouter := mux.NewRouter()

	// Setup the Promethius health middleware
	instrumentation := muxprom.NewDefaultInstrumentation()
	gorillaRouter.Use(instrumentation.Middleware)

	// Configure cors filters
	originsOk := gorillaHandlers.AllowedOrigins([]string{"*"})
	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"})

	// Define the Routes
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/", personHandler.GetPeople).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")
	gorillaRouter.Path("/api/health/").Handler(promhttp.Handler())

	// Start the server with Cors handler
	port := config.GetPort()
	log.Printf("INFO: API Server Version %s", config.ApiVersion)
	log.Printf("INFO: Server Listening at %s", port)
	err := http.ListenAndServe(port, gorillaHandlers.CORS(originsOk, headersOk, methodsOk)(gorillaRouter))
	if err != nil {
		log.Println("ERROR: Server Ending with error", err)
	}
}
