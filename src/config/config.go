/********************************************************************************
** Config
**    This class Manages the mongo database connection for the API
** 		Configuration values are managed by the config which implements a
**      environment, file, default heiaracy for configurable values.
**    This calss also supports a convenience /config endpoint that the person UI
**    uses to get configuraiton information, enumerators, and needed select values.
********************************************************************************/
package config

import (
	"context"
	"log"
	"mentorhub-person-api/src/models"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	From  string `json:"from"`
}

type StoreItem struct {
	CollectionName string `json:"collectionName"`
	Version        string `json:"version"`
}

type Config struct {
	ConfigItems        []*ConfigItem
	Stores             []*StoreItem
	Mentors            []*models.ShortName    `json:"mentors"`
	Partners           []*models.ShortName    `json:"partners"`
	Enumerators        map[string]interface{} `json:"enums"`
	ApiVersion         string
	port               string
	patch              string
	configFolder       string
	databaseName       string
	databaseTimeout    int
	connectionString   string
	client             *mongo.Client
	database           *mongo.Database
	peopleCollection   *mongo.Collection
	enumsCollection    *mongo.Collection
	partnersCollection *mongo.Collection
	cancel             context.CancelFunc
}

const (
	VersionMajor            = "1"
	VersionMinor            = "3"
	DefaultConfigFolder     = "./"
	DefaultConnectionString = "mongodb://root:example@localhost:27017/?tls=false&directConnection=true"
	DefaultDatabaseName     = "agile-learning-institute"
	DefaultPort             = ":8082"
	DefaultTimeout          = 10

	PeoplesCollectionName     = "people"
	PartnersCollectionName    = "partners"
	EnumeratorsCollectionName = "enumerators"
)

func GetAllQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"name": bson.M{"$ne": "VERSION"}},
		{"status": bson.M{"$ne": "Archived"}},
	}}
}

func GetMentorsQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"status": bson.M{"$ne": "Archived"}},
		{"mentor": true},
	}}
}

func GetMentorsProjection() bson.D {
	return bson.D{
		{Key: "ID", Value: "_id"},
		{Key: "name", Value: bson.M{"$concat": bson.A{"$firstName", " ", "$lastName"}}},
	}
}

/**********************************************************************
* Construct a config item by finding all the configuration values
 */
func NewConfig() *Config {
	this := &Config{}

	// Load Confiuration Values
	this.patch = this.findStringValue("PATCH_LEVEL", "LocalDev", false)
	this.configFolder = this.findStringValue("CONFIG_FOLDER", DefaultConfigFolder, false)
	this.connectionString = this.findStringValue("CONNECTION_STRING", DefaultConnectionString, true)
	this.databaseName = this.findStringValue("DATABASE_NAME", DefaultDatabaseName, false)
	this.databaseTimeout = this.findIntValue("CONNECTION_TIMEOUT", DefaultTimeout, false)
	this.port = this.findStringValue("PORT", DefaultPort, false)
	this.ApiVersion = VersionMajor + "." + VersionMinor + "." + this.patch

	return this
}

/********************************************************************************
* Find a configuration value, and build the ConfigItems array
* 	If in an Environment Variable exists use it
* 	If not, and a Config File exists use that
* 	If all else fails use the default value provided
 */
func (cfg *Config) findStringValue(key string, defaultValue string, secret bool) string {
	var theValue string
	var from string

	// Start with default values
	theValue = defaultValue
	from = "default"

	// Check for a file value
	fileValue := cfg.fileValue(key)
	if fileValue != "" {
		theValue = fileValue
		from = "file"
	}

	// Check for Environemt Variable
	envValue, isSet := os.LookupEnv(key)
	if isSet {
		theValue = envValue
		from = "environment"
	}

	// Create the CI and add it to the list
	theItem := &ConfigItem{Name: key, From: from}
	if secret {
		theItem.Value = "Secret"
	} else {
		theItem.Value = theValue
	}
	cfg.ConfigItems = append(cfg.ConfigItems, theItem)

	// Return the config value
	return theValue
}

/**
* Find an Integer configuration value - find the string value and convert it
 */
func (cfg *Config) findIntValue(key string, defaultValue int, secret bool) int {
	theValue := cfg.findStringValue(key, strconv.Itoa(defaultValue), secret)
	theInteger, _ := strconv.Atoi(theValue)
	return theInteger
}

/**
* Return the contents of a file if it exists, or an empty string otherwise
 */
