package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// Server struct represents the API server
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func (server *Server) Start(address string) error {
	// Start the worker manager
	// ctx := context.Background()
	// You can add any background workers or tasks here

	// Start the HTTP server
	return server.router.Run(address)
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	// Set up Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User API routes
	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.GET("/", server.getUsers)
		userRoutes.GET("/:id", server.getUserByID)
		userRoutes.POST("/", server.createUser)
		userRoutes.PUT("/:id", server.updateUser)
		userRoutes.DELETE("/:id", server.deleteUser)
	}

	// Company API routes
	companyRoutes := router.Group("/api/v1/companies")
	{
		companyRoutes.GET("/", server.listCompanies)
		companyRoutes.GET("/:id", server.getCompanyByID)
		companyRoutes.POST("/", server.createCompany)
		companyRoutes.PUT("/:id", server.updateCompany)
		companyRoutes.DELETE("/:id", server.deleteCompany)

		// Get companies by user ID
		companyRoutes.GET("/user/:user_id", server.getCompaniesByUser)
	}

	// Contact API routes
	contactRoutes := router.Group("/api/v1/contacts")
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
	}

	// Add a route to test if the server is up
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Assign configured router to server
	server.router = router
	return server
}
