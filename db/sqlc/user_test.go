package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserCreate tests the creation of a user in the database
func TestUserCreate(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Clear test data before starting
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%@testuser.com'")
	require.NoError(t, err)

	// Generate test data
	username := "testuser_" + randomString(8)
	email := username + "@testuser.com"
	passwordHash := "hashed_password_" + randomString(10)
	firstName := sql.NullString{String: "Test", Valid: true}
	lastName := sql.NullString{String: "User", Valid: true}

	// Execute
	query := `
	INSERT INTO users (username, email, password_hash, first_name, last_name)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING user_id, username, email, password_hash, first_name, last_name, created_at, updated_at
	`

	row := db.QueryRowContext(context.Background(), query, username, email, passwordHash, firstName, lastName)

	// Verify
	var user struct {
		UserID       int32
		Username     string
		Email        string
		PasswordHash string
		FirstName    sql.NullString
		LastName     sql.NullString
		CreatedAt    sql.NullTime
		UpdatedAt    sql.NullTime
	}

	err = row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	require.NoError(t, err)

	// Assert
	assert.NotZero(t, user.UserID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, passwordHash, user.PasswordHash)
	assert.Equal(t, firstName, user.FirstName)
	assert.Equal(t, lastName, user.LastName)
	assert.True(t, user.CreatedAt.Valid)
	assert.True(t, user.UpdatedAt.Valid)
	assert.NotZero(t, user.CreatedAt.Time)
	assert.NotZero(t, user.UpdatedAt.Time)

	// Log success
	t.Logf("Successfully created user: ID=%d, Username=%s", user.UserID, user.Username)
}

// TestUserQuery tests querying a user by email and username
func TestUserQuery(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Generate test data
	username := "queryuser_" + randomString(8)
	email := username + "@testuser.com"
	passwordHash := "hashed_password_" + randomString(10)
	firstName := sql.NullString{String: "Query", Valid: true}
	lastName := sql.NullString{String: "User", Valid: true}

	// Create a test user first
	var userID int32
	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (username, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5) RETURNING user_id`,
		username, email, passwordHash, firstName, lastName,
	).Scan(&userID)
	require.NoError(t, err)
	require.NotZero(t, userID)

	// Test cases for different queries
	testCases := []struct {
		name  string
		query string
		args  []interface{}
	}{
		{
			name:  "Query by ID",
			query: "SELECT * FROM users WHERE user_id = $1",
			args:  []interface{}{userID},
		},
		{
			name:  "Query by username",
			query: "SELECT * FROM users WHERE username = $1",
			args:  []interface{}{username},
		},
		{
			name:  "Query by email",
			query: "SELECT * FROM users WHERE email = $1",
			args:  []interface{}{email},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var user struct {
				UserID       int32
				Username     string
				Email        string
				PasswordHash string
				FirstName    sql.NullString
				LastName     sql.NullString
				CreatedAt    sql.NullTime
				UpdatedAt    sql.NullTime
			}

			err := db.QueryRowContext(context.Background(), tc.query, tc.args...).Scan(
				&user.UserID,
				&user.Username,
				&user.Email,
				&user.PasswordHash,
				&user.FirstName,
				&user.LastName,
				&user.CreatedAt,
				&user.UpdatedAt,
			)

			require.NoError(t, err)
			assert.Equal(t, userID, user.UserID)
			assert.Equal(t, username, user.Username)
			assert.Equal(t, email, user.Email)
			assert.Equal(t, passwordHash, user.PasswordHash)
			assert.Equal(t, firstName, user.FirstName)
			assert.Equal(t, lastName, user.LastName)
		})
	}

	// Clean up
	_, err = db.ExecContext(context.Background(), "DELETE FROM users WHERE user_id = $1", userID)
	require.NoError(t, err)
}

// TestUserUpdate tests updating a user's information
func TestUserUpdate(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Generate test data
	username := "updateuser_" + randomString(8)
	email := username + "@testuser.com"
	passwordHash := "hashed_password_" + randomString(10)
	firstName := sql.NullString{String: "Update", Valid: true}
	lastName := sql.NullString{String: "User", Valid: true}

	// Create a test user first
	var userID int32
	var originalUpdatedAt time.Time

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (username, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5) RETURNING user_id, updated_at`,
		username, email, passwordHash, firstName, lastName,
	).Scan(&userID, &originalUpdatedAt)
	require.NoError(t, err)
	require.NotZero(t, userID)

	// Wait a moment to ensure updated_at timestamp will be different
	time.Sleep(10 * time.Millisecond)

	// Update user
	newFirstName := sql.NullString{String: "UpdatedFirst", Valid: true}
	newLastName := sql.NullString{String: "UpdatedLast", Valid: true}
	newPasswordHash := "updated_hash_" + randomString(10)

	query := `
	UPDATE users
	SET 
		first_name = $1,
		last_name = $2,
		password_hash = $3,
		updated_at = CURRENT_TIMESTAMP
	WHERE user_id = $4
	RETURNING user_id, username, email, password_hash, first_name, last_name, created_at, updated_at
	`

	var updatedUser struct {
		UserID       int32
		Username     string
		Email        string
		PasswordHash string
		FirstName    sql.NullString
		LastName     sql.NullString
		CreatedAt    sql.NullTime
		UpdatedAt    sql.NullTime
	}

	err = db.QueryRowContext(
		context.Background(),
		query,
		newFirstName, newLastName, newPasswordHash, userID,
	).Scan(
		&updatedUser.UserID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.PasswordHash,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	require.NoError(t, err)

	// Verify fields were updated
	assert.Equal(t, userID, updatedUser.UserID)
	assert.Equal(t, username, updatedUser.Username) // Should not change
	assert.Equal(t, email, updatedUser.Email)       // Should not change
	assert.Equal(t, newPasswordHash, updatedUser.PasswordHash)
	assert.Equal(t, newFirstName, updatedUser.FirstName)
	assert.Equal(t, newLastName, updatedUser.LastName)
	assert.True(t, updatedUser.UpdatedAt.Time.After(originalUpdatedAt))

	// Clean up
	_, err = db.ExecContext(context.Background(), "DELETE FROM users WHERE user_id = $1", userID)
	require.NoError(t, err)
}

