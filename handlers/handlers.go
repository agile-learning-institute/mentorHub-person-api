package handlers

import (
	"encoding/json"
	"net/http"

	"product-api/models"
)

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL parameters
	// vars := mux.Vars(r)
	// id := vars["id"]

	// Decode the request body to get the updated product details
	var updatedProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Update the product in your database or data store

	// Return the updated product as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	// Decode the request body to get the new product details
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Add the product to your database or data store

	// Return the new product as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
	// Create a Config object
	config := models.Config{APIVersion: "1.0", DataVersion: "1.0"}

	// Return the Config object as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
