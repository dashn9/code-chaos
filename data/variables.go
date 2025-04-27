package data

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Ishogbon/code-chaos/parser"
)

type testVariables struct {
	variables []parser.Variable
	id        int
}

type generateVariables struct {
	variables []parser.Variable
	id        int
}

type Variables struct {
	testVariables     []testVariables
	generateVariables []generateVariables
}

func (v *Variables) GetGenerateVariable(generateVariablesId int, variableName string) (parser.Variable, error) {
	for _, variable := range v.generateVariables[generateVariablesId].variables {
		if variable.Name == variableName {
			return variable, nil
		}
	}
	return parser.Variable{}, errors.New("variable not found")
}
func (v *Variables) AddGenerateVariables(id int, variables []parser.Variable) {
	// Check if a generateVariable with the given id already exists
	for i, av := range v.generateVariables {
		if av.id == id {
			// Extend the existing variables
			v.generateVariables[i].variables = append(v.generateVariables[i].variables, variables...)
			return
		}
	}

	// If no existing generateVariable found, create a new one
	v.generateVariables = append(v.generateVariables, generateVariables{
		id:        id,
		variables: variables,
	})
}

func (v *Variables) AddTestVariables(id int, variables []parser.Variable) {
	// Check if an actionVariable with the given id already exists
	for i, av := range v.testVariables {
		if av.id == id {
			// Extend the existing variables
			v.testVariables[i].variables = append(v.testVariables[i].variables, variables...)
			return
		}
	}

	// If no existing actionVariable found, create a new one
	v.testVariables = append(v.testVariables, testVariables{
		id:        id,
		variables: variables,
	})
}

var StoredVariables Variables

func StoreVariables(id int, input []string, variableType string) {
	for _, variable := range input {
		v, err := parser.ParseVariable(variable)
		if err != nil {
			panic(err)
		}
		if variableType == "test" {
			StoredVariables.AddTestVariables(id, []parser.Variable{*v})
		} else if variableType == "generate" {
			StoredVariables.AddGenerateVariables(id, []parser.Variable{*v})
		} else {
			panic("Invalid variable type")
		}
	}
}

func FindVariablesInString(input string) []string {
	// Regular expression to match $ followed by an alphabetic character and optional alphanumeric characters
	re := regexp.MustCompile(`\$[a-zA-Z][a-zA-Z0-9]*`)
	matches := re.FindAllString(input, -1)

	// Strip the $ prefix from each match
	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match[1:] // Remove the first character ($)
	}
	return result
}

func FetchVariable(variableName string, procedureType string, procedureID int) parser.Variable {
	if procedureType == "generate" {
		variable, err := StoredVariables.GetGenerateVariable(procedureID, variableName)
		if err != nil {
			panic(err)
		}
		return variable
	}
	// else if procedureType == "test" {
	// 	return StoredVariables.GetTestVariable(procedureID, variableName)
	// }
	return parser.Variable{}
}

func ReplaceVariablesWithValuesInString(input string, procedureType string, procedureID int) string {
	variableNames := FindVariablesInString(input)
	var variables []parser.Variable
	if procedureType == "generate" {
		for _, variableName := range variableNames {
			variable, err := StoredVariables.GetGenerateVariable(procedureID, variableName)
			if err != nil {
				panic(err)
			}
			_, err = ExecuteFunction(variable, procedureType, procedureID)
			if err != nil {
				panic(err)
			}
			variables = append(variables, variable)
		}
	} else if procedureType == "test" {
	}
	if len(variables) != len(variableNames) {
		panic("Some variables you defined are missing in the procedure's scope")
	}

	// Replace each variable with its value
	result := input
	for i, variable := range variables {
		// Add $ back to the variable name for replacement
		searchStr := "$" + variableNames[i]
		result = strings.ReplaceAll(result, searchStr, variable.Value)
	}
	return result
}
