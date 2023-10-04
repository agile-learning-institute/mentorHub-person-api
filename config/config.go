package config

import (
	"context"
	"strconv"
	"time"
)

type ConfigItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	From  string `json:"from"`
}

type Config struct {
	ConfigItems          []*ConfigItem
	Version              string
	DBVersion            string
	patch                string
	configFolder         string
	databaseName         string
	peopleCollectionName string
	databaseTimeout      int
	connectionString     string
}

const (
	VersionMajor                = "1"
	VersionMinor                = "0"
	DefaultConfigFolder         = "/opt/"
	DefaultConnectionString     = "mongodb://root:example@institute-person-db:27017/?tls=false&directConnection=true"
	DefaultDatabaseName         = "agile-learning-institute"
	DefaultPeopleCollectionName = "people"
	DefaultTimeout              = 10
)

/**
* Construct a config item by finding all the configuration values
 */
func NewConfig() *Config {
	this := &Config{}
	this.configFolder = this.findStringValue("CONFIG_FOLDER", DefaultConfigFolder, false)
	this.connectionString = this.findStringValue("CONNECTION_STRING", DefaultConnectionString, true)
	this.databaseName = this.findStringValue("DATABASE_NAME", DefaultDatabaseName, false)
	this.peopleCollectionName = this.findStringValue("PEOPLE_COLLECTION_NAME", DefaultPeopleCollectionName, false)
	this.databaseTimeout = this.findIntValue("CONNECTION_TIMEOUT", DefaultTimeout, false)
	this.patch = this.findStringValue("PATCH_LEVEL", "LocalDev", false)
	this.Version = VersionMajor + "." + VersionMinor + "." + this.patch

	return this
}

/**
* Simple Getters
 */
func (cfg *Config) GetConnectionString() string {
	return cfg.connectionString
}

func (cfg *Config) GetDatabaseName() string {
	return cfg.databaseName
}

func (cfg *Config) GetPeopleCollectionName() string {
	return cfg.peopleCollectionName
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

	// Start with default values
	theValue := defaultValue
	from := "default"

	// Check for Config File
	// if file.exists(cfg.configFolder/key) {
	// 	theValue = slurp file
	// 	from = "file"
	// }

	// Check for Environemt Variable
	// if ENV.exists(key) {
	// 	theValue = ENV KEY
	// 	from = "environment"
	// }

	// Create the CI and add it to the list
	theItem := &ConfigItem{}
	theItem.Name = key
	theItem.From = from
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
