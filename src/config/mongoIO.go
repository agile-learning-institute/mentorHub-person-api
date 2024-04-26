/********************************************************************************
** MongoIO
**    This class is a wrapper for MongoDB interactions. It manages
**    the database connection, and executes all mongoDb calls with
**    easy to mock wrappers. This requires a local mongoDb with test data
**	  You can start this contianer with the command "mh up mongodb"
********************************************************************************/
package config

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortName struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	Name string             `json:"name,omitempty"`
}

type Enums struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Version     int                    `json:"version"`
	Enumerators map[string]interface{} `json:"enumerators"`
}

type MongoIO struct {
	config             *Config
	client             *mongo.Client
	database           *mongo.Database
	peopleCollection   *mongo.Collection
	enumsCollection    *mongo.Collection
	partnersCollection *mongo.Collection
	versionsCollection *mongo.Collection
	cancel             context.CancelFunc
}

type MongoIOInterface interface {
	Connect()
	Disconnect()
	Find(collection *mongo.Collection, query bson.M, opts *options.FindOptions, results interface{}) error
	FindOne(collection *mongo.Collection, query bson.M, results interface{}) error
	InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(collection *mongo.Collection, query bson.M, opts *options.FindOneAndUpdateOptions, update interface{}, results interface{}) error
}

var _ MongoIOInterface = (*MongoIO)(nil)

const (
	PeoplesCollectionName     = "people"
	PartnersCollectionName    = "partners"
	EnumeratorsCollectionName = "enumerators"
	VersionsCollectionName    = "msmCurrentVersions"
)

func NameProjection() bson.D {
	return bson.D{
		{Key: "ID", Value: "_id"},
		{Key: "name", Value: bson.M{"$concat": bson.A{"$firstName", " ", "$lastName"}}},
	}
}

/**********************************************************************
* Constructor - initilize configuration values
 */
func NewMongoIO(config *Config) *MongoIO {
	this := &MongoIO{}
	this.config = config
	return this
}

/********************************************************************************
* Connect to the Database, initilize connected values
 */
func (mongoIO *MongoIO) Connect() {
	// Connect to the database
	ctx, cancel := mongoIO.getTimeoutContext()
	mongoIO.cancel = cancel
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoIO.config.GetConnectionString()))
	if err != nil {
		cancel()
		log.Fatal("Database Connection Failed:", err)
	}

	// Get the database object, and collection objects
	mongoIO.client = client
	mongoIO.database = mongoIO.client.Database(mongoIO.config.GetDatabaseName())
	mongoIO.peopleCollection = mongoIO.database.Collection(PeoplesCollectionName)
	mongoIO.enumsCollection = mongoIO.database.Collection(EnumeratorsCollectionName)
	mongoIO.partnersCollection = mongoIO.database.Collection(PartnersCollectionName)
	mongoIO.versionsCollection = mongoIO.database.Collection(VersionsCollectionName)
}

/********************************************************************************
* Disconnect fromthe Database
 */
func (mongoIO *MongoIO) Disconnect() {
	ctx, cancel := mongoIO.getTimeoutContext()
	defer cancel()
	mongoIO.client.Disconnect(ctx)
	mongoIO.cancel()
	mongoIO.database = nil
	mongoIO.peopleCollection = nil
	mongoIO.enumsCollection = nil
	mongoIO.partnersCollection = nil
	mongoIO.versionsCollection = nil
	cancel()
}

/**
* Find multiple documents
 */
