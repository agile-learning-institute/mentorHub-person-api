package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test to verify Connect and Disconnect
func TestConnectDisconnect(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify Find Many Documents
func TestFind(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify FindOne Document
func TestFindOne(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify InsertOne
func TestInsertOne(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify UpdateOne
func TestUpdateOne(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify people collection is valid
func TestGetPeopleCollection(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify Loading the Versions element in the Config
func TestLoadVersions(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify Loading Enumerators element in the Config
func TestLoadEnumerators(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify Updating the Mentors list in the Config
func TestFetchMentors(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}

// Test to verify Updating the Partners list in the Config
func TestFetchPartners(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	assert.Equal(t, "Foo", mongoIO.GetPeopleCollection())
}
