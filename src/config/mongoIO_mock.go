package config

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoIO struct {
	mock.Mock
	Config *Config // Add the Config object
}

func NewMockMongoIO(cfg *Config) *MockMongoIO {
	mock := &MockMongoIO{}
	// Use cfg to set up mock as needed
	return mock
}

// Implement the MongoIOInterface methods using mock calls
func (m *MockMongoIO) Connect() {
}

func (m *MockMongoIO) Disconnect() {
}

func (m *MockMongoIO) Find(collection *mongo.Collection, query bson.M, opts *options.FindOptions, results interface{}) error {
	args := m.Called(collection, results)
	return args.Error(0)
}

func (m *MockMongoIO) FindOne(collection *mongo.Collection, query bson.M, results interface{}) error {
	args := m.Called(collection, query, results)
	return args.Error(0)
}

func (m *MockMongoIO) InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(collection, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoIO) UpdateOne(collection *mongo.Collection, query bson.M, opts *options.FindOneAndUpdateOptions, update interface{}, results interface{}) error {
	args := m.Called(collection, query, opts, update, results)
	return args.Error(0)
}
