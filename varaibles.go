package main

import "github.com/Ishogbon/code-chaos/parser"

type actionVariables struct {
	variables []parser.Variable
	actionId  int
}

type Variables struct {
	Variables []actionVariables
}

func (v *Variables) AddActionVariables(actionId int, variables []parser.Variable) {
	// Check if an actionVariable with the given actionId already exists
	for i, av := range v.Variables {
		if av.actionId == actionId {
			// Extend the existing variables
			v.Variables[i].variables = append(v.Variables[i].variables, variables...)
			return
		}
	}

	// If no existing actionVariable found, create a new one
	v.Variables = append(v.Variables, actionVariables{
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
