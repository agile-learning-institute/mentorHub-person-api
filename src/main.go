// This application implements a simple Person API
package main

import (
	"log"
	"net/http"

	"institute-person-api/src/config"
	"institute-person-api/src/handlers"
	"institute-person-api/src/stores"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// Setup the Config
	config := config.NewConfig()

	// Connect to the database
	config.Connect()
	defer config.Disconnect()

	// Setup the Stores
	enumStore := stores.NewMongoStore(config, "enumerators", stores.MongoQueryNotVersion())
	mentorStore := stores.NewMongoStore(config, "people", bson.M{"mentor": true})
	partnerStore := stores.NewMongoStore(config, "partners", stores.MongoQueryNotVersion())
	readPersonStore := stores.NewMongoStore(config, "people", stores.MongoQueryNotVersion())
	personStore := stores.NewPersonStore(config)

	// Setup the Handlers
	configHandler := handlers.NewConfigHandler(config)
	enumHandler := handlers.NewMongoHandler(enumStore)
	mentorHandler := handlers.NewMongoHandler(mentorStore)
	partnerHandler := handlers.NewMongoHandler(partnerStore)
	readPersonHandler := handlers.NewMongoHandler(readPersonStore)
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
	gorillaRouter.HandleFunc("/api/person/", readPersonHandler.GetNames).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/", personHandler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/{id}", readPersonHandler.GetOne).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", personHandler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", configHandler.GetConfig).Methods("GET")
	gorillaRouter.HandleFunc("/api/enums/", enumHandler.GetAll).Methods("GET")
	gorillaRouter.HandleFunc("/api/partners/", partnerHandler.GetNames).Methods("GET")
	gorillaRouter.HandleFunc("/api/mentors/", mentorHandler.GetNames).Methods("GET")
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
