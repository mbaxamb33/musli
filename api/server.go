package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
	"github.com/mbaxamb3/nusli/middleware"
	"golang.org/x/oauth2"
)

// Hardcoded AWS Cognito configuration
const (
	cognitoRegion       = "us-east-1"
	cognitoUserPoolID   = "us-east-1_177Be0rjJ"
	cognitoClientID     = "6bo0q3c938g1oa0hjggqbdv0b"
	cognitoClientSecret = "rhlqr4vp21s2v5rfp7fijltlhe8ha3aj4i5oar561h8hgsvslam"
	cognitoRedirectURL  = "http://localhost:8080/callback"
	cognitoIssuerURL    = "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_177Be0rjJ"
)

// Server struct represents the API server
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func (server *Server) Start(address string) error {
	// Start the HTTP server
	return server.router.Run(address)
}

// initializeAuth initializes the authentication components
func initializeAuth() {
	// Load AWS SDK Config with region
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cognitoRegion))
	if err != nil {
		fmt.Println("Failed to load AWS SDK config:", err)
		return
	}
	cognitoClient = cognitoidentityprovider.NewFromConfig(awsCfg)

	// Initialize OIDC provider
	providerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provider, err = oidc.NewProvider(providerCtx, cognitoIssuerURL)
	if err != nil {
		fmt.Println("Failed to initialize OIDC provider:", err)
		return
	}

	oauth2Config = oauth2.Config{
		ClientID:     cognitoClientID,
		ClientSecret: cognitoClientSecret,
		RedirectURL:  cognitoRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"openid", "email", "profile", "aws.cognito.signin.user.admin"},
	}

	fmt.Println("Auth initialization completed successfully")
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	// Initialize authentication systems with hardcoded values
	initializeAuth()

	// Initialize JWKS with hardcoded Cognito user pool region and ID
	middleware.InitJWKS(cognitoRegion, cognitoUserPoolID)

	// Set up Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // For development only - more permissive
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add this log
	fmt.Println("CORS middleware configured with AllowAllOrigins: true")

	// Public Authentication Routes - accessible without authentication
	router.POST("/signup", server.handleSignUp)
	router.POST("/confirm-signup", server.handleConfirmSignUp)
	router.GET("/login", server.handleLogin)
	router.GET("/callback", server.handleCallback)
	router.POST("/refresh-token", server.handleRefreshToken)
	router.GET("/logout", server.handleLogout)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Protected API routes - require authentication
	apiRoutes := router.Group("/api/v1")
	apiRoutes.Use(middleware.AuthMiddleware()) // Apply auth middleware to all /api/v1 routes

	// User API routes - Updated to use cognito_sub instead of id
	userRoutes := apiRoutes.Group("/users")
	{
		userRoutes.GET("/", server.getUsers)
		userRoutes.GET("/:cognito_sub", server.getUserByID) // Changed from :id to :cognito_sub
		userRoutes.POST("/", server.createUser)
		userRoutes.PUT("/:cognito_sub", server.updateUser)    // Changed from :id to :cognito_sub
		userRoutes.DELETE("/:cognito_sub", server.deleteUser) // Changed from :id to :cognito_sub
		userRoutes.GET("/me", server.getCurrentUser)          // Enabled endpoint to get current user
	}

	// Company API routes
	companyRoutes := apiRoutes.Group("/companies")
	{
		companyRoutes.GET("/", server.listCompanies)
		companyRoutes.GET("/:id", server.getCompanyByID)
		companyRoutes.POST("/", server.createCompany)
		companyRoutes.PUT("/:id", server.updateCompany)
		companyRoutes.DELETE("/:id", server.deleteCompany)

		// Get companies by user cognito_sub - Changed from user_id to cognito_sub
		companyRoutes.GET("/user/:cognito_sub", server.getCompaniesByUser) // Changed from :user_id to :cognito_sub

		// Company datasources routes
		companyRoutes.GET("/:id/datasources", server.listCompanyDatasources)
		companyRoutes.POST("/:id/datasources", server.createCompanyDatasource)
		companyRoutes.DELETE("/:id/datasources/:datasource_id", server.deleteCompanyDatasource)

		// Company paragraphs routes
		companyRoutes.GET("/:id/paragraphs", server.listCompanyParagraphs)
		companyRoutes.GET("/:id/paragraphs/search", server.searchCompanyParagraphs)
	}

	// Contact API routes
	contactRoutes := apiRoutes.Group("/contacts")
	{
		contactRoutes.GET("/", server.listContacts)
		contactRoutes.GET("/:id", server.getContactByID)
		contactRoutes.POST("/", server.createContact)
		contactRoutes.PUT("/:id", server.updateContact)
		contactRoutes.DELETE("/:id", server.deleteContact)

		// Get contacts by company ID
		contactRoutes.GET("/company/:company_id", server.listContactsByCompany)

		// Search contacts by name
		contactRoutes.GET("/search", server.searchContactsByName)

		// Contact datasources routes
		contactRoutes.GET("/:id/datasources", server.listContactDatasources)
		contactRoutes.POST("/:id/datasources", server.createContactDatasource)
		contactRoutes.DELETE("/:id/datasources/:datasource_id", server.deleteContactDatasource)

		// Contact paragraphs routes
		contactRoutes.GET("/:id/paragraphs", server.listContactParagraphs)
		contactRoutes.GET("/:id/paragraphs/search", server.searchContactParagraphs)
	}

	// Shared file upload endpoint for both companies and contacts
	apiRoutes.POST("/:entity_type/:id/datasources/upload", server.uploadDatasource)

	// Paragraphs API routes
	paragraphRoutes := apiRoutes.Group("/paragraphs")
	{
		paragraphRoutes.GET("/:id", server.getParagraphByID)
		paragraphRoutes.POST("/", server.createParagraph)
		paragraphRoutes.PUT("/:id", server.updateParagraph)
		paragraphRoutes.DELETE("/:id", server.deleteParagraph)

		// Get paragraphs by datasource ID
		paragraphRoutes.GET("/datasource/:datasource_id", server.listParagraphsByDatasource)
	}

	// Project API routes
	projectRoutes := apiRoutes.Group("/projects")
	{
		projectRoutes.GET("/", server.listProjects) // This will use cognito_sub from authentication
		projectRoutes.GET("/:id", server.getProjectByID)
		projectRoutes.POST("/", server.createProject) // This will use cognito_sub from authentication
		projectRoutes.DELETE("/:id", server.deleteProject)
		// Project datasources routes
		projectRoutes.GET("/:id/datasources", server.listDatasourcesByProject)
		projectRoutes.POST("/:id/datasources", server.createAndAssociateProjectDatasource)
		projectRoutes.POST("/:id/datasources/associate", server.associateDatasourceWithProject)
		projectRoutes.DELETE("/:id/datasources/:datasource_id", server.removeDatasourceFromProject)
	}

	// Datasource processing route
	apiRoutes.POST("/datasources/:id/process", server.processDatasourceByID)

	// Assign configured router to server
	server.router = router
	return server
}

// getCurrentUser returns the authenticated user's details
func (server *Server) getCurrentUser(ctx *gin.Context) {
	// Get the cognito_sub from the context (added by middleware)
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Fetch user information from the database
	user, err := server.store.GetUserByID(ctx, cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// Return user data using the appropriate conversion function
	ctx.JSON(http.StatusOK, convertGetUserRowToResponse(user))
}
