package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// userResponse represents the API response structure for user data
type userResponse struct {
	CognitoSub string    `json:"cognito_sub"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// createUserRequest represents the request body for creating a user
type createUserRequest struct {
	CognitoSub string `json:"cognito_sub" binding:"required"`
	Username   string `json:"username" binding:"required,alphanum,min=3,max=30"`
	Password   string `json:"password" binding:"required,min=6"`
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// convertUserToResponse converts a database user model to an API response
func convertUserToResponse(user db.User) userResponse {
	createdAt := time.Time{}
	if user.CreatedAt.Valid {
		createdAt = user.CreatedAt.Time
	}

	return userResponse{
		CognitoSub: user.CognitoSub,
		Username:   user.Username,
		CreatedAt:  createdAt,
	}
}

// getUsers handles requests to get all users
func (server *Server) getUsers(ctx *gin.Context) {
	// Default pagination settings
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	limitParam := ctx.DefaultQuery("limit", "10")
	offsetParam := ctx.DefaultQuery("offset", "0")

	// Convert string params to integers
	var err error
	limit, err = strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err = strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get users from the database using the store
	users, err := server.store.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Convert users to response format
	userResponses := make([]userResponse, len(users))
	for i, user := range users {
		userResponses[i] = convertUserToResponse(user)
	}

	ctx.JSON(http.StatusOK, userResponses)
}

// getUserByID handles requests to get a specific user
func (server *Server) getUserByID(ctx *gin.Context) {
	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

	// No need to convert to integer, keep as string

	// Get user from database using the store
	user, err := server.store.GetUserByID(ctx, cognitoSub)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Return user response
	ctx.JSON(http.StatusOK, convertUserToResponse(user))
}

// createUser handles requests to create a new user
func (server *Server) createUser(ctx *gin.Context) {
	// Parse request body
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	_, err := server.store.GetUserByUsername(ctx, req.Username)
	if err == nil {
		// User already exists
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	} else if err != sql.ErrNoRows {
		// Database error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
		return
	}

	// Create user in database
	// Note: In a real application, you would hash the password here
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		CognitoSub: req.CognitoSub,
		Username:   req.Username,
		Password:   req.Password, // This should be hashed in a real application
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return created user
	ctx.JSON(http.StatusCreated, convertUserToResponse(user))
}

// updateUser handles requests to update an existing user
func (server *Server) updateUser(ctx *gin.Context) {
	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

	// No need to convert to integer, keep as string

	// Parse request body
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	_, err := server.store.GetUserByID(ctx, cognitoSub)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Update user password
	// Note: In a real application, you would hash the password here
	user, err := server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		CognitoSub: cognitoSub,
		Password:   req.Password, // This should be hashed in a real application
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return updated user
	ctx.JSON(http.StatusOK, convertUserToResponse(user))
}

// deleteUser handles requests to delete a user
func (server *Server) deleteUser(ctx *gin.Context) {
	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

	// No need to convert to integer, keep as string

	// Check if user exists
	_, err := server.store.GetUserByID(ctx, cognitoSub)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Delete user
	err = server.store.DeleteUser(ctx, cognitoSub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
