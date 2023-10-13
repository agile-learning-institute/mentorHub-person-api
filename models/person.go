package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonInterface interface {
	GetPerson(id string) (PersonInterface, error)
	GetAllNames() ([]PersonShort, error)
	PostPerson(body []byte) (PersonInterface, error)
	PatchPerson(id string, body []byte) (PersonInterface, error)
}

type PersonShort struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `json:"name,omitempty"`
}
type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Store       PersonStoreInterface
}

func NewPerson(theStore PersonStoreInterface) PersonInterface {
	this := &Person{}
	this.Store = theStore
	return this
}

func (this *Person) GetPerson(id string) (PersonInterface, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	result, err := this.Store.FindOne(query)
	return result, err
}

func (this *Person) GetAllNames() ([]PersonShort, error) {
	query := bson.M{}
	theOptions := options.Find().SetProjection(bson.D{{Key: "name", Value: 1}})
	result, err := this.Store.FindMany(query, *theOptions)
	return result, err
}

func (this *Person) PostPerson(body []byte) (PersonInterface, error) {
	// Get the values to insert
	var insertValues bson.M
	err := json.Unmarshal(body, &insertValues)
	if err != nil {
		return nil, err
	}

	// Insert the new Person
	result, err := this.Store.Insert(insertValues)
	if err != nil {
		return nil, err
	}

	// Get the new document
	query := bson.M{"_id": result.InsertedID}
	person, err := this.Store.FindOne(query)
	return person, err
}

func (this *Person) PatchPerson(id string, body []byte) (PersonInterface, error) {
	// Build the query on ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	// Create the update set command
	var updateValues bson.M
	err := json.Unmarshal(body, &updateValues)
	if err != nil {
		return nil, err
	}
	update := bson.M{"$set": updateValues}

	// Update the document
	return this.Store.FindOneAndUpdate(query, update)
}
