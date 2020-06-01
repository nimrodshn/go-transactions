package main

import (
	"database/sql"

	// Register the postgresql sql backend
	_ "github.com/lib/pq"
	"golang.org/x/net/context"

	//nolint
	errors "github.com/zgalor/weberr"
)

// db is the global database connection object.
var db *sql.DB

// By default do no roll back transaction only if it was set explicitly using rollback.Flag(ctx).
const defaultRollbackPolicy = false

// transaction represents an sql transaction
type transaction struct {
	// Reference to the database transaction:
	tx *sql.Tx

	// Flag indicating if the transaction should be rolled back:
	rollback bool

	// List of functions to call after committing the transaction:
	postCommitCallbacks []func()

	// List of functions to call after rolling back the transaction:
	postRollbackCallbacks []func()
}

// New constructs a new Transaction object.
func new() (*transaction, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &transaction{
		tx:       tx,
		rollback: defaultRollbackPolicy,
	}, nil
}

// InitDB opens a data base connection.
func InitDB(connStr string) (err error) {
	// InitDB should be called only once
	if db != nil {
		return errors.Errorf("DB object already exists")
	}

	// Open the postgres sql database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return
	}
	return
}

func CheckConnection() error {
	ctx, err := NewContext(context.Background())
	if err != nil {
		return err
	}
	defer Resolve(ctx)
	rows, err := db.QueryContext(ctx, "SELECT 1")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

// Close closes the connection to the DB.
func Close() error {
	return db.Close()
}
