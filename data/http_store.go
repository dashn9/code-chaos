package data

import (
	"encoding/json"
	"fmt"

	"github.com/Ishogbon/code-chaos/schema"
)

// HTTPResponse represents the response from an HTTP action
type HTTPResponse struct {
	StatusCode             int
	Body                   []byte
	Headers                map[string][]string
	ResponseDataTypeMarker *schema.ResponseDataType
}

// ResponseStore manages HTTP responses for both test and generate actions
type ResponseStore struct {
	testResponses     map[string]HTTPResponse // key: ResponseDataType_id
	generateResponses map[string]HTTPResponse // key: ResponseDataType_id
}

var Store = &ResponseStore{
	testResponses:     make(map[string]HTTPResponse),
	generateResponses: make(map[string]HTTPResponse),
}

// StoreTestResponse stores an HTTP response for a test action
func (s *ResponseStore) StoreTestResponse(responseDataTypeMarker *schema.ResponseDataType, testId int, response HTTPResponse) {
	key := fmt.Sprintf("%d_%s", testId, responseDataTypeMarker.ID)
	s.testResponses[key] = response
}

// StoreGenerateResponse stores an HTTP response for a generate action
func (s *ResponseStore) StoreGenerateResponse(responseDataTypeMarker *schema.ResponseDataType, generateId int, response HTTPResponse) {
	key := fmt.Sprintf("%d_%s", generateId, responseDataTypeMarker.ID)
	s.generateResponses[key] = response
}

// GetTestResponse retrieves an HTTP response for a test action
func (s *ResponseStore) GetTestResponse(testID int, responseDataTypeMarkerID string) (HTTPResponse, error) {
	key := fmt.Sprintf("%d_%s", testID, responseDataTypeMarkerID)
	response, exists := s.testResponses[key]
	if !exists {
		return HTTPResponse{}, fmt.Errorf("response not found for test %d, responseDataTypeMarkerID %s", testID, responseDataTypeMarkerID)
	}
	return response, nil
}

// GetGenerateResponse retrieves an HTTP response for a generate action
func (s *ResponseStore) GetGenerateResponse(generateID int, responseDataTypeMarkerID string) (HTTPResponse, error) {
	key := fmt.Sprintf("%d_%s", generateID, responseDataTypeMarkerID)
	response, exists := s.generateResponses[key]
	if !exists {
		return HTTPResponse{}, fmt.Errorf("response not found for generate %d, responseDataTypeMarkerID %s", generateID, responseDataTypeMarkerID)
	}
	return response, nil
}

func (s *ResponseStore) GetResponse(procedureType string, procedureID int, responseDataTypeMarkerID string) (HTTPResponse, error) {
	if procedureType == "test" {
		return s.GetTestResponse(procedureID, responseDataTypeMarkerID)
	} else if procedureType == "generate" {
		return s.GetGenerateResponse(procedureID, responseDataTypeMarkerID)
	}
	return HTTPResponse{}, fmt.Errorf("invalid procedure type: %s", procedureType)
}

func (s *ResponseStore) StoreResponse(procedureType string, procedureID int, responseDataTypeMarker *schema.ResponseDataType, response HTTPResponse) {
	if procedureType == "test" {
		s.StoreTestResponse(responseDataTypeMarker, procedureID, response)
	} else if procedureType == "generate" {
		s.StoreGenerateResponse(responseDataTypeMarker, procedureID, response)
	}
}

func (s *ResponseStore) AccessFieldInResponseData(procedureType string, procedureID int, responseDataTypeMarkerID string, accessPaths []string) (interface{}, error) {
	response, err := s.GetResponse(procedureType, procedureID, responseDataTypeMarkerID)
	if err != nil {
		return nil, err
	}

	if response.ResponseDataTypeMarker.Type != "json" {
		return nil, fmt.Errorf("unsupported response data type: %s", response.ResponseDataTypeMarker.Type)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(response.Body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// Traverse the access path
	var current interface{} = jsonData
	for i, path := range accessPaths {
		obj, ok := current.(map[string]interface{})
		if !ok {
			if i == len(accessPaths)-1 {
				// Final path, return whatever it is (even if not a map)
				return current, nil
			}
			return nil, fmt.Errorf("path '%s' is not a JSON object and more paths remain", path)
		}

		value, exists := obj[path]
		if !exists {
			return nil, fmt.Errorf("path not found: %s", path)
		}
		current = value
	}

	return current, nil
}

// ClearTestResponses clears all stored test responses
func (s *ResponseStore) ClearTestResponses() {
	s.testResponses = make(map[string]HTTPResponse)
}

// ClearGenerateResponses clears all stored generate responses
func (s *ResponseStore) ClearGenerateResponses() {
	s.generateResponses = make(map[string]HTTPResponse)
}

// ClearAllResponses clears all stored responses
func (s *ResponseStore) ClearAllResponses() {
	s.testResponses = make(map[string]HTTPResponse)
	s.generateResponses = make(map[string]HTTPResponse)
}
