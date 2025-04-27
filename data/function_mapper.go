package data

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/Ishogbon/code-chaos/parser"
)

func parseInt(value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return parsed
}

func parseFuncArgValue(arg parser.FuncArg, procedureType string, procedureId int) string {
	if arg.Type == nil {
		return FetchVariable(arg.Value, procedureType, procedureId).Value
	}
	return arg.Value
}

func ExecuteFunction(function parser.Variable, procedureType string, procedureId int) (string, error) {
	if !function.IsFunc {
		return "", nil
	}
	if function.CacheResult && function.Value != "" {
		return function.Value, nil
	}
	switch function.Name {
	case "rng":
		minValue := parseInt(parseFuncArgValue(function.FuncArgs[0], procedureType, procedureId))
		maxValue := parseInt(parseFuncArgValue(function.FuncArgs[1], procedureType, procedureId))
		value := rand.IntN(maxValue-minValue+1) + minValue
		function.Value = strconv.Itoa(value)
		return function.Value, nil
	}
	return "", fmt.Errorf("function not found")
}
