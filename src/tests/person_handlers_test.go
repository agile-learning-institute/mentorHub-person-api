package tests

import (
	 "institute-person-api/src/handlers"
	 "institute-person-api/src/mocks"
	 "institute-person-api/src/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewPersonHandler(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)

	// Invoke NewPerson
	handler := handlers.NewPersonHandler(mockStore)

	// Examine the result
	assert.NotNil(t, handler)
	assert.Equal(t, mockStore, handler.PersonStore)
}

func TestAddPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)
	handler := handlers.NewPersonHandler(mockStore)

	// Tell the Mock how to respond to PostPerson
	personJSON := `{"name": "John Doe", "description": "Test Person"}`
	request := httptest.NewRequest("POST", "/person", strings.NewReader(personJSON))
	responseRecorder := httptest.NewRecorder()
	mockStore.EXPECT().Insert([]byte(personJSON), gomock.Any()).Return(&models.Person{Name: "John Doe", Description: "Test Person"}, nil)

	// Invoke NewPerson
	handler.AddPerson(responseRecorder, request)

	// Examine the result
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"ID\":\"000000000000000000000000\",\"name\":\"John Doe\",\"description\":\"Test Person\"}\n", responseRecorder.Body.String())
}

func TestGetPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)
	handler := handlers.NewPersonHandler(mockStore)

	// Tell the Mock how to respond to GetPerson
	personJSON := `{"name": "John Doe", "description": "Test Person"}`
	request := httptest.NewRequest("GET", "/person/000000000000000000000000/", strings.NewReader(personJSON))
	responseRecorder := httptest.NewRecorder()
	mockStore.EXPECT().FindOne(gomock.Any()).Return(&models.Person{Name: "John Doe", Description: "Test Person"}, nil)

	// Invoke GetPerson
	handler.GetPerson(responseRecorder, request)

	// Examine the result
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"ID\":\"000000000000000000000000\",\"name\":\"John Doe\",\"description\":\"Test Person\"}\n", responseRecorder.Body.String())
}

func TestGetPeople(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)
	handler := handlers.NewPersonHandler(mockStore)

	// Tell the Mock how to respond to GetPerson
	expectedNames := []models.PersonShort{
		{Name: "Mock Name 1"},
		{Name: "Mock Name 2"},
	}
	request := httptest.NewRequest("GET", "/person/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	mockStore.EXPECT().FindMany().Return(expectedNames, nil)

	// Invoke NewPerson
	handler.GetPeople(responseRecorder, request)

	// Examine the result
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "[{\"ID\":\"000000000000000000000000\",\"name\":\"Mock Name 1\"},{\"ID\":\"000000000000000000000000\",\"name\":\"Mock Name 2\"}]\n", responseRecorder.Body.String())
}

func TestUpdatePerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mocks.NewMockPersonStoreInterface(ctrl)
	handler := handlers.NewPersonHandler(mockStore)

	// Tell the Mock how to respond to GetPerson
	personJSON := `"name": "John Doe", "description": "Test Person"}`
	request := httptest.NewRequest("PATCH", "/person/000000000000000000000000/", strings.NewReader(personJSON))
	responseRecorder := httptest.NewRecorder()
	mockStore.EXPECT().FindOneAndUpdate(gomock.Any(), []byte(personJSON), gomock.Any()).Return(&models.Person{Name: "John Doe", Description: "Test Person"}, nil)

	// Invoke NewPerson
	handler.UpdatePerson(responseRecorder, request)

	// Examine the result
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"ID\":\"000000000000000000000000\",\"name\":\"John Doe\",\"description\":\"Test Person\"}\n", responseRecorder.Body.String())
}
