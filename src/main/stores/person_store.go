package stores

import (
	"encoding/json"

	"institute-person-api/src/main/config"
	"institute-person-api/src/main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonStore struct {
	config *config.Config
	Store  *MongoStore
}

const (
	CollectionName = "people"
)

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore(cfg *config.Config) *PersonStore {
	this := &PersonStore{}
	this.config = cfg
	this.Store = NewMongoStore(cfg, CollectionName, MongoQueryNotVersion(), MongoShortNameProjection())
	return this
}

/**
* Insert a new person with the information provided
 */
func (this *PersonStore) Insert(information []byte, crumb *models.BreadCrumb) (*interface{}, error) {
	// Get the document values
	var insertValues bson.M
	err := json.Unmarshal(information, &insertValues)
	if err != nil {
		return nil, err
	}

	// Add the breadcrumb
	insertValues["lastSaved"] = crumb

	// Insert the document
	result, err := this.Store.InsertOne(insertValues)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	// Get the new document
	return this.Store.FindId(id)
}

/**
* Find One person and Update with the data provided
 */
func (this *PersonStore) FindOneAndUpdate(id string, request []byte, crumb *models.BreadCrumb) (*models.Person, error) {
	var thePerson models.Person

	// Build the query on ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	// Create the update set values
	var updateValues bson.M
	err := json.Unmarshal(request, &updateValues)
	if err != nil {
		return nil, err
	}

	// add breadcrumb to update object
	updateValues["lastSaved"] = crumb.AsBson()

	// Create the update object
	update := bson.M{"$set": updateValues}

	// Set Options
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// Update the document
	err = this.Store.FindOneAndUpdate(query, update, options).Decode(&thePerson)
	if err != nil {
		// throw the error up the call stack
		return nil, err
	}

	return &thePerson, nil
}
