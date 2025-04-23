package tester

import (
	"github.com/Ishogbon/code-chaos/data"
	"github.com/Ishogbon/code-chaos/schema"
)

func Generate(test []schema.Generate) {
	for _, generate := range test {
		data.StoreVariables(generate.ID, generate.Variables, "generate")
	}
}
