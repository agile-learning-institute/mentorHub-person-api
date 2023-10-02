package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigJSON struct {
	ConfigFolder    string `json:"configFolder"`
	ApiVersion      string `json:"apiVersion"`
	DataVersion     string `json:"dataVersion"`
	DatabaseName    string `json:"databaseName"`
	DatabaseTimeout int    `json:"databaseTimeout"`
}

type Config struct {
	configFolder     string
	apiVersion       string
	dataVersion      string
	databaseName     string
	peopleCollection string
	databaseTimeout  int
	connectionString string
	client           *mongo.Client
	database         *mongo.Database
	cancel           context.CancelFunc
}

const (
	DefaultConfigFolder     = "/opt/"
	DefaultConnectionString = "mongodb://root:example@localhost:27017/?tls=false&directConnection=true"
	DefaultDatabaseName     = "agile-learning-institute"
	DefaultPeopleCollection = "people"
	DefaultTimeout          = 10
)

func NewConfig() *Config {
	this := &Config{}
	this.configFolder = this.findStringValue("CONFIG_FOLDER", DefaultConfigFolder)
	this.connectionString = this.findStringValue("CONNECTION_STRING", DefaultConnectionString)
	this.databaseName = this.findStringValue("DATABASE_NAME", DefaultDatabaseName)
	this.databaseTimeout = this.findIntValue("CONNECTION_TIMEOUT", DefaultTimeout)
	this.peopleCollection = this.findStringValue("PEPOLE_COLLECTION", DefaultPeopleCollection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(this.databaseTimeout)*time.Second)
	this.cancel = cancel

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(this.connectionString))
	if err != nil {
		log.Fatal(err)
	}
	this.client = client

	this.database = client.Database(this.databaseName)
	this.apiVersion = "LocalDev"  // TODO: Read apiVersion from files
	this.dataVersion = "LocalDev" // TODO: Load DB Version from Database

	return this
}

func (cfg *Config) Disconnect() {
	if err := cfg.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
	cfg.cancel()
}

func (cfg *Config) GetPeopleCollection() *mongo.Collection {
	return cfg.database.Collection(cfg.peopleCollection)
}

func (cfg *Config) GetTimeoutContext() (context.Context, context.CancelFunc) {
	timeout := time.Duration(cfg.databaseTimeout) * time.Second
	return context.WithTimeout(context.Background(), timeout)
}

func (cfg *Config) ToJSONStruct() ConfigJSON {
	return ConfigJSON{
		ConfigFolder:    cfg.configFolder,
		ApiVersion:      cfg.apiVersion,
		DataVersion:     cfg.dataVersion,
		DatabaseName:    cfg.databaseName,
		DatabaseTimeout: cfg.databaseTimeout,
	}
}

func (cfg *Config) findStringValue(key string, defaultValue string) string {
	theValue := defaultValue
	// TODO: Look in Environemnt for key name
	// TODO: Look in cfg.ConfigFolder for key file
	return theValue
}

func (cfg *Config) findIntValue(key string, defaultValue int) int {
	theValue := defaultValue
	// TODO: Look in Environemnt for key name
	// TODO: Look in cfg.ConfigFolder for key file
	return theValue
}
