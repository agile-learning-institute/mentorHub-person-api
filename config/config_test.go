package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg)

	json := cfg.ToJSONStruct()
	assert.NotNil(t, json)
	assert.Equal(t, DefaultConfigFolder, json.ConfigFolder)
	assert.Equal(t, DefaultDatabaseName, json.DatabaseName)
	assert.Equal(t, DefaultTimeout, json.DatabaseTimeout)
	assert.NotNil(t, json.ApiVersion)
	assert.NotNil(t, json.DataVersion)

	collection := cfg.GetPeopleCollection()
	assert.NotNil(t, collection)

	filter := bson.M{}
	ctx, cancel := cfg.GetTimeoutContext()
	defer cancel()
	_, err := collection.CountDocuments(ctx, filter)

	assert.Nil(t, err)

	cfg.Disconnect() // Clean up, disconnect the MongoDB client
}
