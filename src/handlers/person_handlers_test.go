package handlers

import (
	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"
	"mentorhub-person-api/src/stores"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConstructor(t *testing.T) {
	mockStore := new(stores.MockPersonStore)
	personHandler := NewPersonHandler(mockStore)

	assert.NotNil(t, personHandler)
}

func TestGetPerson(t *testing.T) {
	// Setup
	mockStore := new(stores.MockPersonStore)
	personHandler := NewPersonHandler(mockStore)
	request := httptest.NewRequest("GET", "/person/", nil)
	responseRecorder := httptest.NewRecorder()
	person := &models.Person{}

	// Initilize Mock
	mockStore.On("FindId", mock.Anything).Return(person, nil)

	// Invoke GetPerson
	personHandler.GetPerson(responseRecorder, request)

	// Examine the result
	assert.NotNil(t, personHandler)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
}

func TestAddPersonHandler(t *testing.T) {
	// Setup
	mockStore := new(stores.MockPersonStore)
	personHandler := NewPersonHandler(mockStore)
	request := httptest.NewRequest("GET", "/config/", nil)
	responseRecorder := httptest.NewRecorder()
	person := &models.Person{}

	// Initilize Mock
	mockStore.On("Insert", mock.Anything, mock.Anything).Return("000000000000000000000000", nil)
	mockStore.On("FindId", mock.Anything).Return(person, nil)

	// Invoke AddLPerson
	personHandler.AddPerson(responseRecorder, request)

	// Examine the result
	assert.NotNil(t, personHandler)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
}

func TestUpdatePersonHandler(t *testing.T) {
	// Setup
	mockStore := new(stores.MockPersonStore)
	personHandler := NewPersonHandler(mockStore)
	request := httptest.NewRequest("GET", "/config/", nil)
	responseRecorder := httptest.NewRecorder()
	person := &models.Person{}

	// Initilize Mock
	mockStore.On("UpdateId", mock.Anything, mock.Anything, mock.Anything).Return(person, nil)

	// Invoke AddLPerson
	personHandler.UpdatePerson(responseRecorder, request)

	// Examine the result
	assert.NotNil(t, personHandler)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
}

func TestGetPeople(t *testing.T) {
	// Setup
	mockStore := new(stores.MockPersonStore)
	personHandler := NewPersonHandler(mockStore)
	request := httptest.NewRequest("GET", "/person/", nil)
	responseRecorder := httptest.NewRecorder()
	findNamesResult := []config.ShortName{}

	// Initilize Mock
	mockStore.On("FindNames").Return(findNamesResult, nil)

	// Invoke GetPerson
	personHandler.GetPeople(responseRecorder, request)

	// Examine the result
	assert.NotNil(t, personHandler)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
}
