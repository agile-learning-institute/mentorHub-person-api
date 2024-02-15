package models

type Enumerator struct {
	Name         string `json:"name,omitempty"`
	Enumerations any    `json:"enumerations,omitempty"`
}
