package tester

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Ishogbon/code-chaos/brokers"
	"github.com/Ishogbon/code-chaos/data"
	"github.com/Ishogbon/code-chaos/parser"
	"github.com/Ishogbon/code-chaos/schema"
)

// ExecuteAction executes an action and returns the result
func ExecuteAction(action *schema.Action, procedureType string, procedureID int) (string, error) {
	if action.SleepBefore > 0 {
		time.Sleep(time.Duration(action.SleepBefore) * time.Millisecond)
	}
	switch action.Type {
	case "endpoint":
		return executeEndpointAction(action, procedureType, procedureID)
	case "broker:rabbitmq":
		return executeRmqBrokerAction(action, procedureType, procedureID)
	default:
		return "", fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// executeEndpointAction executes an HTTP endpoint action
func executeEndpointAction(action *schema.Action, procedureType string, procedureID int) (string, error) {
	client := &http.Client{}

	// Create request based on method type
	var req *http.Request
	var err error

	url := data.ReplaceVariablesWithValuesInString(action.URL, procedureType, procedureID)

	switch action.MethodType {
	case "get":
		req, err = http.NewRequest("GET", url, nil)
	case "post":
		req, err = http.NewRequest("POST", url, nil)
	case "put":
		req, err = http.NewRequest("PUT", url, nil)
	case "delete":
		req, err = http.NewRequest("DELETE", url, nil)
	default:
		return "", fmt.Errorf("unsupported HTTP method: %s", action.MethodType)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Store the response
	httpResponse := data.HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}

	data.Store.StoreResponse(procedureType, procedureID, &action.ResponseDataType, httpResponse)

	return string(body), nil
}

// executeRmqBrokerAction executes a message broker action
func executeRmqBrokerAction(action *schema.Action, procedureType string, procedureID int) (string, error) {
	rmq_connection := brokers.BrokerConnectionManagerInstance.GetConnection(action.ConnectionID)
	message := data.ReplaceVariablesWithValuesInString(action.Message, procedureType, procedureID)
	err := rmq_connection.RMQ.Publish(action.Exchange, action.RoutingKey, []byte(message))
	if err != nil {
		panic(err)
	}
	return "", nil
}

func ExecuteResult(result *schema.ExpectedResult, procedureType string, procedureID int) (bool, error) {
	success := false
	for _, check := range result.Checks {
		success = false
		variableName, responseDataPath, err := parser.ParseExpectedResultCheck(check)
		if err != nil {
			return false, fmt.Errorf("failed to parse expected result check: %v", err)
		}
		variable := data.FetchVariable(variableName, procedureType, procedureID)
		variableToTestAgainst, err := data.Store.AccessFieldInResponseData(procedureType, procedureID, responseDataPath[0], responseDataPath[1:])
		if err != nil {
			return false, fmt.Errorf("failed to access field in response data: %v", err)
		}
		if result.Condition == "eq" {
			if variable.Value == variableToTestAgainst {
				success = true
			}
		} else if result.Condition == "neq" {
			if variable.Value != variableToTestAgainst {
				success = true
			}
		}
	}
	return success, nil
}
