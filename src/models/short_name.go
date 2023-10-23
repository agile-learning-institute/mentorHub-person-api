package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortName struct {
	ID   primitive.ObjectID `bson:"ID,omitempty"`
	Name string             `json:"name,omitempty"`
}
