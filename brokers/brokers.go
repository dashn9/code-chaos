package brokers

import (
	"fmt"

	"github.com/Ishogbon/code-chaos/schema"
	"github.com/rabbitmq/amqp091-go"
)

type RMQ struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func (r *RMQ) Publish(exchange string, routingKey string, body []byte) error {
	return r.Channel.Publish(exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

type brokerConnection struct {
	ConnectionConfig *schema.Broker
	RMQ              RMQ
}

type BrokerConnectionManager struct {
	Connections map[string]brokerConnection
}

func (b *BrokerConnectionManager) GetConnection(connectionID string) brokerConnection {
	return b.Connections[connectionID]
}

func (b *BrokerConnectionManager) CreateRMQConnection(connectionConfig *schema.Broker) {
	conn, err := amqp091.Dial(connectionConfig.ConnectionURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel: %v", err))
	}

	b.Connections[connectionConfig.ConnectionID] = brokerConnection{
		ConnectionConfig: connectionConfig,
		RMQ:              RMQ{Connection: conn, Channel: ch},
	}
}

func (b *BrokerConnectionManager) CreateConnection(connectionConfig *schema.Broker) {
	if connectionConfig.Type == "rabbitmq" {
		b.CreateRMQConnection(connectionConfig)
	} else {
		panic(fmt.Sprintf("Unsupported broker type: %v", connectionConfig.Type))
	}
}

func (b *BrokerConnectionManager) CloseConnection(connectionID string) {
	conn := b.Connections[connectionID]
	conn.RMQ.Connection.Close()
	conn.RMQ.Channel.Close()
	delete(b.Connections, connectionID)
}

var BrokerConnectionManagerInstance BrokerConnectionManager

func CreateBrokerConnectionManagerInstance(connectionConfigs []schema.Broker) {
	BrokerConnectionManagerInstance = BrokerConnectionManager{
		Connections: make(map[string]brokerConnection),
	}
	for _, connectionConfig := range connectionConfigs {
		BrokerConnectionManagerInstance.CreateConnection(&connectionConfig)
	}
}
