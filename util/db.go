package util

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// ConnectDB establishes a connection to the database
func ConnectDB(dbDriver, dbSource string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}

	return db, nil
}

// SetupDBConnection establishes a database connection using config
func SetupDBConnection(config Config) (*sql.DB, error) {
	db, err := ConnectDB(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, err
	}

	// Log successful connection
	log.Printf("Connected to database: %s", config.DBSource)

	return db, nil
}

// SetupTestDBConnection establishes a test database connection
func SetupTestDBConnection(config Config) (*sql.DB, error) {
	db, err := ConnectDB(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, err
	}

	// Log successful connection
	log.Printf("Connected to test database: %s", config.DBSource)

	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
			return
		}
		log.Println("Database connection closed")
	}
}
