package db

import (
	"fmt"
)

// DBElement represents a single database element (row/document)
type DBElement map[string]interface{}

// DBElements represents a collection of database elements
type DBElements []DBElement

// DB represents a database connection
type DB struct {
	client interface{} // This will hold either *sql.DB or *mongo.Client
	dbType string
}

// DBFactory creates database instances based on type
type DBFactory struct{}

// NewDBFactory creates a new database factory
func NewDBFactory() *DBFactory {
	return &DBFactory{}
}

// CreateDB creates a new database instance based on the type
func (f *DBFactory) CreateDB(dbType string, connectionString string) (*DB, error) {
	var client interface{}
	var err error

	switch dbType {
	case "mysql":
		client, err = NewMySQLDB(connectionString)
	case "mongodb":
		client, err = NewMongoDB(connectionString)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, err
	}

	return &DB{
		client: client,
		dbType: dbType,
	}, nil
}

// Query executes a database query and returns the results
func (d *DB) Query(query string, args ...interface{}) (DBElements, error) {
	if d.client == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	switch d.dbType {
	case "mysql":
		if mysqlDB, ok := d.client.(*MySQLDB); ok {
			return mysqlDB.Query(query, args...)
		}
	case "mongodb":
		if mongoDB, ok := d.client.(*MongoDB); ok {
			return mongoDB.Query(query, args...)
		}
	}

	return nil, fmt.Errorf("invalid database type or client")
}

// Execute executes a database command
func (d *DB) Execute(query string, args ...interface{}) error {
	if d.client == nil {
		return fmt.Errorf("database connection is nil")
	}

	switch d.dbType {
	case "mysql":
		if mysqlDB, ok := d.client.(*MySQLDB); ok {
			return mysqlDB.Execute(query, args...)
		}
	case "mongodb":
		if mongoDB, ok := d.client.(*MongoDB); ok {
			return mongoDB.Execute(query, args...)
		}
	}

	return fmt.Errorf("invalid database type or client")
}

// Close closes the database connection
func (d *DB) Close() error {
	if d.client == nil {
		return nil
	}

	switch d.dbType {
	case "mysql":
		if mysqlDB, ok := d.client.(*MySQLDB); ok {
			return mysqlDB.Close()
		}
	case "mongodb":
		if mongoDB, ok := d.client.(*MongoDB); ok {
			return mongoDB.Close()
		}
	}

	return fmt.Errorf("invalid database type or client")
}
