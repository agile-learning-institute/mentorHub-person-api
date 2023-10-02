package models

import (
	"encoding/json"
	"institute-person-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
}

const (
	PeopleCollectionName = "people"
)

func GetPerson(id string, config *config.Config) *Person {
	var thePerson Person
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	people := config.GetCollection(PeopleCollectionName)
	ctx, cancel := config.GetTimeoutContext()
	defer cancel()
	people.FindOne(ctx, query).Decode(&thePerson)
	return &thePerson
}

func PostPerson(body []byte, config *config.Config) *Person {
	// Get the values to insert
	var insertValues bson.M
	json.Unmarshal(body, &insertValues)

	// Insert the new Person
	people := config.GetCollection(PeopleCollectionName)
	addContext, addCancel := config.GetTimeoutContext()
	defer addCancel()
	result, _ := people.InsertOne(addContext, insertValues)

	// Get the new document
	query := bson.M{"_id": result.InsertedID}
	getContext, getCancel := config.GetTimeoutContext()
	defer getCancel()
	var thePerson Person
	people.FindOne(getContext, query).Decode(&thePerson)

	return &thePerson
}

func PatchPerson(id string, body []byte, config *config.Config) *Person {
	// Get the update values
	var updateValues bson.M
	json.Unmarshal(body, &updateValues)

	// Build the update parameteres
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": updateValues}

	// Update the document
	var updatedPerson Person
	people := config.GetCollection(PeopleCollectionName)
	ctx, cancel := config.GetTimeoutContext()
	defer cancel()
	people.FindOneAndUpdate(ctx, query, update, options).Decode(&updatedPerson)

	return &updatedPerson
}
