/********************************************************************************
** MongoIO_test
**    This testing requires a local mongoDb with test data
**	  You can start this contianer with the command "mh up mongodb"
********************************************************************************/
package config

import (
	"log"
	"mentorhub-person-api/src/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test to verify Connect and Disconnect
func TestConnectDisconnect(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)

	mongoIO.Connect()
	defer mongoIO.Disconnect()
	assert.Equal(t, "people", mongoIO.GetPeopleCollection().Name())

	mongoIO.Disconnect()
	assert.Nil(t, mongoIO.GetPeopleCollection())
}

// Test to verify CRU functionality
func TestInsertOne(t *testing.T) {
	// Initilize Config, connect to database
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	var people = mongoIO.GetPeopleCollection()
	assert.Equal(t, "people", people.Name())

	// Insert a Person
	testName := uuid.New().String()[:31]
	input := struct {
		UserName string `bson:"userName,omitempty"`
	}{UserName: testName}
	log.Print(input.UserName)
	output, err := mongoIO.InsertOne(people, input)
	assert.Nil(t, err)
	assert.NotNil(t, output.InsertedID)

	mongoIO.Disconnect()
}

func TestFindOne(t *testing.T) {
	// Initilize Config, connect to database
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	var people = mongoIO.GetPeopleCollection()
	assert.Equal(t, "people", people.Name())

	// Get a Person from the test data
	ID, err := primitive.ObjectIDFromHex("aaaa00000000000000000000")
	assert.Nil(t, err)

	output := models.Person{}
	query := bson.M{"_id": ID}

	err = mongoIO.FindOne(people, query, &output)
	assert.Nil(t, err)
	assert.Equal(t, "JamesSmith", output.UserName)
	assert.Equal(t, "She was too busy always talking about what she wanted to do to actually do any of it.", output.Description)

	mongoIO.Disconnect()
}

func TestUpdateOne(t *testing.T) {
	// Initilize Config, connect to database
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	var people = mongoIO.GetPeopleCollection()
	assert.Equal(t, "people", people.Name())

	// Update a Person
	ID, err := primitive.ObjectIDFromHex("AAAA00000000000000000000")
	assert.Nil(t, err)

	input := models.Person{Description: "Updated"}
	output := models.Person{}
	query := bson.M{"_id": ID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = mongoIO.UpdateOne(people, query, opts, input, &output)
	assert.NotNil(t, err)
	assert.Equal(t, "Foo", output.UserName)
	assert.Equal(t, "Updated", output.Description)

	mongoIO.Disconnect()
}

// Test to verify Get People Names
func TestFindNames(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	assert.Equal(t, "people", mongoIO.GetPeopleCollection().Name())

	results := []*ShortName{}
	query := bson.M{}
	opts := options.Find()
	opts.SetProjection(NameProjection())
	err := mongoIO.Find(mongoIO.GetPeopleCollection(), query, opts, &results)
	assert.Nil(t, err)
	assert.Greater(t, len(results), 0)

	mongoIO.Disconnect()
}

// Test to verify Loading the Versions element in the Config
func TestLoadVersions(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	assert.Equal(t, "people", mongoIO.GetPeopleCollection().Name())

	mongoIO.LoadVersions()
	assert.Greater(t, len(cfg.Versions), 0)
	assert.NotNil(t, cfg.Versions[0].CollectionName)
	assert.NotNil(t, cfg.Versions[0].CurrentVersion)

	mongoIO.Disconnect()
}

// Test to verify Loading Enumerators element in the Config
func TestLoadEnumerators(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	assert.Equal(t, "people", mongoIO.GetPeopleCollection().Name())

	mongoIO.LoadVersions()
	mongoIO.LoadEnumerators()
	defaultStatus := cfg.Enumerators["defaultStatus"].(map[string]interface{})
	assert.NotNil(t, defaultStatus)
	active := defaultStatus["Active"].(string)
	assert.Equal(t, "Not Deleted", active)

	mongoIO.Disconnect()
}

// Test to verify Updating the Mentors list in the Config
func TestFetchMentors(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	assert.Equal(t, "people", mongoIO.GetPeopleCollection().Name())

	mongoIO.FetchMentors()
	assert.Greater(t, len(cfg.Mentors), 0)
	assert.NotNil(t, cfg.Mentors[0].ID)
	assert.NotNil(t, cfg.Mentors[0].Name)
}

// Test to verify Updating the Partners list in the Config
func TestFetchPartners(t *testing.T) {
	cfg := NewConfig()
	mongoIO := NewMongoIO(cfg)
	mongoIO.Connect()
	defer mongoIO.Disconnect()
	assert.Equal(t, "people", mongoIO.GetPeopleCollection().Name())

	mongoIO.FetchPartners()
	assert.Greater(t, len(cfg.Partners), 0)
	assert.NotNil(t, cfg.Partners[0].ID)
	assert.NotNil(t, cfg.Partners[0].Name)
}
