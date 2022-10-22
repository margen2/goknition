package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var StringConnectionBD string

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", StringConnectionBD)
	if err != nil {
		return nil, fmt.Errorf("connectdb/sql.open: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connectdb/db.ping: %w", err)
	}

	return db, nil
}

func SetConnection(connection string) {
	StringConnectionBD = connection
}
