package main

import (
	"fmt"
	"net/http"

	"institute-person-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// TODO: Initiate mongo connection pool

	gorillaRouter := mux.NewRouter()

	gorillaRouter.HandleFunc("/api/product/", handlers.AddPerson).Methods("POST")
	gorillaRouter.HandleFunc("/api/product/{id}", handlers.GetPerson).Methods("GET")
	gorillaRouter.HandleFunc("/api/product/{id}", handlers.UpdatePerson).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/config/", handlers.GetConfig).Methods("GET")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", gorillaRouter)

	// TODO: Housekeep mongo connection pool
}
