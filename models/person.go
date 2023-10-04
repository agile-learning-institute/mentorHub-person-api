package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonInterface interface {
	GetStore() PersonStoreInterface
	GetPerson(id string) PersonInterface
	GetAllNames() []PersonShort
	PostPerson(body []byte) PersonInterface
	PatchPerson(id string, body []byte) PersonInterface
}

type PersonShort struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `json:"name,omitempty"`
}
type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	store       PersonStoreInterface
	person      PersonInterface
}

const (
	PeopleCollectionName = "people"
)

func NewPerson(theStore PersonStoreInterface) PersonInterface {
	this := &Person{}
	this.store = theStore
	return this
}

func (this *Person) GetStore() PersonStoreInterface {
	return this.store
}
func (this *Person) GetPerson(id string) PersonInterface {
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	return this.store.FindOne(query)
}
func (this *Person) GetAllNames() []PersonShort {
	query := bson.M{}
	theOptions := options.Find().SetProjection(bson.D{{Key: "name", Value: 1}})
	return this.store.FindMany(query, *theOptions)
}

func (this *Person) PostPerson(body []byte) PersonInterface {
	// Get the values to insert
	var insertValues bson.M
	json.Unmarshal(body, &insertValues)

	// Insert the new Person
	result := this.store.Insert(insertValues)
	query := bson.M{"_id": result.InsertedID}

	// Get the new document
	return this.store.FindOne(query)
}

func (this *Person) PatchPerson(id string, body []byte) PersonInterface {
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
