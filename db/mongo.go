package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB(connectionString string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return &MongoDB{client: client}, nil
}

func (m *MongoDB) Open(connectionString string) error {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *MongoDB) Close() error {
	if m.client != nil {
		return m.client.Disconnect(context.Background())
	}
	return nil
}

func (m *MongoDB) Execute(query string, args ...interface{}) error {
	// For MongoDB, we'll treat this as a command
	// The query string should be a JSON string representing the command
	var command bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &command); err != nil {
		return fmt.Errorf("invalid MongoDB command: %v", err)
	}

	// Execute the command
	err := m.client.Database("admin").RunCommand(context.Background(), command).Err()
	if err != nil {
		return fmt.Errorf("failed to execute MongoDB command: %v", err)
	}

	return nil
}

func (m *MongoDB) Query(query string, args ...interface{}) (DBElements, error) {
	// For MongoDB, we'll treat this as a find operation
	// The query string should be a JSON string representing the filter
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, fmt.Errorf("invalid MongoDB filter: %v", err)
	}

	// Get the database and collection from the filter
	database, ok := filter["database"].(string)
	if !ok {
		return nil, fmt.Errorf("database name is required in MongoDB query")
	}
	delete(filter, "database")

	collection, ok := filter["collection"].(string)
	if !ok {
		return nil, fmt.Errorf("collection name is required in MongoDB query")
	}
	delete(filter, "collection")

	// Execute the find operation
	cursor, err := m.client.Database(database).Collection(collection).Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute MongoDB find: %v", err)
	}
	defer cursor.Close(context.Background())

	var results DBElements
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode MongoDB result: %v", err)
		}
		// Convert the MongoDB document to a DBElement
		element := make(DBElement)
		for k, v := range result {
			element[k] = v
		}
		results = append(results, element)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating MongoDB cursor: %v", err)
	}

	return results, nil
}
