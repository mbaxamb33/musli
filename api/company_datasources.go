// api/company_datasources.go

package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// datasourceResponse represents the API response structure for datasource data
type datasourceResponse struct {
	DatasourceID int32             `json:"datasource_id"`
	SourceType   db.DatasourceType `json:"source_type"`
	Link         string            `json:"link,omitempty"`
	FileName     string            `json:"file_name,omitempty"`
	CreatedAt    string            `json:"created_at,omitempty"`
}

// associateDatasourceRequest represents the request to associate a datasource with a company
type associateDatasourceRequest struct {
	DatasourceID int32 `json:"datasource_id" binding:"required"`
}

// createDatasourceRequest represents the request to create a new datasource
type createDatasourceRequest struct {
	SourceType db.DatasourceType `json:"source_type" binding:"required"`
	Link       string            `json:"link,omitempty"`
	FileData   []byte            `json:"file_data,omitempty"`
	FileName   string            `json:"file_name,omitempty"`
}

// convertDatasourceToResponse converts a database datasource row to an API response
func convertDatasourceToResponse(datasource db.ListDatasourcesByCompanyRow) datasourceResponse {
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

	return datasourceResponse{
		DatasourceID: datasource.DatasourceID,
		SourceType:   datasource.SourceType,
		Link:         link,
		FileName:     fileName,
		CreatedAt:    createdAt,
	}
}

// createAndAssociateDatasource handles requests to create a new datasource and associate it with a company
func (server *Server) createAndAssociateDatasource(ctx *gin.Context) {
	// Get company ID from URL param
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if company exists
	_, err = server.store.GetCompanyByID(ctx, int32(companyID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}

	// Parse request body
	var req createDatasourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create datasource
	datasourceArg := db.CreateDatasourceParams{
		SourceType: req.SourceType,
		Link:       sql.NullString{String: req.Link, Valid: req.Link != ""},
		FileData:   req.FileData,
		FileName:   sql.NullString{String: req.FileName, Valid: req.FileName != ""},
	}

	datasource, err := server.store.CreateDatasource(ctx, datasourceArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create datasource"})
		return
	}

	// Associate datasource with company
	associateArg := db.AssociateDatasourceWithCompanyParams{
		CompanyID:    int32(companyID),
		DatasourceID: datasource.DatasourceID,
	}

	err = server.store.AssociateDatasourceWithCompany(ctx, associateArg)
	if err != nil {
		// Rollback datasource creation if association fails
		_ = server.store.DeleteDatasource(ctx, datasource.DatasourceID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with company"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"datasource_id": datasource.DatasourceID,
		"company_id":    companyID,
		"message":       "Datasource created and associated with company successfully",
	})
}

// associateDatasourceWithCompany handles requests to associate an existing datasource with a company
func (server *Server) associateDatasourceWithCompany(ctx *gin.Context) {
	// Get company ID from URL param
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if company exists
	_, err = server.store.GetCompanyByID(ctx, int32(companyID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
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
	_, err = server.store.GetCompanyDatasourceAssociation(ctx, db.GetCompanyDatasourceAssociationParams{
		CompanyID:    int32(companyID),
		DatasourceID: req.DatasourceID,
	})
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datasource is already associated with this company"})
		return
	} else if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing association"})
		return
	}

	// Associate datasource with company
	err = server.store.AssociateDatasourceWithCompany(ctx, db.AssociateDatasourceWithCompanyParams{
		CompanyID:    int32(companyID),
		DatasourceID: req.DatasourceID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with company"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"company_id":    companyID,
		"datasource_id": req.DatasourceID,
		"message":       "Datasource associated with company successfully",
	})
}

// removeDatasourceFromCompany handles requests to remove a datasource association from a company
func (server *Server) removeDatasourceFromCompany(ctx *gin.Context) {
	// Get company ID and datasource ID from URL params
	companyIDParam := ctx.Param("company_id")
	datasourceIDParam := ctx.Param("datasource_id")

	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	datasourceID, err := strconv.Atoi(datasourceIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datasource ID format"})
		return
	}

	// Check if association exists
	_, err = server.store.GetCompanyDatasourceAssociation(ctx, db.GetCompanyDatasourceAssociationParams{
		CompanyID:    int32(companyID),
		DatasourceID: int32(datasourceID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Association between company and datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing association"})
		return
	}

	// Remove association
	err = server.store.RemoveDatasourceFromCompany(ctx, db.RemoveDatasourceFromCompanyParams{
		CompanyID:    int32(companyID),
		DatasourceID: int32(datasourceID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove datasource from company"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Datasource removed from company successfully"})
}

// listDatasourcesByCompany handles requests to get all datasources for a specific company
func (server *Server) listDatasourcesByCompany(ctx *gin.Context) {
	// Get company ID from URL param
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if company exists
	_, err = server.store.GetCompanyByID(ctx, int32(companyID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
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

	// Get datasources for company
	datasources, err := server.store.ListDatasourcesByCompany(ctx, db.ListDatasourcesByCompanyParams{
		CompanyID: int32(companyID),
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
		responses[i] = convertDatasourceToResponse(datasource)
	}

	ctx.JSON(http.StatusOK, responses)
}
