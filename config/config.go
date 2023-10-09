package config

import (
	"context"
	"os"
	"strconv"
	"time"
)

type ConfigItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	From  string `json:"from"`
}

type Config struct {
	ConfigItems          []*configItems
	Version              string
	DBVersion            string
	port                 string
	patch                string
	configFolder         string
	databaseName         string
	peopleCollectionName string
	databaseTimeout      int
	connectionString     string
	// person               models.PersonInterface
	// personStore          models.PersonStoreInterface
}

const (
	VersionMajor                = "1"
	VersionMinor                = "0"
	DefaultConfigFolder         = "./"
	DefaultConnectionString     = "mongodb://root:example@localhost:27017/?tls=false&directConnection=true"
	DefaultDatabaseName         = "agile-learning-institute"
	DefaultPeopleCollectionName = "people"
	DefaultPort                 = ":8080"
	DefaultTimeout              = 10
)

/**
* Construct a config item by finding all the configuration values
 */
func NewConfig() *Config {
	this := &Config{}
	this.patch = this.findStringValue("PATCH_LEVEL", "LocalDev", false)
	this.configFolder = this.findStringValue("CONFIG_FOLDER", DefaultConfigFolder, false)
	this.connectionString = this.findStringValue("CONNECTION_STRING", DefaultConnectionString, true)
	this.databaseName = this.findStringValue("DATABASE_NAME", DefaultDatabaseName, false)
	this.peopleCollectionName = this.findStringValue("PEOPLE_COLLECTION_NAME", DefaultPeopleCollectionName, false)
	this.databaseTimeout = this.findIntValue("CONNECTION_TIMEOUT", DefaultTimeout, false)
	this.port = this.findStringValue("PORT", DefaultPort, false)
	this.Version = VersionMajor + "." + VersionMinor + "." + this.patch
	this.DBVersion = "Pending"
	return this
}

/**
* Simple Getters - Read Only attributes
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

func (cfg *Config) GetPort() string {
	return cfg.port
}

/**
* the one and only Setter
 */
func (cfg *Config) SetDbVersion(theVersion string) {
	cfg.DBVersion = theVersion
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
