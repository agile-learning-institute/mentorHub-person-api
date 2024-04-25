// mongo_io_mock.go
package stores

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoIO struct {
	mock.Mock
}

func (m *MockMongoIO) InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(collection, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoIO) FindOne(collection *mongo.Collection, query interface{}, result interface{}) error {
	args := m.Called(collection, query, result)
	return args.Error(0)
}

func (m *MockMongoIO) UpdateOne(collection *mongo.Collection, query interface{}, opts *options.FindOneAndUpdateOptions, update interface{}, result interface{}) error {
	args := m.Called(collection, query, opts, update, result)
	return args.Error(0)
}
