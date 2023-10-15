package tests

import (
	"institute-person-api/mocks"
	"institute-person-api/models"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)

	// Invoke NewPerson
	person := models.NewPerson(mockStore)

	// Examine the result
	assert.NotNil(t, person)
}

func TestGetPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)

	// Configure Mock Response
	id := primitive.NewObjectID().Hex()
	newId, _ := primitive.ObjectIDFromHex(id)
	match := gomock.Eq(bson.M{"_id": newId})
	expectedPerson := &models.Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		Store:       mockStore,
	}
	mockStore.EXPECT().FindOne(match).Return(expectedPerson, nil)

	// Invoke GetPerson
	person := models.NewPerson(mockStore)
	result, err := person.GetPerson(id)

	// Examine the results of the invocation
	log.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, expectedPerson, result)
}

func TestGetAllNames(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)

	// Configure Mock return value
	match := bson.M{"name": bson.M{"$ne": "VERSION"}}
	options := gomock.Not(gomock.Nil())
	expectedNames := []models.PersonShort{
		{Name: "Mock Name 1"},
		{Name: "Mock Name 2"},
	}
	mockStore.EXPECT().FindMany(match, options).Return(expectedNames, nil)

	// Invoke NewPerson
	person := models.NewPerson(mockStore)
	result, err := person.GetAllNames()

	// Examine the results of the invocation
	assert.Nil(t, err)
	assert.Equal(t, expectedNames, result)
}

func TestPostPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)

	// Configure Mock return values for Insert
	id := primitive.NewObjectID().Hex()
	newId, _ := primitive.ObjectIDFromHex(id)
	expectedResult := &mongo.InsertOneResult{InsertedID: newId}
	mockStore.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

	// Configure Mock return values for FindOne
	expectedPerson := &models.Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		Store:       mockStore,
	}
	match := gomock.Eq(bson.M{"_id": newId})
	mockStore.EXPECT().FindOne(match).Return(expectedPerson, nil)

	// Invoke Post Person
	json := "{}"
	body := []byte(json)
	person := models.NewPerson(mockStore)
	result, err := person.PostPerson(body, "")

	// Examine the results of the invocation
	assert.Nil(t, err)
	assert.Equal(t, expectedPerson, result)
}

func TestPatchPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)

	// Configure Mock response to FindOneAndUpdate
	id := primitive.NewObjectID().Hex()
	newId, _ := primitive.ObjectIDFromHex(id)
	match := gomock.Eq(bson.M{"_id": newId})
	expectedPerson := &models.Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		Store:       mockStore,
	}
	mockStore.EXPECT().FindOneAndUpdate(match, gomock.Any(), gomock.Any()).Return(expectedPerson, nil)

	// Invoke Patch Person
	json := "{}"
	body := []byte(json)
	person := models.NewPerson(mockStore)
	result, err := person.PatchPerson(id, body, "")

	// Examine the results of the invocation
	assert.Nil(t, err)
	assert.Equal(t, expectedPerson, result)
}
