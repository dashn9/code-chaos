package reader

import (
	"fmt"
	"os"

	"github.com/Ishogbon/code-chaos/schema"
	"gopkg.in/yaml.v2"
)

func LoadTestFromYAML(filePath string) (*schema.Test, error) {
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var test schema.Test
	err = yaml.Unmarshal(yamlData, &test)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &test, nil
}
