package models

type Person struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Config struct {
	APIVersion  string `json:"apiVersion"`
	DataVersion string `json:"dataVersion"`
}
