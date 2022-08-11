package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/margen2/goknition/config"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.StringConnectionBD)
	if err != nil {
		return nil, fmt.Errorf("connectdb/sql.open: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connectdb/db.ping: %w", err)
	}

	return db, nil
}
