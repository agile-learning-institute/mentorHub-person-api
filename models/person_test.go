package models

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewPerson(t *testing.T) {
	// Prepare the ritual circle
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Conjure the mock PersonStore
	mockStore := NewMockPersonStore(ctrl)

	// Invoke NewPerson
	person := NewPerson(mockStore)

	// Examine the conjured Person
	assert.NotNil(t, person)
	assert.Equal(t, mockStore, person.store)
}

func TestGetPerson(t *testing.T) {
	// Prepare the ritual circle
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Conjure the mock PersonStore
	mockStore := NewMockPersonStore(ctrl)

	// Prepare the ritual components
	id := primitive.NewObjectID().Hex()
	expectedPerson := &Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
		store:       mockStore,
	}

	// Foretell the behavior of the mock PersonStore
	newId, _ := primitive.ObjectIDFromHex(id)
	match := bson.M{"_id": gomock.Eq(newId)}
	mockStore.EXPECT().FindOne(match).Return(expectedPerson)

	// Invoke GetPerson
	person := NewPerson(mockStore)
	result := person.GetPerson(id)

	// Examine the results of the invocation
	assert.Equal(t, expectedPerson, result)
}
