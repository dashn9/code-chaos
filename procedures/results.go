package procedures

import (
	"fmt"

	"github.com/Ishogbon/code-chaos/schema"
)

type Results struct {
	results map[int]schema.ExpectedResult
}

var ResultStore = &Results{
	results: make(map[int]schema.ExpectedResult),
}

func (s *Results) Add(result schema.ExpectedResult) {
	// Check if a result with the given id already exists
	if _, exists := s.results[result.ID]; exists {
		panic(fmt.Sprintf("attempt to overwrite existing result with ID %d", result.ID))
	}

	// Store the new result
	s.results[result.ID] = result
}

func (s *Results) Get(id int) (*schema.ExpectedResult, error) {
	result, exists := s.results[id]
	if !exists {
		return nil, fmt.Errorf("result with ID %d not found", id)
	}
	return &result, nil
}

func StoreResult(result schema.ExpectedResult) {
	ResultStore.Add(result)
}

func StoreResults(results []schema.ExpectedResult) {
	for _, result := range results {
		StoreResult(result)
	}
}

func GetResult(id int) (*schema.ExpectedResult, error) {
	return ResultStore.Get(id)
}
