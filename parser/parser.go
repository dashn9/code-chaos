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

// FuncArg represents a function argument with its type and value
type FuncArg struct {
	Type  *string // Nullable type field to handle variable references
	Value string
}

// Variable represents a parsed expression that can be either a variable declaration
// or a function call. The structure holds all necessary information to process
// the expression further.
type Variable struct {
	// Name is the identifier for variables or function name
	Name string
	// Type represents the data type for variables or parameter types for functions
	Type string
	// Value holds the actual value for variables or the function's return value
	Value string
	// IsFunc indicates whether this is a function call (true) or variable declaration (false)
	IsFunc bool
	// FuncArgs holds the arguments for function calls with their types
	FuncArgs []FuncArg
	// CacheResult indicates whether the function result should be cached (true) or evaluated each time (false)
	CacheResult bool
}

// Regex patterns for parsing different expressions
var (
	// varPattern matches variable declarations: name(type):value
	// Example: x(int):42, price(float):19.99, name(string):John
	varPattern = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\(\s*(\w+)\s*\)\s*:\s*(.+)$`)

	// funcPattern matches function calls: value(type):func(value1(type1),value2(type2))
	// Example: random(string):rng(42(int),100(int)), 100.5(float):max(19.99(float),100.5(float))
	// With caching marker: random(string):rng!(42(int),100(int))
	funcPattern = regexp.MustCompile(`^([^(]+)\s*\(\s*(\w+)\s*\)\s*:\s*([a-zA-Z]{3,4})(!?)\s*\(\s*([^)]+)\s*\)$`)

	// argPattern matches function arguments: value(type) or variable reference
	// Example: 42(int), 19.99(float), "hello"(string), true(bool), myVar
	argPattern = regexp.MustCompile(`^([^(]+)(?:\s*\(\s*(\w+)\s*\))?$`)
)

// ParseVariable parses a string expression and returns a Variable struct.
// The input can be either:
//  1. Variable declaration: name(type):value
//     Example: x(int):42
//  2. Function call: value(type):func(value1(type1),value2(type2))
//     Example: random(string):rng(42(int),100(int))
//     With caching marker: random(string):rng!(42(int),100(int))
//     The '!' after the function name indicates the result should be cached
//
// Returns:
// - *Variable: The parsed expression
// - error: Any parsing errors that occurred
func ParseVariable(input string) (*Variable, error) {
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
func parseVariableMatch(matches []string) (*Variable, error) {
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

	return &Variable{
		Name:   name,
		Type:   typeStr,
		Value:  value,
		IsFunc: false,
	}, nil
}

// parseFunctionMatch processes a function call match.
// Expected format: value(type):func(value1(type1),value2(type2))
// Example: random(string):rng(42(int),100(int))
//
// The function validates:
// - Function name length (must be 3-4 letters)
// - Argument format (value(type) or variable reference value())
// - Argument types (must be valid types: int, float, string, bool)
// - Argument values (must be valid for their types)
// - Return value type (must be a valid type)
//
// Returns a Variable struct with:
// - Name: function name (3-4 letters)
// - Type: return value type
// - Value: function return value
// - IsFunc: true
// - FuncArgs: slice of FuncArg structs containing type and value information
func parseFunctionMatch(matches []string) (*Variable, error) {
	if len(matches) != 6 {
		return nil, errors.New("invalid function match")
	}

	value := strings.TrimSpace(matches[1])
	returnType := strings.TrimSpace(matches[2])
	name := matches[3]
	cacheMarker := matches[4]
	argsStr := matches[5]

	// Validate return type
	if !ValidTypes[returnType] {
		return nil, fmt.Errorf("unsupported return type: %s. Valid types are: int, float, string, bool", returnType)
	}

	// Validate return value
	if err := validateValue(value, returnType); err != nil {
		return nil, fmt.Errorf("invalid return value: %w", err)
	}

	// Validate function name length (3-4 letters)
	if len(name) < 3 || len(name) > 4 {
		return nil, fmt.Errorf("function name must be 3-4 letters long: %s", name)
	}

	// Split and validate arguments
	argStrings := strings.Split(argsStr, ",")
	funcArgs := make([]FuncArg, 0, len(argStrings))

	for _, argStr := range argStrings {
		argStr = strings.TrimSpace(argStr)
		argMatches := argPattern.FindStringSubmatch(argStr)
		if argMatches == nil {
			return nil, fmt.Errorf("invalid argument format: %s", argStr)
		}

		argValue := strings.TrimSpace(argMatches[1])

		// Check if this is a variable reference (no type specified)
		if len(argMatches) == 2 || argMatches[2] == "" {
			// This is a variable reference
			funcArgs = append(funcArgs, FuncArg{
				Type:  nil, // Type is nil for variable references
				Value: argValue,
			})
			continue
		}

		// This is a direct value with type
		argType := strings.TrimSpace(argMatches[2])

		// Validate the type
		if !ValidTypes[argType] {
			return nil, fmt.Errorf("unsupported argument type: %s. Valid types are: int, float, string, bool", argType)
		}

		// Validate the value against the type
		if err := validateValue(argValue, argType); err != nil {
			return nil, fmt.Errorf("invalid argument value: %w", err)
		}

		// Create a pointer to the type string
		typePtr := &argType
		funcArgs = append(funcArgs, FuncArg{
			Type:  typePtr,
			Value: argValue,
		})
	}

	return &Variable{
		Name:        name,
		Type:        returnType,
		Value:       value,
		IsFunc:      true,
		FuncArgs:    funcArgs,
		CacheResult: cacheMarker == "!",
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
