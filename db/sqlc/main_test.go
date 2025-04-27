package db

import (
	"database/sql"

	"log"

	"os"

	"testing"

	_ "github.com/lib/pq" // This is important - import the PostgreSQL driver
)

const (
	dbDriver = "postgres"

	dbSource = "postgresql://root:secret@localhost:5432/musli?sslmode=disable"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {

	var err error

	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {

		log.Fatalf("Cannot open database: %v", err)

	}

	// Test the connection with a ping

	err = testDB.Ping()

	if err != nil {

		log.Fatalf("Cannot ping database: %v", err)

	}

	testQueries = New(testDB)

	exitCode := m.Run()

	testDB.Close()

	os.Exit(exitCode)

}
