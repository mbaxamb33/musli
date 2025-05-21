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

// Convert different user struct types to the consistent userResponse format

// convertListUserRowToResponse converts a ListUsersRow to userResponse
func convertListUserRowToResponse(user db.ListUsersRow) userResponse {
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

// convertGetUserRowToResponse converts a GetUserByIDRow to userResponse
func convertGetUserRowToResponse(user db.GetUserByIDRow) userResponse {
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

// convertCreateUserRowToResponse converts a CreateUserRow to userResponse
func convertCreateUserRowToResponse(user db.CreateUserRow) userResponse {
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

// convertUpdateUserRowToResponse converts an UpdateUserPasswordRow to userResponse
func convertUpdateUserRowToResponse(user db.UpdateUserPasswordRow) userResponse {
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

// Keep the original function for backward compatibility if needed
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
	// Get authenticated user's cognito_sub from context
	_, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

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

	// Convert users to response format using the appropriate conversion function
	userResponses := make([]userResponse, len(users))
	for i, user := range users {
		userResponses[i] = convertListUserRowToResponse(user)
	}

	ctx.JSON(http.StatusOK, userResponses)
}

// getUserByID handles requests to get a specific user
func (server *Server) getUserByID(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	authedCognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

	// Only allow users to view their own information unless they're admins
	// TODO: Add admin check if admin functionality is needed
	if cognitoSub != authedCognitoSub.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this user's information"})
		return
	}

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

	// Return user response with the appropriate conversion function
	ctx.JSON(http.StatusOK, convertGetUserRowToResponse(user))
}

// createUser handles requests to create a new user
func (server *Server) createUser(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	_, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

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

	// Return created user with the appropriate conversion function
	ctx.JSON(http.StatusCreated, convertCreateUserRowToResponse(user))
}

// updateUser handles requests to update an existing user
func (server *Server) updateUser(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	authedCognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

	// Only allow users to update their own information unless they're admins
	// TODO: Add admin check if admin functionality is needed
	if cognitoSub != authedCognitoSub.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this user's information"})
		return
	}

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

	// Return updated user with the appropriate conversion function
	ctx.JSON(http.StatusOK, convertUpdateUserRowToResponse(user))
}

// deleteUser handles requests to delete a user
func (server *Server) deleteUser(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	authedCognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

	// Only allow users to delete their own account unless they're admins
	// TODO: Add admin check if admin functionality is needed
	if cognitoSub != authedCognitoSub.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this user"})
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

	// Delete user
	err = server.store.DeleteUser(ctx, cognitoSub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
