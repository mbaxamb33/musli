// middleware/auth.go
package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Global variables
var (
	jwksURL string
	jwks    map[string]*rsa.PublicKey
	jwksMu  sync.RWMutex
)

// Initialize JWKS
func InitJWKS(region, userPoolID string) {
	jwksURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
	fmt.Println("Setting JWKS URL to:", jwksURL)

	jwksMu.Lock()
	jwks = make(map[string]*rsa.PublicKey)
	jwksMu.Unlock()

	refreshJWKS() // Initial load

	// Auto-refresh JWKS every 10 minutes
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			refreshJWKS()
		}
	}()
}

// Fetch JWKS keys from Cognito
func refreshJWKS() {
	fmt.Println("Fetching JWKS keys from:", jwksURL)

	if jwksURL == "" {
		fmt.Println("JWKS URL is not initialized. Call InitJWKS first.")
		return
	}

	resp, err := http.Get(jwksURL)
	if err != nil {
		fmt.Println("Failed to fetch JWKS:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch JWKS. Status:", resp.Status)
		return
	}

	// Parse JSON response
	var jwksResponse struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
			Kty string `json:"kty"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwksResponse); err != nil {
		fmt.Println("Failed to parse JWKS:", err)
		return
	}

	if len(jwksResponse.Keys) == 0 {
		fmt.Println("Warning: JWKS response has no keys.")
		return
	}

	// Create new map to store parsed keys
	newJwks := make(map[string]*rsa.PublicKey)
	fmt.Println("JWKS keys retrieved:", len(jwksResponse.Keys))

	for _, key := range jwksResponse.Keys {
		fmt.Println("Available JWKS kid:", key.Kid)

		if key.Kty != "RSA" {
			fmt.Println("Skipping non-RSA key:", key.Kid)
			continue
		}

		// Decode modulus (N) from Base64 URL encoding
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			fmt.Println("Failed to decode modulus (N) for kid:", key.Kid, "Error:", err)
			continue
		}

		// Decode exponent (E) from Base64 URL encoding
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			fmt.Println("Failed to decode exponent (E) for kid:", key.Kid, "Error:", err)
			continue
		}

		// Convert exponent bytes to int
		var eInt int
		for i := 0; i < len(eBytes); i++ {
			eInt = eInt<<8 + int(eBytes[i])
		}

		// Create RSA Public Key
		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: eInt,
		}

		// Store the key in the map
		newJwks[key.Kid] = pubKey
		fmt.Printf("Added key %s to JWKS cache\n", key.Kid)
	}

	// Ensure we loaded at least one valid key
	if len(newJwks) == 0 {
		fmt.Println("No valid JWKS keys were loaded")
		return
	}

	// Replace old keys with new ones (thread-safe)
	jwksMu.Lock()
	jwks = newJwks
	jwksMu.Unlock()

	fmt.Println("JWKS keys successfully loaded:", len(jwks))
	for kid := range jwks {
		fmt.Println("Available JWKS kid:", kid)
	}
}

// Verify JWT token
func verifyAccessToken(accessToken string) (*jwt.Token, error) {
	fmt.Println("Verifying access token...")

	// Check if JWKS is initialized
	jwksMu.RLock()
	keysAvailable := len(jwks) > 0
	jwksMu.RUnlock()

	if !keysAvailable {
		fmt.Println("JWKS keys not loaded. Call InitJWKS and ensure it completes successfully.")
		return nil, errors.New("JWKS keys not loaded")
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Extract kid from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			fmt.Println("Missing 'kid' in token header")
			return nil, errors.New("missing kid in token header")
		}

		fmt.Println("Token 'kid':", kid)

		// Get public key for kid
		jwksMu.RLock()
		pubKey, exists := jwks[kid]
		jwksMu.RUnlock()

		if !exists {
			fmt.Println("Invalid 'kid', key not found:", kid)
			fmt.Println("Available JWKS kids:", getAvailableKids())
			return nil, errors.New("invalid kid, key not found")
		}

		fmt.Println("Found valid JWKS key for:", kid)
		return pubKey, nil
	})

	// Handle parsing errors
	if err != nil {
		fmt.Println("Token verification failed:", err)
		return nil, err
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Verify token has not expired
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			fmt.Println("Token has expired")
			return nil, errors.New("token has expired")
		}

		// Verify token issuer matches Cognito issuer
		issuer, ok := claims["iss"].(string)
		if !ok || !strings.Contains(issuer, "cognito-idp") {
			fmt.Println("Invalid token issuer:", issuer)
			return nil, errors.New("invalid token issuer")
		}
	}

	return token, nil
}

// Helper function to get available kids
func getAvailableKids() []string {
	jwksMu.RLock()
	defer jwksMu.RUnlock()

	keys := make([]string, 0, len(jwks))
	for k := range jwks {
		keys = append(keys, k)
	}
	return keys
}

// AuthMiddleware verifies JWT tokens and extracts user info
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// Hide full token in logs for security
		var logHeader string
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				logHeader = parts[0] + " [TOKEN_REDACTED]"
			} else {
				logHeader = authHeader
			}
		}
		fmt.Println("Received Authorization Header:", logHeader)

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			ctx.Abort()
			return
		}

		accessToken := parts[1]
		token, err := verifyAccessToken(accessToken)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		// Extract Cognito sub (user ID) and add to context
		cognitoSub, exists := claims["sub"].(string)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Cognito sub in token"})
			ctx.Abort()
			return
		}

		fmt.Println("Token verified successfully for sub:", cognitoSub)
		ctx.Set("cognito_sub", cognitoSub)
		ctx.Next()
	}
}
