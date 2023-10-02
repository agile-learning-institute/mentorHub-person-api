package main

import (
	"fmt"
	"net/http"

	"institute-person-api/config"
	"institute-person-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	config := config.NewConfig()
	handler := handlers.NewHandler(config)
	gorillaRouter := mux.NewRouter()

	gorillaRouter.HandleFunc("/api/person/", handler.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/person/{id}", handler.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/person/{id}", handler.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", handler.GetConfig).Methods("GET")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", gorillaRouter)
	defer config.Disconnect()
}
