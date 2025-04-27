package tester

import (
	"github.com/Ishogbon/code-chaos/schema"
)

// Test executes a test with the given action and expected results
func Test(test *schema.Test, action *schema.Action, results []schema.ExpectedResult) error {
	// TODO: Implement test execution logic

	// For now, just return success
	return nil
}

// ValidateResult validates the result against the expected results
func ValidateResult(result string, expectedResults []schema.ExpectedResult) error {
	// TODO: Implement result validation logic
	return nil
}
