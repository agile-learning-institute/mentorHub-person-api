package models

type Enumerator struct {
	Name         string   `json:"name,omitempty"`
	Enumerations []string `json:"enumerations,omitempty"`
}
