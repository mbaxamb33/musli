package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: randomString(6),
		Password: randomString(10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	// Create a random user first
	user1 := createRandomUser(t)

	// Retrieve the user by ID
	user2, err := testQueries.GetUserByID(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestGetUserByUsername(t *testing.T) {
	// Create a random user first
	user1 := createRandomUser(t)

	// Retrieve the user by username
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUserPassword(t *testing.T) {
	// Create a random user first
	user1 := createRandomUser(t)

	// Prepare new password
	newPassword := randomString(12)

	// Update user's password
	user2, err := testQueries.UpdateUserPassword(context.Background(), UpdateUserPasswordParams{
		UserID:   user1.UserID,
		Password: newPassword,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	// Verify updated password
	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, newPassword, user2.Password)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	// Create a random user first
	user1 := createRandomUser(t)

	// Delete the user
	err := testQueries.DeleteUser(context.Background(), user1.UserID)
	require.NoError(t, err)

	// Try to retrieve the deleted user (should fail)
	_, err = testQueries.GetUserByID(context.Background(), user1.UserID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestListUsers(t *testing.T) {
	// Create multiple users
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	// List users with pagination
	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestCreateUserWithDuplicateUsername(t *testing.T) {
	// Create an initial user
	user1 := createRandomUser(t)

	// Try to create another user with the same username
	arg := CreateUserParams{
		Username: user1.Username,
		Password: randomString(10),
	}

	_, err := testQueries.CreateUser(context.Background(), arg)
	require.Error(t, err) // Expect a unique constraint violation
}
