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
	mockMongoIO.On("InsertOne",
		mock.Anything,
		mock.Anything,
	).Return(insertResult, nil)

	// Run the test
	resultID, err := store.Insert(information, crumb)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, resultID)
	assert.Equal(t, insertedID.Hex(), resultID)
	mockMongoIO.AssertExpectations(t)
}

func TestFindId(t *testing.T) {
	cfg := config.NewConfig()
	mockMongoIO := config.NewMockMongoIO(cfg)
	store := NewPersonStore(mockMongoIO)

	// Setup the Mock
	findResult := &models.Person{
		UserName:    "someUser",
		Description: "Some Document",
	}
	mockMongoIO.On("FindOne",
		mock.AnythingOfType("*mongo.Collection"),
		mock.Anything,
		mock.AnythingOfType("*models.Person"),
	).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*models.Person)
		*arg = *findResult
	}).Return(nil)

	// Run the test
	result, err := store.FindId("123456789012345678901234")

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "someUser", result.UserName)
	assert.Equal(t, "Some Document", result.Description)
}

func TestUpdateId(t *testing.T) {
	cfg := config.NewConfig()
	mockMongoIO := config.NewMockMongoIO(cfg)
	store := NewPersonStore(mockMongoIO)

	// Prepare the input and other test setups
	information := []byte(`{"userName":"Test", "description":"Updated User"}`)
	crumb, err := models.NewBreadCrumb("127.0.0.1", "123456789012345678901234", "CORRID")
	assert.Nil(t, err)

	// Setup the Mock
	findOneAndUpdateResult := &models.Person{
		UserName:    "someUser",
		Description: "Some Document",
	}
	mockMongoIO.On("UpdateOne",
		mock.AnythingOfType("*mongo.Collection"),
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*models.Person"),
	).Run(func(args mock.Arguments) {
		arg := args.Get(4).(*models.Person)
		*arg = *findOneAndUpdateResult
	}).Return(nil)

	// Run the test
	result := &models.Person{}
	result, err = store.UpdateId("123456789012345678901234", information, crumb)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "someUser", result.UserName)
	assert.Equal(t, "Some Document", result.Description)
}

func TestFindNames(t *testing.T) {
	cfg := config.NewConfig()
	mockMongoIO := config.NewMockMongoIO(cfg)
	store := NewPersonStore(mockMongoIO)

	// Setup the Teting Data
	findResult := []*config.ShortName{}
	fooName := &config.ShortName{
		Name: "foo",
		ID:   primitive.NewObjectID(),
	}
	barName := &config.ShortName{
		Name: "bar",
		ID:   primitive.NewObjectID(),
	}
	findResult = append(findResult, fooName)
	findResult = append(findResult, barName)

	// setup the mock
	mockMongoIO.On("Find",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*[]*config.ShortName"),
	).Run(func(args mock.Arguments) {
		// Populate the results in the slice pointed to by the 4th argument
		arg := args.Get(3).(*[]*config.ShortName)
		*arg = findResult
	}).Return(nil)

	// Run the test
	result, err := store.FindNames()

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "foo", result[0].Name)
	assert.Equal(t, "bar", result[1].Name)
}
