package tester

import (
	"strconv"

	"github.com/Ishogbon/code-chaos/data"
	"github.com/Ishogbon/code-chaos/procedures"
	"github.com/Ishogbon/code-chaos/schema"
)

func Generate(test []schema.Generate) {
	for _, generate := range test {
		for i := 0; i < generate.Count; i++ {
			data.StoreVariables(generate.ID, []string{"count(int):" + strconv.Itoa(i+1)}, "generate")
			data.StoreVariables(generate.ID, generate.Variables, "generate")
			action, err := procedures.GetAction(generate.Action)
			if err != nil {
				panic(err)
			}
			ExecuteAction(action, "generate", generate.ID)
		}
	}
}
