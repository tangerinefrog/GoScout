package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func Init() (*DB, error) {
	dsn := os.Getenv("DATABASE_CONN")
	if dsn == "" {
		return nil, errors.New("could not get db connection string from the .env file")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening a connection to a db: %w", err)
	}
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
