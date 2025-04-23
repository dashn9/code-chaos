package reader

import (
	"fmt"
	"os"

	"github.com/Ishogbon/code-chaos/schema"
	"gopkg.in/yaml.v2"
)

// LoadTestFromYAML loads a test file from a YAML file
func LoadTestFromYAML(filePath string) (*schema.TestFile, error) {
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var testFile schema.TestFile
	err = yaml.Unmarshal(yamlData, &testFile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &testFile, nil
}

// LoadCodeChaosProceduresYAML is deprecated, use LoadTestFromYAML instead
func LoadCodeChaosProceduresYAML(filePath string) (*schema.Test, error) {
	testFile, err := LoadTestFromYAML(filePath)
	if err != nil {
		return nil, err
	}

	// For backward compatibility, return the first test
	if len(testFile.Tests) == 0 {
		return nil, fmt.Errorf("no tests found in the YAML file")
	}

	return &testFile.Tests[0], nil
}

// ProcessTest processes a test and returns the result
func ProcessTest(test *schema.Test) error {
	// TODO: Implement test processing logic
	return nil
}
