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
