package models

import (
	"context"
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
	FindOne(query bson.M) (PersonInterface, error)
	FindMany(query bson.M, options options.FindOptions) ([]PersonShort, error)
	Insert(information bson.M, crumb *BreadCrumb) (*mongo.InsertOneResult, error)
	FindOneAndUpdate(query bson.M, update bson.M, crumb *BreadCrumb) (PersonInterface, error)
	Disconnect()
}
type PersonStore struct {
	config     *config.Config
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	cancel     context.CancelFunc
	store      PersonStoreInterface
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
func (store *PersonStore) Disconnect() {
	defer store.cancel()
	err := store.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

/**
* Insert a new person with the information provided
 */
func (store *PersonStore) Insert(information bson.M, crumb *BreadCrumb) (*mongo.InsertOneResult, error) {
	var result *mongo.InsertOneResult
	var err error

	// Add the breadcrumb
	information["lastSaved"] = crumb

	// Insert the document
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	result, err = store.collection.InsertOne(context, information)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
* Find a single person by _id
 */
func (store *PersonStore) FindOne(query bson.M) (PersonInterface, error) {
	var thePerson Person
	var err error
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = store.collection.FindOne(context, query).Decode(&thePerson)
	if err != nil {
		return nil, err
	}
	return &thePerson, nil
}

/**
* Find may people by a matcher
 */
func (store *PersonStore) FindMany(query bson.M, options options.FindOptions) ([]PersonShort, error) {
	var people []PersonShort
	var err error

	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	cursor, err := store.collection.Find(context, query, &options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context, &people)
	if err != nil {
		return nil, err
	}

	return people, nil
}

/**
* Find One person and Update with the data provided
 */
func (store *PersonStore) FindOneAndUpdate(query bson.M, updateValues bson.M, crumb *BreadCrumb) (PersonInterface, error) {
	var thePerson Person
	var err error

	// add breadcrumb to update object
	updateValues["lastSaved"] = crumb.AsBson()

	// Create the update object
	update := bson.M{"$set": updateValues}

	// Update the document
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	ctx, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = store.collection.FindOneAndUpdate(ctx, query, update, options).Decode(&thePerson)
	if err != nil {
		// throw the error up the call stack
		return nil, err
	}

	return &thePerson, nil
}

func (store *PersonStore) GetDatabaseVersion() (string, error) {
	var theVersion VersionInfo
	var err error

	query := bson.M{"name": "VERSION"}
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = store.collection.FindOne(context, query).Decode(&theVersion)
	if err != nil {
		// throw the error up the call stack
		return "", err
	}
	return theVersion.Version, nil
}
