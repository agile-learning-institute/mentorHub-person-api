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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"
)

func main() {
	// Setup the ConfigHandler
	config := config.NewConfig()

	// Setup the EnumeratorStore
	var enumStore *models.EnumeratorStore
	var err error
	enumStore, err = models.NewEnumeratorStore(config)
	if err != nil {
		log.Fatal("PersonStore Construction Error:", err)
	}

	// Setup the PersonStore
	var personStore models.PersonStoreInterface
	personStore, err = models.NewPersonStore(config)
	defer personStore.Disconnect()
	if err != nil {
		log.Fatal("PersonStore Construction Error:", err)
	}

	// Setup the Handlers
	personHandler := handlers.NewPersonHandler(personStore)
	configHandler := handlers.NewConfigHandler(config)
	enumHandler := handlers.NewEnumHandler(enumStore)

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
	gorillaRouter.Path("/api/health/").Handler(promhttp.Handler())
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/", personHandler.GetPeople).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")
	gorillaRouter.HandleFunc("/api/enums/", enumHandler.GetEnums).Methods("GET")

	// Start the server with Cors handler
	port := config.GetPort()
	log.Printf("INFO: API Server Version %s", config.ApiVersion)
	log.Printf("INFO: Server Listening at %s", port)
	err = http.ListenAndServe(port, gorillaHandlers.CORS(originsOk, headersOk, methodsOk)(gorillaRouter))
	if err != nil {
		log.Println("ERROR: Server Ending with error", err)
	}
}
