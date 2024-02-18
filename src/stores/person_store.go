package stores

import (
	"encoding/json"

	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonStore struct {
	config     *config.Config
	MongoStore *MongoStore
}

const (
	CollectionName = "people"
)

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore(cfg *config.Config) *PersonStore {
	store := &PersonStore{}
	store.config = cfg
	store.MongoStore = NewMongoStore(cfg, "people", nil)
	return store
}

/**
* Insert a new person with the information provided
 */
func (store *PersonStore) Insert(information []byte, crumb *models.BreadCrumb) (*map[string]interface{}, error) {
	// Get the document values
	var insertValues bson.M
	err := json.Unmarshal(information, &insertValues)
	if err != nil {
		return nil, err
	}

	// Add the breadcrumb
	insertValues["lastSaved"] = crumb

	// Insert the document
	result, err := store.MongoStore.InsertOne(insertValues)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	// Get the new document
	return store.MongoStore.FindId(id)
}

/**
* Find One person and Update with the data provided
 */
func (store *PersonStore) FindOneAndUpdate(id string, request []byte, crumb *models.BreadCrumb) (*models.Person, error) {
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
	err = store.MongoStore.FindOneAndUpdate(query, update, options).Decode(&thePerson)
	if err != nil {
		// throw the error up the call stack
		return nil, err
	}

	return &thePerson, nil
}
