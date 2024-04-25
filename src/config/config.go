/********************************************************************************
** Config
**    The Config class manages configuration values that are initilized from
**	    descrete configuration files, the environment, or defalut values.
********************************************************************************/
package config

import (
	"log"
	"mentorhub-person-api/src/models"
	"os"
	"strconv"
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
)

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
	* Then add it to the ConfigItems array
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

	/**********************************************************************
	* Initilize all configuration values
	 */
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
* Simple Getters
 */
func (cfg *Config) GetPort() string {
	return cfg.port
}

func (cfg *Config) GetPatch() string {
	return cfg.patch
}

func (cfg *Config) GetConfigFolder() string {
	return cfg.configFolder
}

func (cfg *Config) GetConnectionString() string {
	return cfg.connectionString
}

func (cfg *Config) GetDatabaseName() string {
	return cfg.databaseName
}

func (cfg *Config) GetDatabaseTimeout() int {
	return cfg.databaseTimeout
}
