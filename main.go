package main

import (
	"fmt"
	"net/http"

	"product-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	gorillaRouter := mux.NewRouter()

	gorillaRouter.HandleFunc("/api/product/{id}", handlers.UpdateProduct).Methods("PATCH")
	gorillaRouter.HandleFunc("/api/product/", handlers.AddProduct).Methods("POST")
	gorillaRouter.HandleFunc("/api/config/", handlers.GetConfig).Methods("GET")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", gorillaRouter)
}
