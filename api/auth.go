// api/auth.go
package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
	"golang.org/x/oauth2"
)

const (
	clientID     = "6bo0q3c938g1oa0hjggqbdv0b"
	clientSecret = "rhlqr4vp21s2v5rfp7fijltlhe8ha3aj4i5oar561h8hgsvslam"
	redirectURL  = "http://localhost:8080/callback"
	issuerURL    = "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_177Be0rjJ"
	region       = "us-east-1"
)

var (
	provider      *oidc.Provider
	oauth2Config  oauth2.Config
	cognitoClient *cognitoidentityprovider.Client
)

// Initialize OAuth and Cognito
func InitAuth() {
	var err error

	// Load AWS SDK Config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Using hardcoded region
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region))
	if err != nil {
		fmt.Println("Failed to load AWS SDK config:", err)
		return
	}
	cognitoClient = cognitoidentityprovider.NewFromConfig(awsCfg)

	// Initialize OIDC provider
	providerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provider, err = oidc.NewProvider(providerCtx, issuerURL)
	if err != nil {
		fmt.Println("Failed to initialize OIDC provider:", err)
		return
	}

	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"openid", "email", "profile", "aws.cognito.signin.user.admin"},
	}

	fmt.Println("Auth initialization completed successfully")
}

// Generate Cognito secret hash
func generateSecretHash(username string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// SignUp request structure
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

// Handle user signup
func (server *Server) handleSignUp(ctx *gin.Context) {
	var req SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate secret hash
	secretHash := generateSecretHash(req.Email)

	// Register user in Cognito
	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientID),
		Username:   &req.Email,
		Password:   &req.Password,
		SecretHash: &secretHash,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
			{Name: aws.String("name"), Value: &req.FullName},
		},
	}

	resp, err := cognitoClient.SignUp(context.TODO(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed: " + err.Error()})
		return
	}

	// Get Cognito user sub
	cognitoSub := aws.ToString(resp.UserSub)

	// Store user in your database
	// Modify this to match your database structure
	_, dbErr := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:   req.Email,
		CognitoSub: cognitoSub,
	})

	if dbErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user: " + dbErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered. Check your email for the confirmation code."})
}

// Confirm signup request structure
type ConfirmSignUpRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// Handle signup confirmation
func (server *Server) handleConfirmSignUp(ctx *gin.Context) {
	var req ConfirmSignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(clientID),
		Username:         &req.Email,
		ConfirmationCode: &req.Code,
		SecretHash:       aws.String(generateSecretHash(req.Email)),
	}

	_, err := cognitoClient.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Account confirmation failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User confirmed successfully. You can now log in."})
}

// Handle login redirect
func (server *Server) handleLogin(ctx *gin.Context) {
	// Get the return URL from the query parameters, with a default value of "/"
	returnURL := ctx.DefaultQuery("returnUrl", "/")

	// Use the return URL as the state parameter to preserve it through the OAuth flow
	state := returnURL

	// Start OAuth flow with the state parameter
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusFound, url)
}

// Handle OAuth callback
func (server *Server) handleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	// Get the state parameter (which contains the return URL)
	state := ctx.Query("state")
	returnURL := "/"
	if state != "" {
		returnURL = state
	}

	// Exchange authorization code for tokens
	exchangeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := oauth2Config.Exchange(exchangeCtx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	// Extract tokens from response
	accessToken := token.AccessToken
	refreshToken, ok := token.Extra("refresh_token").(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Missing refresh token"})
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Missing ID token"})
		return
	}

	// Redirect to frontend with tokens and the return URL
	frontendURL := "http://localhost:5173/callback"
	redirectURL := fmt.Sprintf("%s?access_token=%s&id_token=%s&refresh_token=%s&return_url=%s",
		frontendURL,
		url.QueryEscape(accessToken),
		url.QueryEscape(rawIDToken),
		url.QueryEscape(refreshToken),
		url.QueryEscape(returnURL))

	ctx.Redirect(http.StatusFound, redirectURL)
}

// Handle token refresh
func (server *Server) handleRefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call Cognito to refresh tokens
	refreshCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := oauth2Config.TokenSource(refreshCtx, &oauth2.Token{RefreshToken: req.RefreshToken}).Token()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken,
		"id_token":     token.Extra("id_token"),
	})
}

// Handle logout
func (server *Server) handleLogout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
