package data

import "github.com/Ishogbon/code-chaos/parser"

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
