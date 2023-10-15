package tests

import (
	"institute-person-api/handlers"
	"institute-person-api/mocks"
	"institute-person-api/models"
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
	mockPerson := mocks.NewMockPersonInterface(ctrl)

	// Invoke NewPerson
	handler := handlers.NewPersonHandler(mockPerson)

	// Examine the result
	assert.NotNil(t, handler)
	assert.Equal(t, mockPerson, handler.Person)
}

func TestAddPerson(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPerson := mocks.NewMockPersonInterface(ctrl)
	handler := handlers.PersonHandler{Person: mockPerson}

	// Tell the Mock how to respond to PostPerson
	personJSON := `{"name": "John Doe", "description": "Test Person"}`
	request := httptest.NewRequest("POST", "/person", strings.NewReader(personJSON))
	responseRecorder := httptest.NewRecorder()
	mockPerson.EXPECT().PostPerson([]byte(personJSON), gomock.Any(), gomock.Any()).Return(&models.Person{Name: "John Doe", Description: "Test Person"}, nil)

	// Invoke NewPerson
	handler.AddPerson(responseRecorder, request)

	// Examine the result
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"ID\":\"000000000000000000000000\",\"name\":\"John Doe\",\"description\":\"Test Person\"}\n", responseRecorder.Body.String())
}

func TestGetPersonWithHandler(t *testing.T) {
	// Setup the Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPerson := mocks.NewMockPersonInterface(ctrl)
	handler := handlers.PersonHandler{Person: mockPerson}

	// Tell the Mock how to respond to GetPerson
	// id := primitive.NewObjectID().Hex()
	personJSON := `{"name": "John Doe", "description": "Test Person"}`
	request := httptest.NewRequest("GET", "/person/000000000000000000000000/", strings.NewReader(personJSON))
	responseRecorder := httptest.NewRecorder()
	mockPerson.EXPECT().GetPerson(gomock.Any()).Return(&models.Person{Name: "John Doe", Description: "Test Person"}, nil)

	// Invoke NewPerson
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
	mockPerson := mocks.NewMockPersonInterface(ctrl)
	handler := handlers.PersonHandler{Person: mockPerson}

	// Tell the Mock how to respond to GetPerson
	// id := primitive.NewObjectID().Hex()
	expectedNames := []models.PersonShort{
		{Name: "Mock Name 1"},
		{Name: "Mock Name 2"},
	}
	request := httptest.NewRequest("GET", "/person/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	mockPerson.EXPECT().GetAllNames().Return(expectedNames, nil)

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
	mockPerson := mocks.NewMockPersonInterface(ctrl)
	handler := handlers.PersonHandler{Person: mockPerson}

	// Tell the Mock how to respond to GetPerson
	personJSON := `"name": "John Doe", "description": "Test Person"}`
	request := httptest.NewRequest("PATCH", "/person/000000000000000000000000/", strings.NewReader(personJSON))
	responseRecorder := httptest.NewRecorder()
	mockPerson.EXPECT().PatchPerson(gomock.Any(), []byte(personJSON), gomock.Any(), gomock.Any()).Return(&models.Person{Name: "John Doe", Description: "Test Person"}, nil)

	// Invoke NewPerson
	handler.UpdatePerson(responseRecorder, request)

	// Examine the result
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"ID\":\"000000000000000000000000\",\"name\":\"John Doe\",\"description\":\"Test Person\"}\n", responseRecorder.Body.String())
}
