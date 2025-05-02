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

// Update the NewServer function in api/server.go

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

		// Company datasources routes
		companyRoutes.GET("/:id/datasources", server.listCompanyDatasources)
		companyRoutes.POST("/:id/datasources", server.createCompanyDatasource)
		companyRoutes.DELETE("/:id/datasources/:datasource_id", server.deleteCompanyDatasource)

		// Company paragraphs routes
		companyRoutes.GET("/:id/paragraphs", server.listCompanyParagraphs)
		companyRoutes.GET("/:id/paragraphs/search", server.searchCompanyParagraphs)
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

		// Contact datasources routes
		contactRoutes.GET("/:id/datasources", server.listContactDatasources)
		contactRoutes.POST("/:id/datasources", server.createContactDatasource)
		contactRoutes.DELETE("/:id/datasources/:datasource_id", server.deleteContactDatasource)

		// Contact paragraphs routes
		contactRoutes.GET("/:id/paragraphs", server.listContactParagraphs)
		contactRoutes.GET("/:id/paragraphs/search", server.searchContactParagraphs)
	}

	// Shared file upload endpoint for both companies and contacts
	router.POST("/api/v1/:entity_type/:id/datasources/upload", server.uploadDatasource)

	// Paragraphs API routes
	paragraphRoutes := router.Group("/api/v1/paragraphs")
	{
		paragraphRoutes.GET("/:id", server.getParagraphByID)
		paragraphRoutes.POST("/", server.createParagraph)
		paragraphRoutes.PUT("/:id", server.updateParagraph)
		paragraphRoutes.DELETE("/:id", server.deleteParagraph)

		// Get paragraphs by datasource ID
		paragraphRoutes.GET("/datasource/:datasource_id", server.listParagraphsByDatasource)
	}

	// Add these new routes to the NewServer function in api/server.go
	// Add below the existing route groups

	// Project API routes
	projectRoutes := router.Group("/api/v1/projects")
	{
		projectRoutes.GET("/", server.listProjects)
		projectRoutes.GET("/:id", server.getProjectByID)
		projectRoutes.POST("/", server.createProject)
		// Project datasources routes
		projectRoutes.GET("/:project_id/datasources", server.listDatasourcesByProject)
		projectRoutes.POST("/:project_id/datasources", server.createAndAssociateProjectDatasource)
		projectRoutes.POST("/:project_id/datasources/associate", server.associateDatasourceWithProject)
		projectRoutes.DELETE("/:project_id/datasources/:datasource_id", server.removeDatasourceFromProject)
	}

	// Add a route to test if the server is up
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.POST("/api/v1/datasources/:id/process", server.processDatasourceByID)

	// Assign configured router to server
	server.router = router
	return server
}
