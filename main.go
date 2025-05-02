package main

import (
	"log"

	"github.com/mbaxamb3/nusli/api"
	db "github.com/mbaxamb3/nusli/db/sqlc"
	"github.com/mbaxamb3/nusli/util"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Set explicit database connection
	dbDriver := "postgres"
	dbSource := "postgresql://root:secret@localhost:5432/musli?sslmode=disable"

	log.Printf("DB driver: %s", dbDriver)
	log.Printf("Connecting to database with driver: %s", dbDriver)
	log.Printf("Database source: %s", dbSource)

	conn, err := util.ConnectDB(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(":8080")
	if err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
