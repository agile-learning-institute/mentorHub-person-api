package models

import (
	"context"
	"log"

	"institute-person-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EnumeratorStore struct {
	config      *config.Config
	client      *mongo.Client
	database    *mongo.Database
	collection  *mongo.Collection
	cancel      context.CancelFunc
	Enumerators []*Enumerator
}

/**
* Construct a PersonStore to handle person database io
 */
func NewEnumeratorStore(cfg *config.Config) (*EnumeratorStore, error) {
	this := &EnumeratorStore{}
	var err error

	// get Configuration Values
	this.config = cfg

	// Connect to the database
	ctx, cancelConnect := this.config.GetTimeoutContext()
	defer cancelConnect()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(this.config.GetConnectionString()))
	if err != nil {
		cancelConnect()
		return nil, err
	}

	// Get the database and collection objects
	this.client = client
	this.database = this.client.Database(this.config.GetDatabaseName())
	this.collection = this.database.Collection(this.config.GetEnumeratorCollectionName())
	log.Println("Constructor Collection Name:", this.collection.Name())
	log.Println("Constructor Name:", this.config.GetEnumeratorCollectionName())
	// Put the database Version in the Config
	this.ReportVersion()

	// Load all the Enumerators
	err = this.LoadEnumerators()
	if err != nil {
		log.Fatal(err)
	}

	// Disconnect
	err = this.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return this, nil
}

func (this *EnumeratorStore) GetDatabaseVersion() (string, error) {
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

/**
* Find may people by a matcher
 */
func (this *EnumeratorStore) LoadEnumerators() error {
	var err error

	// Setup the query
	query := bson.M{"name": bson.M{"$ne": "VERSION"}}

	// Query the database
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	cursor, err := this.collection.Find(context, query)
	if err != nil {
		log.Println("Collection Name:", this.collection.Name())
		log.Println("Query Error:", err.Error())
		return err
	}

	// Fetch all the results
	err = cursor.All(context, &this.Enumerators)
	if err != nil {
		log.Println("Fetch Error:", err.Error())
		return err
	}
	return nil
}

func (this *EnumeratorStore) ReportVersion() {
	var theVersion VersionInfo
	var err error

	query := bson.M{"name": "VERSION"}
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	err = this.collection.FindOne(context, query).Decode(&theVersion)
	if err != nil {
		this.config.SetEnumVersion(err.Error())
	}
	this.config.SetEnumVersion(theVersion.Version)
}
