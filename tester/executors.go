package tester

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ishogbon/code-chaos/brokers"
	"github.com/Ishogbon/code-chaos/data"
	"github.com/Ishogbon/code-chaos/db"
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
	case "db:sql":
		return executeSQLDBAction(action, procedureType, procedureID)
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

// executeDBAction executes a database action
func executeSQLDBAction(action *schema.Action, procedureType string, procedureID int) (string, error) {
	query := action.Query
	if query == "" {
		return "", fmt.Errorf("query is required for database actions")
	}

	query = data.ReplaceVariablesWithValuesInString(query, procedureType, procedureID)
	// Get the database connection using the connection manager
	dbConn := db.DBConnectionManagerInstance.GetConnection(action.ConnectionID)
	if dbConn.DB == nil {
		return "", fmt.Errorf("database connection is not initialized for connection ID: %s", action.ConnectionID)
	}

	// Execute the query
	results, err := dbConn.DB.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	// Store the response
	dbResponse := data.DBResponse{
		Elements:               results,
		ResponseDataTypeMarker: &action.ResponseDataType,
	}

	data.DB.StoreResponse(procedureType, procedureID, &action.ResponseDataType, dbResponse)

	return "", nil
}

// compareValues performs a comparison between two values based on the given condition
func compareValues(condition string, value1, value2 interface{}) bool {
	switch condition {
	case "eq":
		return value1 == value2
	case "neq":
		return value1 != value2
	case "gt":
		if v1, ok := value1.(int); ok {
			if v2, ok := value2.(int); ok {
				return v1 > v2
			}
		}
		return false
	case "lt":
		if v1, ok := value1.(int); ok {
			if v2, ok := value2.(int); ok {
				return v1 < v2
			}
		}
		return false
	case "gte":
		if v1, ok := value1.(int); ok {
			if v2, ok := value2.(int); ok {
				return v1 >= v2
			}
		}
		return false
	case "lte":
		if v1, ok := value1.(int); ok {
			if v2, ok := value2.(int); ok {
				return v1 <= v2
			}
		}
		return false
	default:
		return false
	}
}

// parseArrayLengthCheck parses an array length check string and returns the condition and quantity
func parseArrayLengthCheck(value string) (string, int, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 4 || parts[0] != "arr" || parts[1] != "len" {
		return "", 0, fmt.Errorf("invalid array length check format: %s", value)
	}

	quantity, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", 0, fmt.Errorf("invalid array length quantity: %s", parts[3])
	}

	return parts[2], quantity, nil
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

		// Handle array length checks
		if strings.HasPrefix(result.Condition, "arr:len:") {
			condition, expectedLen, err := parseArrayLengthCheck(result.Condition)
			if err != nil {
				return false, err
			}

			// Check if the value is an array/slice
			actualArray, ok := variableToTestAgainst.([]interface{})
			if !ok {
				return false, fmt.Errorf("value is not an array: %v", variableToTestAgainst)
			}

			// Check array length
			actualLen := len(actualArray)
			success = compareValues(condition, actualLen, expectedLen)
		} else {
			// Handle regular value comparisons
			success = compareValues(result.Condition, variable.Value, variableToTestAgainst)
		}
	}
	return success, nil
}
