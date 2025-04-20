package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidTypes defines the allowed types in the system
var ValidTypes = map[string]bool{
	"int":    true,
	"float":  true,
	"string": true,
	"bool":   true,
}

// TypeValue represents a parsed expression that can be either a variable declaration
// or a function call. The structure holds all necessary information to process
// the expression further.
type TypeValue struct {
	// Name is the identifier for variables or function name
	Name string
	// Type represents the data type for variables or parameter types for functions
	Type string
	// Value holds the actual value for variables or the function's return value
	Value string
	// IsFunc indicates whether this is a function call (true) or variable declaration (false)
	IsFunc bool
	// FuncArgs holds the arguments for function calls
	FuncArgs []string
}

// Regex patterns for parsing different expressions
var (
	// varPattern matches variable declarations: name(type):value
	// Example: x(int):42, price(float):19.99, name(string):John
	varPattern = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\(\s*(\w+)\s*\)\s*:\s*(.+)$`)

	// funcPattern matches function calls: func(type1,type2):value
	// Example: rng(int,int):random, max(float,float):100.5
	funcPattern = regexp.MustCompile(`^([a-zA-Z]{3,4})\s*\(\s*([^)]+)\s*\)\s*:\s*(.+)$`)
)

// ParseTypeValue parses a string expression and returns a TypeValue struct.
// The input can be either:
//  1. Variable declaration: name(type):value
//     Example: x(int):42
//  2. Function call: func(type1,type2):value
//     Example: rng(int,int):random
//
// Returns:
// - *TypeValue: The parsed expression
// - error: Any parsing errors that occurred
func ParseTypeValue(input string) (*TypeValue, error) {
	if input == "" {
		return nil, errors.New("input is empty")
	}

	// Try to match function pattern first (more specific)
	if matches := funcPattern.FindStringSubmatch(input); matches != nil {
		return parseFunctionMatch(matches)
	}

	// Try to match variable pattern
	if matches := varPattern.FindStringSubmatch(input); matches != nil {
		return parseVariableMatch(matches)
	}

	return nil, errors.New("invalid format: does not match expected patterns")
}

// parseVariableMatch processes a variable declaration match.
// Expected format: name(type):value
// Example: x(int):42
func parseVariableMatch(matches []string) (*TypeValue, error) {
	if len(matches) != 4 {
		return nil, errors.New("invalid variable match")
	}

	name := matches[1]
	typeStr := matches[2]
	value := matches[3]

	// Validate the type
	if !ValidTypes[typeStr] {
		return nil, fmt.Errorf("unsupported type: %s. Valid types are: int, float, string, bool", typeStr)
	}

	// Validate the value against the type
	if err := validateValue(value, typeStr); err != nil {
		return nil, err
	}

	return &TypeValue{
		Name:   name,
		Type:   typeStr,
		Value:  value,
		IsFunc: false,
	}, nil
}

// parseFunctionMatch processes a function call match.
// Expected format: func(type1,type2):value
// Example: rng(int,int):random
func parseFunctionMatch(matches []string) (*TypeValue, error) {
	if len(matches) != 4 {
		return nil, errors.New("invalid function match")
	}

	name := matches[1]
	typeStr := matches[2]
	value := matches[3]

	// Validate function name length (3-4 letters)
	if len(name) < 3 || len(name) > 4 {
		return nil, fmt.Errorf("function name must be 3-4 letters long: %s", name)
	}

	// Split and validate argument types
	argTypes := strings.Split(typeStr, ",")
	for _, argType := range argTypes {
		argType = strings.TrimSpace(argType)
		if !ValidTypes[argType] {
			return nil, fmt.Errorf("unsupported argument type: %s. Valid types are: int, float, string, bool", argType)
		}
	}

	return &TypeValue{
		Name:     name,
		Type:     typeStr,
		Value:    value,
		IsFunc:   true,
		FuncArgs: argTypes,
	}, nil
}

// validateValue checks if a value is valid for a given type
func validateValue(value, typeStr string) error {
	switch typeStr {
	case "int":
		_, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
	case "float":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", value)
		}
	case "bool":
		_, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "string":
		// Strings are always valid
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", typeStr)
	}
	return nil
}
