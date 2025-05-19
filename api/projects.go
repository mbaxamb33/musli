package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// createProjectRequest represents the request body for creating a project
type createProjectRequest struct {
	ProjectName string `json:"project_name" binding:"required"`
	MainIdea    string `json:"main_idea" binding:"omitempty"`
}

// projectResponse represents the API response structure for project data
type projectResponse struct {
	ProjectID   int32  `json:"project_id"`
	UserID      int32  `json:"user_id"`
	ProjectName string `json:"project_name"`
	MainIdea    string `json:"main_idea,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// convertProjectToResponse converts a database project model to an API response
func convertProjectToResponse(project db.Project) projectResponse {
	createdAt := ""
	if project.CreatedAt.Valid {
		createdAt = project.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	updatedAt := ""
	if project.UpdatedAt.Valid {
		updatedAt = project.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	mainIdea := ""
	if project.MainIdea.Valid {
		mainIdea = project.MainIdea.String
	}

	return projectResponse{
		ProjectID:   project.ProjectID,
		UserID:      project.UserID,
		ProjectName: project.ProjectName,
		MainIdea:    mainIdea,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// createProject handles requests to create a new project
func (server *Server) createProject(ctx *gin.Context) {
	var req createProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from authentication context or default to a test user
	// In a real application, you would get this from authenticated user
	userID := int32(1) // TODO: Replace with actual authenticated user ID

	// Convert request to database params
	arg := db.CreateProjectParams{
		UserID:      userID,
		ProjectName: req.ProjectName,
		MainIdea:    sql.NullString{String: req.MainIdea, Valid: req.MainIdea != ""},
	}

	// Create project in database
	project, err := server.store.CreateProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Return created project as response
	ctx.JSON(http.StatusCreated, convertProjectToResponse(project))
}

// getProjectByID handles requests to get a specific project
func (server *Server) getProjectByID(ctx *gin.Context) {
	// Get project ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Get project from database
	project, err := server.store.GetProjectByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	// Return project response
	ctx.JSON(http.StatusOK, convertProjectToResponse(project))
}

// listProjects handles requests to get projects for a user with pagination
func (server *Server) listProjects(ctx *gin.Context) {
	// Get user ID from authentication context or default to a test user
	// In a real application, you would get this from authenticated user
	userID := int32(1) // TODO: Replace with actual authenticated user ID

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

	// Get projects from database
	projects, err := server.store.ListProjectsByUserID(ctx, db.ListProjectsByUserIDParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	// Convert projects to response format
	projectResponses := make([]projectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = convertProjectToResponse(project)
	}

	ctx.JSON(http.StatusOK, projectResponses)
}

// deleteProject handles requests to delete a project
func (server *Server) deleteProject(ctx *gin.Context) {
	// Get project ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Delete project from database
	err = server.store.DeleteProject(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
