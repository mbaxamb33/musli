package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/musli?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain is the entry point for testing
func TestMain(m *testing.M) {
	var err error

	// Connect to the database
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	// Create a new Queries object
	testQueries = New(testDB)

	// Run tests
	exitCode := m.Run()

	// Close the database connection
	testDB.Close()

	// Exit with the test result code
	os.Exit(exitCode)
}

// createRandomUser creates a random user for testing
func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
		PasswordHash: "password-hash",
		FirstName: sql.NullString{
			String: "Test",
			Valid:  true,
		},
		LastName: sql.NullString{
			String: "User",
			Valid:  true,
		},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create user: %v", err)
	}

	return user
}

// createRandomProject creates a random project for testing
func createRandomProject(t *testing.T, userID int32) Project {
	arg := CreateProjectParams{
		UserID: userID,
		Name:   fmt.Sprintf("Test Project %d", time.Now().UnixNano()),
		Description: sql.NullString{
			String: "Test project description",
			Valid:  true,
		},
		Status: sql.NullString{
			String: "active",
			Valid:  true,
		},
	}

	project, err := testQueries.CreateProject(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create project: %v", err)
	}

	return project
}

// createRandomCompany creates a random company for testing
func createRandomCompany(t *testing.T) Company {
	arg := CreateCompanyParams{
		Name: fmt.Sprintf("Test Company %d", time.Now().UnixNano()),
		Industry: sql.NullString{
			String: "Technology",
			Valid:  true,
		},
		Size: sql.NullString{
			String: "Medium",
			Valid:  true,
		},
		Location: sql.NullString{
			String: "Test Location",
			Valid:  true,
		},
		Website: sql.NullString{
			String: "https://example.com",
			Valid:  true,
		},
		Description: sql.NullString{
			String: "Test company description",
			Valid:  true,
		},
	}

	company, err := testQueries.CreateCompany(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create company: %v", err)
	}

	return company
}

// createRandomContact creates a random contact for testing
func createRandomContact(t *testing.T) Contact {
	arg := CreateContactParams{
		FirstName: "Test",
		LastName:  fmt.Sprintf("Contact %d", time.Now().UnixNano()),
		Title: sql.NullString{
			String: "Test Title",
			Valid:  true,
		},
		Email: sql.NullString{
			String: fmt.Sprintf("contact_%d@example.com", time.Now().UnixNano()),
			Valid:  true,
		},
		Phone: sql.NullString{
			String: "123-456-7890",
			Valid:  true,
		},
	}

	contact, err := testQueries.CreateContact(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create contact: %v", err)
	}

	return contact
}

// createRandomResourceCategory creates a random resource category for testing
func createRandomResourceCategory(t *testing.T) ResourceCategory {
	arg := CreateResourceCategoryParams{
		Name: fmt.Sprintf("Test Category %d", time.Now().UnixNano()),
		Description: sql.NullString{
			String: "Test category description",
			Valid:  true,
		},
	}

	category, err := testQueries.CreateResourceCategory(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create resource category: %v", err)
	}

	return category
}

// createRandomResource creates a random resource for testing
func createRandomResource(t *testing.T, categoryID sql.NullInt32) Resource {
	arg := CreateResourceParams{
		Name: fmt.Sprintf("Test Resource %d", time.Now().UnixNano()),
		Description: sql.NullString{
			String: "Test resource description",
			Valid:  true,
		},
		CategoryID: categoryID,
		Unit: sql.NullString{
			String: "pcs",
			Valid:  true,
		},
		CostPerUnit: sql.NullString{
			String: "10.00",
			Valid:  true,
		},
	}

	resource, err := testQueries.CreateResource(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create resource: %v", err)
	}

	return resource
}
