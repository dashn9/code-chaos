package db

import "database/sql"

type DBElement interface {
}

type DBElements interface {
}

type DBInterface interface {
	Open(connectionString string) (DB, error)
	Close() error
	Execute(query string, args ...interface{}) error
	Query(query string, args ...interface{}) (DBElements, error)
}

type DB struct {
	db *sql.DB
}

func (db *DB) Open(connectionString string) error {
	sqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	db.db = sqlDB
	return nil
}
