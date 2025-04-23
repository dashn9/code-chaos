package main

import "github.com/Ishogbon/code-chaos/parser"

type testVariables struct {
	variables []parser.Variable
	actionId  int
}

type generateVariables struct {
	variables []parser.Variable
	actionId  int
}

type Variables struct {
	testVariables     []testVariables
	generateVariables []generateVariables
}

func (v *Variables) AddGenerateVariables(actionId int, variables []parser.Variable) {
	// Check if an generateVariable with the given actionId already exists
	for i, av := range v.testVariables {
		if av.actionId == actionId {
			// Extend the existing variables
			v.generateVariables[i].variables = append(v.generateVariables[i].variables, variables...)
			return
		}
	}

	// If no existing generateVariable found, create a new one
	v.generateVariables = append(v.generateVariables, generateVariables{
		actionId:  actionId,
		variables: variables,
	})
}

func (v *Variables) AddActionVariables(actionId int, variables []parser.Variable) {
	// Check if an actionVariable with the given actionId already exists
	for i, av := range v.testVariables {
		if av.actionId == actionId {
			// Extend the existing variables
			v.testVariables[i].variables = append(v.testVariables[i].variables, variables...)
			return
		}
	}

	// If no existing actionVariable found, create a new one
	v.testVariables = append(v.testVariables, testVariables{
		actionId:  actionId,
		variables: variables,
	})
}

var StoredVariables Variables

func parseVariables(input []string) {
	for _, variable := range input {
		v, err := parser.ParseVariable(variable)
		if err != nil {
			panic(err)
		}
		StoredVariables.AddActionVariables(0, []parser.Variable{*v})
	}
}
