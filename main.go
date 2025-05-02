package main

import (
	"database/sql"
	"log"

	"github.com/mbaxamb3/nusli/api"
	db "github.com/mbaxamb3/nusli/db/sqlc"
	"github.com/mbaxamb3/nusli/util"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	// Make sure we have a DB driver set
	if config.DBDriver == "" {
		config.DBDriver = "postgres"
		log.Printf("DB driver not specified in config, using postgres")
	}

	log.Printf("Connecting to database with driver: %s", config.DBDriver)
	log.Printf("Database source: %s", config.DBSource)

	// Connect to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	// Verify database connection
	err = conn.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Printf("Successfully connected to database")

	// Create a new store
	store := db.NewStore(conn)

	// Create server
	server := api.NewServer(store)

	// Use a safe port if config doesn't specify one or uses port 80
	serverAddress := config.ServerAddress
	if serverAddress == "" || serverAddress == ":80" {
		serverAddress = ":8080"
		log.Printf("Using default port :8080 instead of port 80 which requires admin privileges")
	}

	// Start server
	log.Printf("Server starting on %s", serverAddress)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
