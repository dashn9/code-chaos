package data

import (
	"fmt"

	"github.com/Ishogbon/code-chaos/schema"
)

// StoreWrapper provides a unified interface for accessing both HTTP and DB responses
type StoreWrapper struct {
	httpStore *HTTPStore
	dbStore   *DBStore
}

var Store = NewStoreWrapper(HTTP, DB)

// NewStoreWrapper creates a new store wrapper instance
func NewStoreWrapper(httpStore *HTTPStore, dbStore *DBStore) *StoreWrapper {
	return &StoreWrapper{
		httpStore: httpStore,
		dbStore:   dbStore,
	}
}

// Get retrieves a response from either HTTP or DB store based on the response data type
func (s *StoreWrapper) Get(procedureType string, procedureID int, responseDataTypeMarkerID string) (interface{}, error) {
	// Try HTTP store first
	httpResponse, httpErr := s.httpStore.GetResponse(procedureType, procedureID, responseDataTypeMarkerID)
	if httpErr == nil {
		return httpResponse, nil
	}

	// If HTTP store fails, try DB store
	dbResponse, dbErr := s.dbStore.GetResponse(procedureType, procedureID, responseDataTypeMarkerID)
	if dbErr == nil {
		return dbResponse, nil
	}

	// If both stores fail, return an error
	return nil, fmt.Errorf("response not found in either store for procedure type %s, ID %d, marker %s", procedureType, procedureID, responseDataTypeMarkerID)
}

// AccessFieldInResponseData accesses a specific field in the stored response data
func (s *StoreWrapper) AccessFieldInResponseData(procedureType string, procedureID int, responseDataTypeMarkerID string, accessPaths []string) (interface{}, error) {
	// Try HTTP store first
	_, httpErr := s.httpStore.GetResponse(procedureType, procedureID, responseDataTypeMarkerID)
	if httpErr == nil {
		return s.httpStore.AccessFieldInResponseData(procedureType, procedureID, responseDataTypeMarkerID, accessPaths)
	}

	// If HTTP store fails, try DB store
	_, dbErr := s.dbStore.GetResponse(procedureType, procedureID, responseDataTypeMarkerID)
	if dbErr == nil {
		return s.dbStore.AccessFieldInResponseData(procedureType, procedureID, responseDataTypeMarkerID, accessPaths)
	}

	// If both stores fail, return an error
	return nil, fmt.Errorf("response not found in either store for procedure type %s, ID %d, marker %s", procedureType, procedureID, responseDataTypeMarkerID)
}

// StoreResponse stores a response in the appropriate store based on its type
func (s *StoreWrapper) StoreResponse(procedureType string, procedureID int, responseDataTypeMarker *schema.ResponseDataType, response interface{}) {
	switch resp := response.(type) {
	case HTTPResponse:
		s.httpStore.StoreResponse(procedureType, procedureID, responseDataTypeMarker, resp)
	case DBResponse:
		s.dbStore.StoreResponse(procedureType, procedureID, responseDataTypeMarker, resp)
	default:
		panic(fmt.Sprintf("unsupported response type: %T", response))
	}
}

// ClearAllResponses clears all responses from both stores
func (s *StoreWrapper) ClearAllResponses() {
	s.httpStore.ClearAllResponses()
	s.dbStore.ClearAllResponses()
}
