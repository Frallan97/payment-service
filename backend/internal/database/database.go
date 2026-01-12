package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// Connect establishes a connection to PostgreSQL with retry logic
func Connect(databaseURL string) (*DB, error) {
	var db *sql.DB
	var err error

	// Retry connection with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		// Configure connection pool
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(30 * time.Minute)
		db.SetConnMaxIdleTime(5 * time.Minute)

		// Test connection
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return &DB{db}, nil
		}

		// Retry with exponential backoff
		waitTime := time.Duration(i+1) * time.Second
		log.Printf("Database connection attempt %d/%d failed: %v. Retrying in %v...", i+1, maxRetries, err, waitTime)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
