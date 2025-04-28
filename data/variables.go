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

func (t *testVariables) addVariable(variable parser.Variable) {
	t.variables = append(t.variables, variable)
}

func (t *testVariables) addVariables(variables []parser.Variable) {
	t.variables = append(t.variables, variables...)
}

type generateVariables struct {
	variables []parser.Variable
	id        int
}

func (g *generateVariables) addVariable(variable parser.Variable) {
	g.variables = append(g.variables, variable)
}

func (g *generateVariables) addVariables(variables []parser.Variable) {
	g.variables = append(g.variables, variables...)
}

type Variables struct {
	testVariables     map[int]*testVariables
	generateVariables map[int]*generateVariables
}

func (v *Variables) GetGenerateVariable(generateVariablesId int, variableName string) (parser.Variable, error) {
	// Check if the map entry exists
	genVars, exists := v.generateVariables[generateVariablesId]
	if !exists || genVars == nil {
		return parser.Variable{}, errors.New("variable not found: " + variableName)
	}

	for _, variable := range genVars.variables {
		if variable.Name == variableName {
			return variable, nil
		}
	}
	return parser.Variable{}, errors.New("variable not found: " + variableName)
}

func (v *Variables) AddGenerateVariables(id int, variables []parser.Variable) {
	// Check if the map entry exists, if not create it
	if _, exists := v.generateVariables[id]; !exists {
		v.generateVariables[id] = &generateVariables{
			id:        id,
			variables: []parser.Variable{},
		}
	}

	v.generateVariables[id].addVariables(variables)
}

func (v *Variables) AddTestVariables(id int, variables []parser.Variable) {
	// Check if the map entry exists, if not create it
	if _, exists := v.testVariables[id]; !exists {
		v.testVariables[id] = &testVariables{
			id:        id,
			variables: []parser.Variable{},
		}
	}

	v.testVariables[id].addVariables(variables)
}

var StoredVariables = Variables{
	testVariables:     make(map[int]*testVariables),
	generateVariables: make(map[int]*generateVariables),
}

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
	var variables []*parser.Variable
	if procedureType == "generate" {
		for _, variableName := range variableNames {
			variable, err := StoredVariables.GetGenerateVariable(procedureID, variableName)
			if err != nil {
				panic(err)
			}
			_, err = ExecuteFunction(&variable, procedureType, procedureID)
			if err != nil {
				panic(err)
			}
			variables = append(variables, &variable)
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
