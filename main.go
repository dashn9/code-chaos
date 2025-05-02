package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ishogbon/code-chaos/brokers"
	"github.com/Ishogbon/code-chaos/db"
	"github.com/Ishogbon/code-chaos/procedures"
	"github.com/Ishogbon/code-chaos/reader"
	"github.com/Ishogbon/code-chaos/tester"
)

func main() {
	// Parse command line arguments
	yamlFile := flag.String("file", "", "Path to the YAML test file")
	flag.Parse()

	if *yamlFile == "" {
		fmt.Println("Please provide a YAML file path using the -file flag")
		flag.Usage()
		os.Exit(1)
	}

	// Load the test file
	testFile, err := reader.LoadTestFromYAML(*yamlFile)
	if err != nil {
		log.Fatalf("Failed to load test file: %v", err)
	}

	brokers.CreateBrokerConnectionManagerInstance(testFile.Globals.Brokers)
	db.CreateDBConnectionManagerInstance(testFile.Globals.DBs)

	procedures.StoreActions(testFile.Actions)
	procedures.StoreResults(testFile.Results)

	tester.Generate(testFile.Generates)

	// Process each test
	// for i, test := range testFile.Tests {
	// 	fmt.Printf("Running test %d (ID: %d)...\n", i+1, test.ID)

	// 	// Find the action for this test
	// 	var action *schema.Action
	// 	for _, a := range testFile.Actions {
	// 		if a.ID == test.Action {
	// 			action = &a
	// 			break
	// 		}
	// 	}

	// 	if action == nil {
	// 		log.Printf("Warning: Action with ID %d not found for test %d", test.Action, test.ID)
	// 		continue
	// 	}

	// 	// Process the test
	// 	err := tester.Test(&test, action, testFile.Results)
	// 	if err != nil {
	// 		log.Printf("Test %d failed: %v", test.ID, err)

	// 		if test.BehaviourOnFail == "stop" {
	// 			log.Fatalf("Test %d failed with 'stop' behavior, exiting", test.ID)
	// 		}
	// 	} else {
	// 		fmt.Printf("Test %d completed successfully\n", test.ID)
	// 	}
	// }

	// fmt.Println("All tests completed")
}
