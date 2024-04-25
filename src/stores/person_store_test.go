// person_store_test.go
package stores

import (
	"testing"

	"mentorhub-person-api/src/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInsert(t *testing.T) {
	mockMongoIO := new(MockMongoIO)
	store := NewPersonStore(mockMongoIO)
	// mockCollection := &mongo.Collection{}

	// Prepare the input
	information := []byte(`{"description":"Updated","mentorId":"5f3fc1ef1e9a3225a4a39bf4"}`)
	crumb := &models.BreadCrumb{}

	// Mock InsertOne
	insertResult := &mongo.InsertOneResult{
		InsertedID: primitive.NewObjectID(),
	}
	mockMongoIO.On("InsertOne", mock.AnythingOfType("*mongo.Collection"), mock.Anything).Return(insertResult, nil)

	// Run the test
	person, err := store.Insert(information, crumb)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, person)
	mockMongoIO.AssertExpectations(t)
}

func TestFindId(t *testing.T) {
	mockMongoIO := new(MockMongoIO)
	store := NewPersonStore(mockMongoIO)
	// mockCollection := &mongo.Collection{} // Mock collection

	id := "5f3fc1ef1e9a3225a4a39bf4"
	expectedPerson := &models.Person{UserName: "TestUser"}

	// Mock FindOne
	mockMongoIO.On("FindOne", mock.AnythingOfType("*mongo.Collection"), mock.Anything, mock.AnythingOfType("*models.Person")).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*models.Person)
		*arg = *expectedPerson
	}).Return(nil)

	// Run the test
	person, err := store.FindId(id)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedPerson, person)
	mockMongoIO.AssertExpectations(t)
}

// Further tests can be written for FindOneAndUpdate and other methods in a similar fashion.
