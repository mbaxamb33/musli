package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// createProjectDatasourceRequest represents the request to create a new datasource for a project
type createProjectDatasourceRequest struct {
	SourceType db.DatasourceType `json:"source_type" binding:"required"`
	Link       string            `json:"link,omitempty"`
	FileName   string            `json:"file_name,omitempty"`
}

// createAndAssociateProjectDatasource handles requests to create a new datasource and associate it with a project
func (server *Server) createAndAssociateProjectDatasource(ctx *gin.Context) {
	// Get project ID from URL param
	projectIDParam := ctx.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Check if project exists
	_, err = server.store.GetProjectByID(ctx, int32(projectID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	// Parse request body
	var req createProjectDatasourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create datasource
	datasourceArg := db.CreateDatasourceParams{
		SourceType: req.SourceType,
		Link:       sql.NullString{String: req.Link, Valid: req.Link != ""},
		FileData:   []byte{}, // Empty for non-file datasources
		FileName:   sql.NullString{String: req.FileName, Valid: req.FileName != ""},
	}

	datasource, err := server.store.CreateDatasource(ctx, datasourceArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create datasource"})
		return
	}

	// Associate datasource with project
	associateArg := db.AssociateDatasourceWithProjectParams{
		ProjectID:    int32(projectID),
		DatasourceID: datasource.DatasourceID,
	}

	err = server.store.AssociateDatasourceWithProject(ctx, associateArg)
	if err != nil {
		// Rollback datasource creation if association fails
		_ = server.store.DeleteDatasource(ctx, datasource.DatasourceID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with project"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"datasource_id": datasource.DatasourceID,
		"project_id":    projectID,
		"message":       "Datasource created and associated with project successfully",
	})
}

// associateDatasourceWithProject handles requests to associate an existing datasource with a project
func (server *Server) associateDatasourceWithProject(ctx *gin.Context) {
	// Get project ID from URL param
	projectIDParam := ctx.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Check if project exists
	_, err = server.store.GetProjectByID(ctx, int32(projectID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	// Parse request body
	var req associateDatasourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if datasource exists
	_, err = server.store.GetDatasourceByID(ctx, req.DatasourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasource"})
		return
	}

	// Check if association already exists
	_, err = server.store.GetProjectDatasourceAssociation(ctx, db.GetProjectDatasourceAssociationParams{
		ProjectID:    int32(projectID),
		DatasourceID: req.DatasourceID,
	})
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datasource is already associated with this project"})
		return
	} else if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing association"})
		return
	}

	// Associate datasource with project
	err = server.store.AssociateDatasourceWithProject(ctx, db.AssociateDatasourceWithProjectParams{
		ProjectID:    int32(projectID),
		DatasourceID: req.DatasourceID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with project"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"project_id":    projectID,
		"datasource_id": req.DatasourceID,
		"message":       "Datasource associated with project successfully",
	})
}

// removeDatasourceFromProject handles requests to remove a datasource association from a project
func (server *Server) removeDatasourceFromProject(ctx *gin.Context) {
	// Get project ID and datasource ID from URL params
	projectIDParam := ctx.Param("id")
	datasourceIDParam := ctx.Param("datasource_id")

	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	datasourceID, err := strconv.Atoi(datasourceIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datasource ID format"})
		return
	}

	// Check if association exists
	_, err = server.store.GetProjectDatasourceAssociation(ctx, db.GetProjectDatasourceAssociationParams{
		ProjectID:    int32(projectID),
		DatasourceID: int32(datasourceID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Association between project and datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing association"})
		return
	}

	// Remove association
	err = server.store.RemoveDatasourceFromProject(ctx, db.RemoveDatasourceFromProjectParams{
		ProjectID:    int32(projectID),
		DatasourceID: int32(datasourceID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove datasource from project"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Datasource removed from project successfully"})
}

// listDatasourcesByProject handles requests to get all datasources for a specific project
func (server *Server) listDatasourcesByProject(ctx *gin.Context) {
	// Get project ID from URL param
	projectIDParam := ctx.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Check if project exists
	_, err = server.store.GetProjectByID(ctx, int32(projectID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	// Default pagination settings
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	limitParam := ctx.DefaultQuery("limit", "10")
	offsetParam := ctx.DefaultQuery("offset", "0")

	// Convert string params to integers
	limit, err = strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err = strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get datasources for project
	datasources, err := server.store.ListDatasourcesByProject(ctx, db.ListDatasourcesByProjectParams{
		ProjectID: int32(projectID),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasources"})
		return
	}

	// Convert datasources to response format
	responses := make([]datasourceResponse, len(datasources))
	for i, datasource := range datasources {
		link := ""
		if datasource.Link.Valid {
			link = datasource.Link.String
		}

		fileName := ""
		if datasource.FileName.Valid {
			fileName = datasource.FileName.String
		}

		createdAt := ""
		if datasource.CreatedAt.Valid {
			createdAt = datasource.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		responses[i] = datasourceResponse{
			DatasourceID: datasource.DatasourceID,
			SourceType:   datasource.SourceType,
			Link:         link,
			FileName:     fileName,
			CreatedAt:    createdAt,
		}
	}

	ctx.JSON(http.StatusOK, responses)
}
