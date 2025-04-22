package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username:     "testuser1",
		Email:        "test1@example.com",
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
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	// First create a user
	createdUser := createRandomUser(t)

	// Then fetch the user
	fetchedUser, err := testQueries.GetUser(context.Background(), createdUser.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, createdUser.UserID, fetchedUser.UserID)
	require.Equal(t, createdUser.Username, fetchedUser.Username)
	require.Equal(t, createdUser.Email, fetchedUser.Email)
	require.Equal(t, createdUser.PasswordHash, fetchedUser.PasswordHash)
	require.Equal(t, createdUser.FirstName, fetchedUser.FirstName)
	require.Equal(t, createdUser.LastName, fetchedUser.LastName)

	// Check timestamps within a reasonable range
	require.WithinDuration(t, createdUser.CreatedAt.Time, fetchedUser.CreatedAt.Time, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	// First create a user
	createdUser := createRandomUser(t)

	// Then fetch the user by email
	fetchedUser, err := testQueries.GetUserByEmail(context.Background(), createdUser.Email)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, createdUser.UserID, fetchedUser.UserID)
	require.Equal(t, createdUser.Username, fetchedUser.Username)
	require.Equal(t, createdUser.Email, fetchedUser.Email)
}

func TestGetUserByUsername(t *testing.T) {
	// First create a user
	createdUser := createRandomUser(t)

	// Then fetch the user by username
	fetchedUser, err := testQueries.GetUserByUsername(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, createdUser.UserID, fetchedUser.UserID)
	require.Equal(t, createdUser.Username, fetchedUser.Username)
	require.Equal(t, createdUser.Email, fetchedUser.Email)
}

func TestUpdateUser(t *testing.T) {
	// First create a user
	createdUser := createRandomUser(t)

	// Define update parameters
	updateArg := UpdateUserParams{
		UserID:       createdUser.UserID,
		Username:     "updated-username",
		Email:        "updated@example.com",
		PasswordHash: "updated-hash",
		FirstName: sql.NullString{
			String: "Updated",
			Valid:  true,
		},
		LastName: sql.NullString{
			String: "Name",
			Valid:  true,
		},
	}

	// Update the user
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updateArg.UserID, updatedUser.UserID)
	require.Equal(t, updateArg.Username, updatedUser.Username)
	require.Equal(t, updateArg.Email, updatedUser.Email)
	require.Equal(t, updateArg.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, updateArg.FirstName, updatedUser.FirstName)
	require.Equal(t, updateArg.LastName, updatedUser.LastName)

	require.Equal(t, createdUser.CreatedAt, updatedUser.CreatedAt)
	require.True(t, updatedUser.UpdatedAt.Valid)
	require.NotEqual(t, createdUser.UpdatedAt, updatedUser.UpdatedAt)
}

func TestListUsers(t *testing.T) {
	// Create multiple users
	for i := 0; i < 5; i++ {
		createRandomUser(t)
	}

	// Test listing users with pagination
	arg := ListUsersParams{
		Limit:  3,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 3)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

	// Test next page
	arg = ListUsersParams{
		Limit:  3,
		Offset: 3,
	}

	users, err = testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
}

func TestDeleteUser(t *testing.T) {
	// Create a user
	user := createRandomUser(t)

	// Delete the user
	err := testQueries.DeleteUser(context.Background(), user.UserID)
	require.NoError(t, err)

	// Try to fetch the deleted user, should return error
	_, err = testQueries.GetUser(context.Background(), user.UserID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestCountUsers(t *testing.T) {
	// Create a user
	createRandomUser(t)

	// Count users
	count, err := testQueries.CountUsers(context.Background())
	require.NoError(t, err)
	require.True(t, count > 0)
}
