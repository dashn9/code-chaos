package tester

import (
	"fmt"

	"github.com/Ishogbon/code-chaos/brokers"
	"github.com/Ishogbon/code-chaos/data"
	"github.com/Ishogbon/code-chaos/schema"
)

// ExecuteAction executes an action and returns the result
func ExecuteAction(action *schema.Action, procedureType string, procedureID int) (string, error) {
	switch action.Type {
	case "endpoint":
		return executeEndpointAction(action)
	case "broker:rabbitmq":
		return executeRmqBrokerAction(action, procedureType, procedureID)
	default:
		return "", fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// executeEndpointAction executes an HTTP endpoint action
func executeEndpointAction(action *schema.Action) (string, error) {
	return "", nil
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
