// person_store_test.go
package stores

import (
	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInsert(t *testing.T) {
	cfg := config.NewConfig()
	mockMongoIO := config.NewMockMongoIO(cfg)
	store := NewPersonStore(mockMongoIO)

	// Prepare the input and other test setups
	information := []byte(`{"userName":"Test", "description":"New User"}`)
	crumb, err := models.NewBreadCrumb("127.0.0.1", "123456789012345678901234", "CORRID")
	assert.Nil(t, err)

	// Setup the mock expectations and return values
	insertedID := primitive.NewObjectID()
	insertResult := &mongo.InsertOneResult{InsertedID: insertedID}
	mockMongoIO.On("InsertOne", mock.AnythingOfType("*mongo.Collection"), mock.AnythingOfType("bson.M")).Return(insertResult, nil)

	// Run the test
	resultPerson, err := store.Insert(information, crumb)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, resultPerson)
	assert.Equal(t, insertedID.Hex(), resultPerson.ID.Hex())
	mockMongoIO.AssertExpectations(t)
}
