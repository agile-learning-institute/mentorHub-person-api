package stores

import (
	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"

	"github.com/stretchr/testify/mock"
)

// MockPersonStore is a mock type for the IPersonStore interface
type MockPersonStore struct {
	mock.Mock
}

func (m *MockPersonStore) Insert(information []byte, crumb *models.BreadCrumb) (string, error) {
	args := m.Called(information, crumb)
	return args.String(0), args.Error(1)
}

func (m *MockPersonStore) FindId(id string) (*models.Person, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonStore) UpdateId(id string, request []byte, crumb *models.BreadCrumb) (*models.Person, error) {
	args := m.Called(id, request, crumb)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonStore) FindNames() ([]*config.ShortName, error) {
	args := m.Called()
	return args.Get(0).([]*config.ShortName), args.Error(1)
}
