package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	store       *PersonStore
}

const (
	PeopleCollectionName = "people"
)

func NewPerson(theStore *PersonStore) *Person {
	this := &Person{}
	this.store = theStore
	return this
}

func (this *Person) GetPerson(id string) *Person {
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	return this.store.FindOne(query)
}

func (this *Person) PostPerson(body []byte) *Person {
	// Get the values to insert
	var insertValues bson.M
	json.Unmarshal(body, &insertValues)

	// Insert the new Person
	result := this.store.Insert(insertValues)
	query := bson.M{"_id": result.InsertedID}

	// Get the new document
	return this.store.FindOne(query)
}

func (this *Person) PatchPerson(id string, body []byte) *Person {
	// Build the query on ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	// Create the update set command
	var updateValues bson.M
	json.Unmarshal(body, &updateValues)
	update := bson.M{"$set": updateValues}

	// Update the document
	return this.store.FindOneAndUpdate(query, update)
}
