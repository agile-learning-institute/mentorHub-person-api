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
	"strings"
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

type CurrentVersion struct {
	CollectionName string `json:"collectionName"`
	CurrentVersion string `json:"currentVersion"`
}

type Enumerators struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Version     int                    `json:"version"`
	Enumerators map[string]interface{} `json:"enumerators"`
}

type Config struct {
	// Configuration values exposed on the /config endpoint
	ConfigItems []*ConfigItem       `json:"ConfigItems"`
	Versions    []*CurrentVersion   `json:"versions"`
	Enumerators interface{}         `json:"enums"`
	Mentors     []*models.ShortName `json:"mentors"`
	Partners    []*models.ShortName `json:"partners"`
	ApiVersion  string

	// Configurations returned by getters
	port             string
	patch            string
	configFolder     string
	databaseName     string
	databaseTimeout  int
	connectionString string

	// MongoDB Variables
	client             *mongo.Client
	database           *mongo.Database
	peopleCollection   *mongo.Collection
	enumsCollection    *mongo.Collection
	partnersCollection *mongo.Collection
	versionsCollection *mongo.Collection
	cancel             context.CancelFunc
}

const (
	// Version Major.Minog PATCH loaded at runtime
	VersionMajor = "1"
	VersionMinor = "3"

	// Default Values for Configuration Values
	DefaultConfigFolder     = "./"
	DefaultConnectionString = "mongodb://root:example@localhost:27017/?tls=false&directConnection=true"
	DefaultDatabaseName     = "mentorHub"
	DefaultPort             = ":8082"
	DefaultTimeout          = "10"

	// Collection Names
	PeoplesCollectionName     = "people"
	PartnersCollectionName    = "partners"
	EnumeratorsCollectionName = "enumerators"
	VersionsCollectionName    = "msmCurrentVersions"
)

// Common MongoDB Query Objects
func GetAllQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"name": bson.M{"$ne": "VERSION"}},
		{"status": bson.M{"$ne": "Archived"}},
	}}
}

func GetMentorsQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"status": bson.M{"$ne": "Archived"}},
		{"roles": "Mentor"},
	}}
}

func GetMentorsProjection() bson.D {
	return bson.D{
		{Key: "ID", Value: "_id"},
		{Key: "name", Value: bson.M{"$concat": bson.A{"$firstName", " ", "$lastName"}}},
	}
}

/**********************************************************************
* Constructor - initilize configuration values
 */
func NewConfig() *Config {
	this := &Config{}

	/**********************************************************************
	* Find a value in the following order
	*   If a file value exists, use that
	*	Else if an environment variable exists, use that
	* 	Else use the default value provided
	 */
	findValue := func(key string, defaultValue string, secret bool) string {
		var theValue string
		var from string

		// Start with default values
		theValue = defaultValue
		from = "default"

		// Check for Environemt Variable
		envValue, isSet := os.LookupEnv(key)
		if isSet {
			theValue = envValue
			from = "environment"
		}

		// Check for a file value
		var theFile = this.configFolder + key
		_, error := os.Stat(theFile)
		if error == nil {
			fileContent, err := os.ReadFile(theFile)
			if err == nil {
				theValue = string(fileContent)
				from = "file"
			}
		}

		// Create the ConfigItem and add it to the list
		theItem := &ConfigItem{Name: key, From: from}
		if secret {
			theItem.Value = "Secret"
		} else {
			theItem.Value = theValue
		}
		this.ConfigItems = append(this.ConfigItems, theItem)

		// Return the config value
		return theValue
	}

	// Load Confiuration Values
	var err error
	this.patch = findValue("PATCH_LEVEL", "LocalDev", false)
	this.configFolder = findValue("CONFIG_FOLDER", DefaultConfigFolder, false)
	this.connectionString = findValue("CONNECTION_STRING", DefaultConnectionString, true)
	this.databaseName = findValue("DATABASE_NAME", DefaultDatabaseName, false)
	this.port = findValue("PORT", DefaultPort, false)
	this.ApiVersion = VersionMajor + "." + VersionMinor + "." + this.patch
	this.databaseTimeout, err = strconv.Atoi(findValue("CONNECTION_TIMEOUT", DefaultTimeout, false))
	if err != nil {
		log.Fatal("Invalid CONNECTION_TIMEOUT value", err)
	}

	return this
}

/********************************************************************************
* Connect to the Database, load the Versions and Enumerators Config values
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

	// Get the database object, and collection objects
	cfg.client = client
	cfg.database = cfg.client.Database(cfg.databaseName)
	cfg.peopleCollection = cfg.database.Collection(PeoplesCollectionName)
	cfg.enumsCollection = cfg.database.Collection(EnumeratorsCollectionName)
	cfg.partnersCollection = cfg.database.Collection(PartnersCollectionName)
	cfg.versionsCollection = cfg.database.Collection(VersionsCollectionName)

	// Find Versions
	query := bson.M{}
	opts := options.Find()
	cfg.FindDocuments(cfg.versionsCollection, query, opts, &cfg.Versions)
	if len(cfg.Versions) == 0 {
		log.Fatal("Versions not found")
	}

	// Find the current version for the people collection
	versionString := ""
	for _, v := range cfg.Versions {
		if v.CollectionName == PeoplesCollectionName {
			versionString = v.CurrentVersion
		}
	}
	if versionString == "" {
		log.Fatalf("People Collection not found %d", len(cfg.Versions))
	}

	// Extract the enumerators version number from the version string
	var versionNumber int
	parts := strings.Split(versionString, ".")
	versionNumber, err = strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		log.Fatalf("Version Number not an integer: %s", versionString)
	}

	// Query Enumerators
	results := []*Enumerators{}
	query = bson.M{"version": versionNumber}
	opts = options.Find()
	opts.SetSort(bson.D{{Key: "version", Value: 1}})
	cfg.FindDocuments(cfg.enumsCollection, query, opts, &results)

	// Set enumerators
	if len(results) == 1 {
		cfg.Enumerators = results[0].Enumerators
	} else {
		log.Fatalf("Enumerators List is wrong size %d", len(results))
	}
}

func (cfg *Config) FindDocuments(collection *mongo.Collection, query bson.M, opts *options.FindOptions, results interface{}) {
	context, cancel := cfg.GetTimeoutContext()
	defer cancel()

	cursor, err := collection.Find(context, query, opts)
	if err != nil {
		log.Fatalf("Query Failed: %s %v", collection.Name(), err)
	}
	defer cursor.Close(context)

	// Directly pass results to cursor.All
	if err = cursor.All(context, results); err != nil {
		log.Fatalf("Fetch Failed: %s %v", collection.Name(), err)
	}
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