func (cfg *Config) fileValue(key string) string {
	// Check for Config in a File
	var theFile = cfg.configFolder + key
	_, error := os.Stat(theFile)
	if error == nil {
		fileContent, err := os.ReadFile(theFile)
		if err == nil {
			return string(fileContent)
		}
	}
	return ""
}

/********************************************************************************
* Connect to the Database
 */
func (cfg *Config) Connect() {
	// Connect to the database
	ctx, cancel := cfg.GetTimeoutContext()
	cfg.cancel = cancel
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.connectionString))
	if err != nil {
		cancel()
		log.Fatal("Database Connection Failed:", err)
	}

	// Get the database object
	cfg.client = client
	cfg.database = cfg.client.Database(cfg.databaseName)

	// Initilize Collections
	cfg.peopleCollection = cfg.registerCollection(PeoplesCollectionName)
	cfg.partnersCollection = cfg.registerCollection(PartnersCollectionName)
	cfg.enumsCollection = cfg.registerCollection(EnumeratorsCollectionName)

	// Load Enumerators
	cfg.getEnumerators()
}

/**
* Register a Collection Store
 */
func (cfg *Config) registerCollection(collectionName string) *mongo.Collection {
	collection := cfg.database.Collection(collectionName)
	version := cfg.GetVersion(collection)
	var storeItem = StoreItem{
		CollectionName: collectionName,
		Version:        version,
	}
	cfg.Stores = append(cfg.Stores, &storeItem)

	return collection
}

/**
* Get the collection schema version
 */
func (cfg *Config) GetVersion(collection *mongo.Collection) string {
	var theVersion models.VersionInfo
	var err error

	query := bson.M{"name": "VERSION"}
	context, cancel := cfg.GetTimeoutContext()
	defer cancel()
	err = collection.FindOne(context, query).Decode(&theVersion)
	if err != nil {
		return err.Error()
	}
	return theVersion.Version
}

/**
* Load Enumerators
 */
func (cfg *Config) getEnumerators() {
	// Query Enumerators
	opts := options.Find()
	context, cancel := cfg.GetTimeoutContext()
	defer cancel()
	cursor, err := cfg.enumsCollection.Find(context, GetAllQuery(), opts)
	if err != nil {
		cancel()
		log.Fatal("Query Enumerators Failed:", err)
	}

	// Fetch Enumerators
	var result []map[string]interface{}
	err = cursor.All(context, &result)
	if err != nil {
		cancel()
		log.Fatal("Fetch Enumerators Failed:", err)
	}
	cfg.Enumerators = result[0]
}

/********************************************************************************
* Disconnect fromthe Database
 */
func (cfg *Config) Disconnect() {
	ctx, cancel := cfg.GetTimeoutContext()
	defer cancel()
	cfg.client.Disconnect(ctx)
	cfg.cancel()
}

/********************************************************************************
* Simple Getters
 */
func (cfg *Config) GetPort() string {
	return cfg.port
}

func (cfg *Config) GetPersonCollection() *mongo.Collection {
	return cfg.peopleCollection
}

func (cfg *Config) GetTimeoutContext() (context.Context, context.CancelFunc) {
	timeout := time.Duration(cfg.databaseTimeout) * time.Second
	return context.WithTimeout(context.Background(), timeout)
}

/********************************************************************************
* Load the Mentors and Partners list
 */
func (cfg *Config) LoadLists() error {

	// Fetch Mentors
	mentors, err := cfg.findNames(
		cfg.peopleCollection,
		options.Find().SetProjection(GetMentorsProjection()),
		GetMentorsQuery(),
	)

	if err != nil {
		log.Printf("ERROR: Load Mentors failed %s", err)
		return err
	}
	cfg.Mentors = mentors

	// Fetch Partners
	partners, err := cfg.findNames(
		cfg.partnersCollection,
		options.Find(),
		GetAllQuery(),
	)

	if err != nil {
		log.Printf("ERROR: Load Partners failed %s", err)
		return err
	}
	cfg.Partners = partners

	return nil
}

/**
* Get a list of people names based on the query and options provided
 */
func (cfg *Config) findNames(collection *mongo.Collection, opts *options.FindOptions, query bson.M) ([]*models.ShortName, error) {
	var results []*models.ShortName
	var err error

	// Query the database
	context, cancel := cfg.GetTimeoutContext()
	defer cancel()
	cursor, err := collection.Find(context, query, opts)
	if err != nil {
		return nil, err
	}

	// Fetch all the results
	context, cancel = cfg.GetTimeoutContext()
	defer cancel()
	err = cursor.All(context, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
