package tester

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Ishogbon/code-chaos/data"
	"github.com/Ishogbon/code-chaos/procedures"
	"github.com/Ishogbon/code-chaos/schema"
)

func Generate(test []schema.Generate) {
	for _, generate := range test {
		log.Printf("Procedure :: Generation :: %d, Count :: %d", generate.ID, generate.Count)
		for i := 0; i < generate.Count; i++ {
			log.Printf("Generating Resource ID: %d, Count: %d", generate.ID, i+1)
			data.StoreVariables(generate.ID, []string{"count(int):" + strconv.Itoa(i+1)}, "generate")
			data.StoreVariables(generate.ID, generate.Variables, "generate")
			for _, actionID := range generate.Actions {
				action, err := procedures.GetAction(actionID)
				if err != nil {
					panic(err)
				}
				_, err = ExecuteAction(action, "generate", generate.ID)
				if err != nil {
					panic(fmt.Sprintf("Error executing action %d: %v", actionID, err))
				}
			}
			for _, expectedResultID := range generate.ExpectedResults {
				expectedResult, err := procedures.GetResult(expectedResultID)
				if err != nil {
					panic(err)
				}
				result, err := ExecuteResult(expectedResult, "generate", generate.ID)
				if err != nil {
					panic(err)
				}
				if !result {
					panic("Expected result not met")
				} else {
					log.Printf("Expected result met for ID: %d", expectedResultID)
				}
			}
		}
	}
}
