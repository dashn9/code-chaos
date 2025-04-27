package brokers

import (
	"github.com/Ishogbon/code-chaos/schema"
	"github.com/rabbitmq/amqp091-go"
)

type brokerConnection struct {
	ConnectionConfig *schema.Broker
	RMQConnection    *amqp091.Connection
	RMQChannel       *amqp091.Channel
}

type BrokerConnectionManager struct {
	Connections map[string]brokerConnection
}

func (b *BrokerConnectionManager) addConnection(connection brokerConnection) {
	b.Connections[connection.ConnectionID] = connection
}

func (b *BrokerConnectionManager) GetConnection(connectionID string) brokerConnection {
	return b.Connections[connectionID]
}

func (b *BrokerConnectionManager) CreateConnection(connectionID string, connectionType string, connectionURL string) {
	conn, err := amqp091.Dial(connectionURL)
	b.Connections[connectionID] = brokerConnection{
		ConnectionID:   connectionID,
		ConnectionType: connectionType,
		ConnectionURL:  connectionURL,
	}
}
