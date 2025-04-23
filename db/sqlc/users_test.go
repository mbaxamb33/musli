package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomUser creates a user with random values for testing
func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:     randomString(8),
		Email:        randomEmail(),
		PasswordHash: randomString(16),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)

	return user
}

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestGetUserByID tests the GetUserByID function
func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByID(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

// TestGetUserByEmail tests the GetUserByEmail function
func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

// TestGetUserByUsername tests the GetUserByUsername function
func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

// TestUpdateUserEmail tests the UpdateUserEmail function
func TestUpdateUserEmail(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserEmailParams{
		UserID: user1.UserID,
		Email:  randomEmail(),
	}

	user2, err := testQueries.UpdateUserEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

// TestListUsers tests the ListUsers function
func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

	// Test pagination
	arg2 := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users2, err := testQueries.ListUsers(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, users2)
	require.Len(t, users2, 5)

	// Make sure the two sets of users are different
	for i := range users {
		require.NotEqual(t, users[i].UserID, users2[i].UserID)
	}
}

// TestDeleteUser tests the DeleteUser function
func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.UserID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user1.UserID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
