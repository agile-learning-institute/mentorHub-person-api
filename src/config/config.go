package config

import (
	"context"
	"log"
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
	Filter         bson.M `json:"filter"`
}

type Config struct {
	ConfigItems      []*ConfigItem
	Stores           []*StoreItem
	ApiVersion       string
	port             string
	patch            string
	configFolder     string
	databaseName     string
	databaseTimeout  int
	connectionString string
	client           *mongo.Client
	database         *mongo.Database
	cancel           context.CancelFunc
}

const (
	VersionMajor            = "1"
	VersionMinor            = "3"
	DefaultConfigFolder     = "./"
	DefaultConnectionString = "mongodb://root:example@localhost:27017/?tls=false&directConnection=true"
	DefaultDatabaseName     = "agile-learning-institute"
	DefaultPort             = ":8080"
	DefaultTimeout          = 10
)

/**
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

/**
* Disconnect fromthe Database
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

	// Get the database and collection objects
	cfg.client = client
	cfg.database = cfg.client.Database(cfg.databaseName)
}

/**
* Disconnect fromthe Database
 */
func (cfg *Config) Disconnect() {
	ctx, cancel := cfg.GetTimeoutContext()
	defer cancel()
	cfg.client.Disconnect(ctx)
	cfg.cancel()
}

/**
* Get the port config value
 */
func (cfg *Config) GetPort() string {
	return cfg.port
}

/**
* Get mongo Collection
 */
func (cfg *Config) GetCollection(name string) mongo.Collection {
	return *cfg.database.Collection(name)
}

/**
* Register a Config Store
 */
func (cfg *Config) AddConfigStore(theStore *StoreItem) {
	cfg.Stores = append(cfg.Stores, theStore)
}

/**
* Get a Timeout Context using the configured defalut wait
 */
func (cfg *Config) GetTimeoutContext() (context.Context, context.CancelFunc) {
	timeout := time.Duration(cfg.databaseTimeout) * time.Second
	return context.WithTimeout(context.Background(), timeout)
}

/**
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
