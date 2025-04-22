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
	v.Variables = append(v.Variables, actionVariables{
		actionId:  actionId,
		variables: variables,
	})
}

var variables Variables

func parseVariables(input []string) {
	for _, variable := range input {
		v, err := parser.ParseVariable(variable)
		if err != nil {
			panic(err)
		}
	}
}
