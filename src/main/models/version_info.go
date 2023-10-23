package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type VersionInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Version     string             `json:"version,omitempty"`
}