// TestUserList tests listing and pagination of users
func TestUserList(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Clear previous test users
	_, err = db.ExecContext(context.Background(), "DELETE FROM users WHERE username LIKE 'listuser_%'")
	require.NoError(t, err)

	// Insert multiple test users for listing
	userCount := 10
	userIDs := make([]int32, userCount)

	for i := 0; i < userCount; i++ {
		username := fmt.Sprintf("listuser_%d_%s", i, randomString(5))
		email := username + "@testuser.com"
		passwordHash := "password_" + randomString(8)

		err = db.QueryRowContext(
			context.Background(),
			`INSERT INTO users (username, email, password_hash)
			VALUES ($1, $2, $3) RETURNING user_id`,
			username, email, passwordHash,
		).Scan(&userIDs[i])
		require.NoError(t, err)
	}

	// Test pagination
	testCases := []struct {
		name     string
		limit    int
		offset   int
		expected int
	}{
		{
			name:     "First page (5 users)",
			limit:    5,
			offset:   0,
			expected: 5,
		},
		{
			name:     "Second page (5 users)",
			limit:    5,
			offset:   5,
			expected: 5,
		},
		{
			name:     "Partial page",
			limit:    15,
			offset:   7,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rows, err := db.QueryContext(
				context.Background(),
				"SELECT user_id, username, email FROM users WHERE username LIKE 'listuser_%' ORDER BY created_at DESC LIMIT $1 OFFSET $2",
				tc.limit, tc.offset,
			)
			require.NoError(t, err)
			defer rows.Close()

			// Count results
			var users []struct {
				UserID   int32
				Username string
				Email    string
			}

			for rows.Next() {
				var user struct {
					UserID   int32
					Username string
					Email    string
				}
				err := rows.Scan(&user.UserID, &user.Username, &user.Email)
				require.NoError(t, err)
				users = append(users, user)
			}
			require.NoError(t, rows.Err())

			// Check count matches expected
			assert.Equal(t, tc.expected, len(users))

			// Check all returned users have valid fields
			for _, user := range users {
				assert.NotZero(t, user.UserID)
				assert.Contains(t, user.Username, "listuser_")
				assert.Contains(t, user.Email, "@testuser.com")
			}
		})
	}

	// Clean up
	_, err = db.ExecContext(context.Background(), "DELETE FROM users WHERE username LIKE 'listuser_%'")
	require.NoError(t, err)
}

// Helper function to generate random strings
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
	}
	return string(b)
}

// You'll need to add this import if running this code
