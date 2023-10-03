package main

// Import the GoMock library
import (
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the mock
type MockPersonStore struct {
	*gomock.Controller
	store PersonStoreInterface
}

// MockPersonStoreConstructor
func NewMockPersonStore(ctrl *gomock.Controller) *MockPersonStore {
	return &MockPersonStore{Controller: ctrl}
}

// Implement the PersonStore interface
func (m *MockPersonStore) FindOne(query interface{}) *Person {
	// Mock implementation to return a fixed Person object
	return &Person{
		ID:          primitive.NewObjectID(),
		Name:        "Mock Name",
		Description: "Mock Description",
	}
}

func (m *MockPersonStore) FindMany(query interface{}, opts *options.FindOptions) *[]PersonShort {
	// Mock implementation to return a fixed slice of PersonShort objects
	return &[]PersonShort{
		{ID: primitive.NewObjectID(), Name: "Mock Name 1"},
		{ID: primitive.NewObjectID(), Name: "Mock Name 2"},
	}
}

func (m *MockPersonStore) Insert(values interface{}) *mongo.InsertOneResult {
	// Mock implementation to return a fixed InsertOneResult object
	return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}
}

func (m *MockPersonStore) FindOneAndUpdate(query interface{}, update interface{}) *Person {
	// Mock implementation to return a fixed Person object with updated fields
	updateValues := update.(bson.M)["$set"].(bson.M)
	return &Person{
		ID:          query.(bson.M)["_id"].(primitive.ObjectID),
		Name:        updateValues["name"].(string),
		Description: updateValues["description"].(string),
	}
}
