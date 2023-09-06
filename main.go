package main

import (
	"fmt"
	"net/http"

	"product-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/product/{id}", handlers.UpdateProduct).Methods("PATCH")
	r.HandleFunc("/api/product/", handlers.AddProduct).Methods("POST")
	r.HandleFunc("/api/config/", handlers.GetConfig).Methods("GET")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}
