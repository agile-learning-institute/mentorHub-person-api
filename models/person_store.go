package models

import (
	"context"
	"log"

	"institute-person-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonStore struct {
	config     *config.Config
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	cancel     context.CancelFunc
}

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore() *PersonStore {
	this := &PersonStore{}

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
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	result, _ := store.collection.InsertOne(context, information)
	return result
}

/**
* Find a single person by _id
 */
func (store *PersonStore) FindOne(query bson.M) *Person {
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	var thePerson Person
	store.collection.FindOne(context, query).Decode(&thePerson)
	return &thePerson
}

/**
* Find may people by a matcher
 */
func (store *PersonStore) FindMany(query bson.M) *[]Person {
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	cursor, _ := store.collection.Find(context, query)
	var people []Person
	cursor.All(context, &people)
	return &people
}

/**
* Find One person and Update with the data provided
 */
func (store *PersonStore) FindOneAndUpdate(query bson.M, update bson.M) *Person {
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	ctx, cancel := store.config.GetTimeoutContext()
	defer cancel()
	var thePerson Person
	store.collection.FindOneAndUpdate(ctx, query, update, options).Decode(&thePerson)
	return &thePerson
}

func (store *PersonStore) GetDatabaseVersion() string {
	// TODO: - get the database schema version
	return ""
}
