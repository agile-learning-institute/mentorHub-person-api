package models

import (
	"context"
	"encoding/json"
	"log"

	"institute-person-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VersionInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Version     string             `json:"version,omitempty"`
}

// Define the PersonStoreInterface interface
type PersonStoreInterface interface {
	FindOne(id string) (*Person, error)
	FindMany() ([]PersonShort, error)
	Insert(information []byte, crumb *BreadCrumb) (*Person, error)
	FindOneAndUpdate(id string, information []byte, crumb *BreadCrumb) (*Person, error)
	Disconnect()
}
type PersonStore struct {
	config     *config.Config
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	cancel     context.CancelFunc
}

const (
	PeopleCollectionName = "people"
	VersionDocumentName  = "VERSION"
)

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore(cfg *config.Config) (PersonStoreInterface, error) {
	this := &PersonStore{}
	var err error

	// get Configuration Values
	this.config = cfg

	// Connect to the database
	ctx, cancel := this.config.GetTimeoutContext()
	this.cancel = cancel
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(this.config.GetConnectionString()))
	if err != nil {
		cancel()
		return nil, err
	}

	// Get the database and collection objects
	this.client = client
	this.database = this.client.Database(this.config.GetDatabaseName())
	this.collection = this.database.Collection(this.config.GetPeopleCollectionName())

	// Put the database Version in the Config
	version, err := this.GetDatabaseVersion()
	if err != nil {
		cancel()
		return nil, err
	}
	this.config.SetDbVersion(version)

	return this, err
}

/**
* Disconnect from the database and housekeep
 */
func (this *PersonStore) Disconnect() {
	defer this.cancel()
	err := this.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

/**
* Insert a new person with the information provided
 */
func (this *PersonStore) Insert(information []byte, crumb *BreadCrumb) (*Person, error) {
	// Get the document values
	var insertValues bson.M
	err := json.Unmarshal(information, &insertValues)
	if err != nil {
		return nil, err
	}

	// Add the breadcrumb
	insertValues["lastSaved"] = crumb

	// Insert the document
	var result *mongo.InsertOneResult
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	result, err = this.collection.InsertOne(context, insertValues)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	// Get the new document
	person, err := this.FindOne(id)
	return person, err
}

/**
* Find a single person by _id
 */
func (this *PersonStore) FindOne(id string) (*Person, error) {
	// get the bson ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	var thePerson Person
	var err error
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	err = this.collection.FindOne(context, query).Decode(&thePerson)
	if err != nil {
		return nil, err
	}
	return &thePerson, nil
}

/**
* Find may people by a matcher
 */
func (this *PersonStore) FindMany() ([]PersonShort, error) {
	var people []PersonShort
	var err error

	// Setup the query
	query := bson.M{"name": bson.M{"$ne": "VERSION"}}
	options := options.Find().SetProjection(bson.D{{Key: "name", Value: 1}})

	// Query the database
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	cursor, err := this.collection.Find(context, query, options)
	if err != nil {
		return nil, err
	}

	// Fetch all the results
	err = cursor.All(context, &people)
	if err != nil {
		return nil, err
	}

	return people, nil
}

/**
* Find One person and Update with the data provided
 */
func (this *PersonStore) FindOneAndUpdate(id string, request []byte, crumb *BreadCrumb) (*Person, error) {
	var thePerson Person

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

	// Update the document
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	ctx, cancel := this.config.GetTimeoutContext()
	defer cancel()
	err = this.collection.FindOneAndUpdate(ctx, query, update, options).Decode(&thePerson)
	if err != nil {
		// throw the error up the call stack
		return nil, err
	}

	return &thePerson, nil
}

func (this *PersonStore) GetDatabaseVersion() (string, error) {
	var theVersion VersionInfo
	var err error

	query := bson.M{"name": "VERSION"}
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	err = this.collection.FindOne(context, query).Decode(&theVersion)
	if err != nil {
		// throw the error up the call stack
		return "", err
	}
	return theVersion.Version, nil
}