func (mongoIO *MongoIO) Find(collection *mongo.Collection, query bson.M, opts *options.FindOptions, results interface{}) error {
	context, cancel := mongoIO.getTimeoutContext()
	defer cancel()

	cursor, err := collection.Find(context, query, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(context)

	if err = cursor.All(context, results); err != nil {
		return err
	}
	return nil
}

/**
* Find one document
 */
func (mongoIO *MongoIO) FindOne(collection *mongo.Collection, query bson.M, results interface{}) error {
	context, cancel := mongoIO.getTimeoutContext()
	defer cancel()
	isok := collection.FindOne(context, query).Decode(results)
	return isok
}

/**
* Insert one document
 */
func (mongoIO *MongoIO) InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	context, cancel := mongoIO.getTimeoutContext()
	defer cancel()
	log.Print(document)
	result, err := collection.InsertOne(context, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
* Update one document
 */
func (mongoIO *MongoIO) UpdateOne(collection *mongo.Collection, query bson.M, opts *options.FindOneAndUpdateOptions, update interface{}, results interface{}) error {
	ctx, cancel := mongoIO.getTimeoutContext()
	defer cancel()
	result := collection.FindOneAndUpdate(ctx, query, update, opts)
	if result.Err() != nil {
		return result.Err()
	}
	err := result.Decode(results)
	if err != nil {
		return err
	}

	return nil
}

/**
* Simple getter for peopleCollection
 */
func (mongoIO *MongoIO) GetPeopleCollection() *mongo.Collection {
	return mongoIO.peopleCollection
}

/**
* Get a timeout to be used with mongo calls
 */
func (mongoIO *MongoIO) getTimeoutContext() (context.Context, context.CancelFunc) {
	timeout := time.Duration(mongoIO.config.GetDatabaseTimeout()) * time.Second
	return context.WithTimeout(context.Background(), timeout)
}

/**
*
 */
func (mongoIO *MongoIO) LoadVersions() {
	// Find Versions
	query := bson.M{}
	opts := options.Find()
	err := mongoIO.Find(mongoIO.versionsCollection, query, opts, &mongoIO.config.Versions)
	if err != nil {
		log.Fatal("Error geting Versions")
	}

	if len(mongoIO.config.Versions) == 0 {
		log.Fatal("No Versions not found")
	}
}

/**
* Load the Enumerators that match the person collection
 */
func (mongoIO *MongoIO) LoadEnumerators() {
	// Find the current version for the people collection
	versionString := ""
	for _, v := range mongoIO.config.Versions {
		if v.CollectionName == PeoplesCollectionName {
			versionString = v.CurrentVersion
		}
	}
	if versionString == "" {
		log.Fatalf("People Collection not found %d", len(mongoIO.config.Versions))
	}

	// Extract the enumerators version number from the version string
	parts := strings.Split(versionString, ".")
	versionNumber, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		log.Fatalf("Version Number not an integer: %s", versionString)
	}

	// Query Enumerators
	results := []*Enums{}
	query := bson.M{"version": versionNumber}
	opts := options.Find()
	err = mongoIO.Find(mongoIO.enumsCollection, query, opts, &results)
	if err != nil {
		log.Fatalf("Error getting Enumerators %s", err)
	}

	// Set enumerators
	if len(results) == 1 {
		mongoIO.config.Enumerators = results[0].Enumerators
	} else {
		log.Fatalf("Enumerators List is wrong size %d", len(results))
	}
}

/********************************************************************************
*
 */
func (mongoIO *MongoIO) FetchMentors() error {

	// Fetch Mentors
	results := []*ShortName{}
	query := bson.M{"$and": []bson.M{
		{"status": bson.M{"$ne": "Archived"}},
		{"roles": "Mentor"},
	}}
	opts := options.Find()
	opts.SetProjection(NameProjection())
	err := mongoIO.Find(mongoIO.peopleCollection, query, opts, &results)
	if err != nil {
		log.Fatalf("Error getting Mentors %s", err)
	}
	mongoIO.config.Mentors = results
	return nil
}

/********************************************************************************
*
 */
func (mongoIO *MongoIO) FetchPartners() error {
	// Fetch Partners
	results := []*ShortName{}
	query := bson.M{"status": bson.M{"$ne": "Archived"}}
	opts := options.Find()
	err := mongoIO.Find(mongoIO.partnersCollection, query, opts, &results)
	if err != nil {
		log.Fatalf("Error getting Partners %s", err)
	}
	mongoIO.config.Partners = results

	return nil
}
