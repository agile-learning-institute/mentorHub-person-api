package models

import (
	"context"
	"log"

	"institute-person-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the PersonStoreInterface interface
type PersonStoreInterface interface {
	FindOne(query bson.M) PersonInterface
	FindMany(query bson.M, options options.FindOptions) []PersonShort
	Insert(information bson.M) *mongo.InsertOneResult
	FindOneAndUpdate(query bson.M, update bson.M) PersonInterface
}
type PersonStore struct {
	config     *config.Config
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	cancel     context.CancelFunc
	store      PersonStoreInterface
}

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore() PersonStore {
	this := PersonStore{}

	// get Configuration Values
	this.config = config.NewConfig()

	// Connect to the database
	ctx, cancel := this.config.GetTimeoutContext()
	this.cancel = cancel
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(this.config.GetConnectionString()))
	if err != nil {
		cancel()
		log.Fatal(err)
	}

	// Get the database and collection objects
	this.client = client
	this.database = this.client.Database(this.config.GetDatabaseName())
	this.collection = this.database.Collection(this.config.GetPeopleCollectionName())

	// Get the database Version
	this.config.DBVersion = this.GetDatabaseVersion()

	return this
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
func (store *PersonStore) Insert(information bson.M) *mongo.InsertOneResult {
	var result *mongo.InsertOneResult
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	result, _ = store.collection.InsertOne(context, information)
	return result
}

/**
* Find a single person by _id
 */
func (store *PersonStore) FindOne(query bson.M) PersonInterface {
	var thePerson PersonInterface
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	store.collection.FindOne(context, query).Decode(&thePerson)
	return thePerson
}

/**
* Find may people by a matcher
 */
func (store *PersonStore) FindMany(query bson.M, options options.FindOptions) []PersonShort {
	var people []PersonShort
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	cursor, _ := store.collection.Find(context, query, &options)
	cursor.All(context, &people)
	return people
}

/**
* Find One person and Update with the data provided
 */
func (store *PersonStore) FindOneAndUpdate(query bson.M, update bson.M) PersonInterface {
	var thePerson PersonInterface
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	ctx, cancel := store.config.GetTimeoutContext()
	defer cancel()
	store.collection.FindOneAndUpdate(ctx, query, update, options).Decode(&thePerson)
	return thePerson
}

func (store *PersonStore) GetDatabaseVersion() string {
	// TODO: - get the database schema version
	return "1.0.Dev"
}
