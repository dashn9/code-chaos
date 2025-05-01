package data

import (
	"encoding/json"
	"fmt"

	"github.com/Ishogbon/code-chaos/db"
	"github.com/Ishogbon/code-chaos/schema"
)

// DBResponse represents the response from a database action
type DBResponse struct {
	Data                   []byte
	ResponseDataTypeMarker *schema.ResponseDataType
}

// DBStore manages database responses for both test and generate actions
type DBStore struct {
	testResponses     map[string]DBResponse // key: ResponseDataType_id
	generateResponses map[string]DBResponse // key: ResponseDataType_id
	db                *db.DB
}

var DB = &DBStore{
	testResponses:     make(map[string]DBResponse),
	generateResponses: make(map[string]DBResponse),
}

// NewDBStore creates a new database store instance
func NewDBStore(dbType string, connectionString string) (*DBStore, error) {
	factory := db.NewDBFactory()
	dbInstance, err := factory.CreateDB(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	DB.db = dbInstance
	return DB, nil
}

// StoreTestResponse stores a database response for a test action
func (s *DBStore) StoreTestResponse(responseDataTypeMarker *schema.ResponseDataType, testId int, response DBResponse) {
	key := fmt.Sprintf("%d_%s", testId, responseDataTypeMarker.ID)
	s.testResponses[key] = response
}

// StoreGenerateResponse stores a database response for a generate action
func (s *DBStore) StoreGenerateResponse(responseDataTypeMarker *schema.ResponseDataType, generateId int, response DBResponse) {
	key := fmt.Sprintf("%d_%s", generateId, responseDataTypeMarker.ID)
	s.generateResponses[key] = response
}

// GetTestResponse retrieves a database response for a test action
func (s *DBStore) GetTestResponse(testID int, responseDataTypeMarkerID string) (DBResponse, error) {
	key := fmt.Sprintf("%d_%s", testID, responseDataTypeMarkerID)
	response, exists := s.testResponses[key]
	if !exists {
		return DBResponse{}, fmt.Errorf("response not found for test %d, responseDataTypeMarkerID %s", testID, responseDataTypeMarkerID)
	}
	return response, nil
}

// GetGenerateResponse retrieves a database response for a generate action
func (s *DBStore) GetGenerateResponse(generateID int, responseDataTypeMarkerID string) (DBResponse, error) {
	key := fmt.Sprintf("%d_%s", generateID, responseDataTypeMarkerID)
	response, exists := s.generateResponses[key]
	if !exists {
		return DBResponse{}, fmt.Errorf("response not found for generate %d, responseDataTypeMarkerID %s", generateID, responseDataTypeMarkerID)
	}
	return response, nil
}

// GetResponse retrieves a database response based on procedure type
func (s *DBStore) GetResponse(procedureType string, procedureID int, responseDataTypeMarkerID string) (DBResponse, error) {
	if procedureType == "test" {
		return s.GetTestResponse(procedureID, responseDataTypeMarkerID)
	} else if procedureType == "generate" {
		return s.GetGenerateResponse(procedureID, responseDataTypeMarkerID)
	}
	return DBResponse{}, fmt.Errorf("invalid procedure type: %s", procedureType)
}

// StoreResponse stores a database response based on procedure type
func (s *DBStore) StoreResponse(procedureType string, procedureID int, responseDataTypeMarker *schema.ResponseDataType, response DBResponse) {
	if procedureType == "test" {
		s.StoreTestResponse(responseDataTypeMarker, procedureID, response)
	} else if procedureType == "generate" {
		s.StoreGenerateResponse(responseDataTypeMarker, procedureID, response)
	}
}

// ProcessAction processes a database action and stores the results
func (s *DBStore) ProcessAction(procedureType string, procedureID int, action schema.Action) error {
	if action.Type != "db:sql" && action.Type != "db:mongo" {
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}

	if action.Query == "" {
		return fmt.Errorf("query is required for database actions")
	}

	// Execute the query and get results
	results, err := s.db.Query(action.Query)
	if err != nil {
		return err
	}

	// Convert results to JSON
	resultBytes, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	response := DBResponse{
		Data:                   resultBytes,
		ResponseDataTypeMarker: &action.ResponseDataType,
	}

	s.StoreResponse(procedureType, procedureID, &action.ResponseDataType, response)
	return nil
}

// AccessFieldInResponseData accesses a specific field in the stored response data
func (s *DBStore) AccessFieldInResponseData(procedureType string, procedureID int, responseDataTypeMarkerID string, accessPaths []string) (interface{}, error) {
	response, err := s.GetResponse(procedureType, procedureID, responseDataTypeMarkerID)
	if err != nil {
		return nil, err
	}

	if response.ResponseDataTypeMarker.Type != "json" {
		return nil, fmt.Errorf("unsupported response data type: %s", response.ResponseDataTypeMarker.Type)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(response.Data, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
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
func (s *DBStore) ClearTestResponses() {
	s.testResponses = make(map[string]DBResponse)
}

// ClearGenerateResponses clears all stored generate responses
func (s *DBStore) ClearGenerateResponses() {
	s.generateResponses = make(map[string]DBResponse)
}

// ClearAllResponses clears all stored responses
func (s *DBStore) ClearAllResponses() {
	s.testResponses = make(map[string]DBResponse)
	s.generateResponses = make(map[string]DBResponse)
}

// Close closes the database connection
func (s *DBStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
