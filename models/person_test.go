package models

import (
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
	mockStore := NewMockPersonStoreInterface(ctrl)

	// Invoke NewPerson
	person := NewPerson(mockStore)

	// Examine the result
	assert.NotNil(t, person)
}

func TestGetPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockPersonStoreInterface(ctrl)

	// Configure Mock Response
	id := primitive.NewObjectID().Hex()
	newId, _ := primitive.ObjectIDFromHex(id)
	match := gomock.Eq(bson.M{"_id": newId})
	expectedPerson := &Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		store:       mockStore,
	}
	mockStore.EXPECT().FindOne(match).Return(expectedPerson)

	// Invoke GetPerson
	person := NewPerson(mockStore)
	result := person.GetPerson(id)

	// Examine the results of the invocation
	assert.Equal(t, expectedPerson, result)
}

func TestGetAllNames(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockPersonStoreInterface(ctrl)

	// Configure Mock return value
	match := bson.M{}
	options := gomock.Not(gomock.Nil())
	expectedNames := []PersonShort{
		{Name: "Mock Name 1"},
		{Name: "Mock Name 2"},
	}
	mockStore.EXPECT().FindMany(match, options).Return(expectedNames)

	// Invoke NewPerson
	person := NewPerson(mockStore)
	result := person.GetAllNames()

	// Examine the results of the invocation
	assert.Equal(t, expectedNames, result)
}

func TestPostPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockPersonStoreInterface(ctrl)

	// Configure Mock return values for Insert
	id := primitive.NewObjectID().Hex()
	newId, _ := primitive.ObjectIDFromHex(id)
	expectedResult := &mongo.InsertOneResult{InsertedID: newId}
	mockStore.EXPECT().Insert(gomock.Any()).Return(expectedResult)

	// Configure Mock return values for FindOne
	expectedPerson := &Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		store:       mockStore,
	}
	match := gomock.Eq(bson.M{"_id": newId})
	mockStore.EXPECT().FindOne(match).Return(expectedPerson)

	// Invoke Post Person
	var body []byte
	person := NewPerson(mockStore)
	result := person.PostPerson(body)

	// Examine the results of the invocation
	assert.Equal(t, expectedPerson, result)
}

func TestPatchPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockPersonStoreInterface(ctrl)

	// Configure Mock response to FindOneAndUpdate
	id := primitive.NewObjectID().Hex()
	newId, _ := primitive.ObjectIDFromHex(id)
	match := gomock.Eq(bson.M{"_id": newId})
	expectedPerson := &Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		store:       mockStore,
	}
	mockStore.EXPECT().FindOneAndUpdate(match, gomock.Any()).Return(expectedPerson)

	// Invoke Patch Person
	var body []byte
	person := NewPerson(mockStore)
	result := person.PatchPerson(id, body)

	// Examine the results of the invocation
	assert.Equal(t, expectedPerson, result)
}
