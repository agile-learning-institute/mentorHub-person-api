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
	// Setup companies
	cfg := config.NewConfig()
	mongoIO := config.NewMongoIO(cfg)

	// Connect to the database, and load static-ish data
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	mongoIO.LoadVersions()
	mongoIO.LoadEnumerators()

	// Setup store and handlers
	configHandler := handlers.NewConfigHandler(cfg, mongoIO)
	personStore := stores.NewPersonStore(mongoIO)
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
	port := cfg.GetPort()
	log.Printf("INFO: API Server Version %s", cfg.ApiVersion)
	log.Printf("INFO: Server Listening at %s", port)
	err := http.ListenAndServe(port, gorillaHandlers.CORS(originsOk, headersOk, methodsOk)(gorillaRouter))
	if err != nil {
		log.Println("ERROR: Server Ending with error", err)
	}
}
