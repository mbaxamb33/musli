// api/contact_datasources.go

package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// convertContactDatasourceToResponse converts a database datasource row to an API response
func convertContactDatasourceToResponse(datasource db.ListDatasourcesByContactRow) datasourceResponse {
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

// createAndAssociateContactDatasource handles requests to create a new datasource and associate it with a contact
func (server *Server) createAndAssociateContactDatasource(ctx *gin.Context) {
	// Get contact ID from URL param
	contactIDParam := ctx.Param("contact_id")
	contactID, err := strconv.Atoi(contactIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Check if contact exists
	_, err = server.store.GetContactByID(ctx, int32(contactID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
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

	// Associate datasource with contact
	associateArg := db.AssociateDatasourceWithContactParams{
		ContactID:    int32(contactID),
		DatasourceID: datasource.DatasourceID,
	}

	err = server.store.AssociateDatasourceWithContact(ctx, associateArg)
	if err != nil {
		// Rollback datasource creation if association fails
		_ = server.store.DeleteDatasource(ctx, datasource.DatasourceID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with contact"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"datasource_id": datasource.DatasourceID,
		"contact_id":    contactID,
		"message":       "Datasource created and associated with contact successfully",
	})
}

// associateDatasourceWithContact handles requests to associate an existing datasource with a contact
func (server *Server) associateDatasourceWithContact(ctx *gin.Context) {
	// Get contact ID from URL param
	contactIDParam := ctx.Param("contact_id")
	contactID, err := strconv.Atoi(contactIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Check if contact exists
	_, err = server.store.GetContactByID(ctx, int32(contactID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
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
	_, err = server.store.GetContactDatasourceAssociation(ctx, db.GetContactDatasourceAssociationParams{
		ContactID:    int32(contactID),
		DatasourceID: req.DatasourceID,
	})
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datasource is already associated with this contact"})
		return
	} else if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing association"})
		return
	}

	// Associate datasource with contact
	err = server.store.AssociateDatasourceWithContact(ctx, db.AssociateDatasourceWithContactParams{
		ContactID:    int32(contactID),
		DatasourceID: req.DatasourceID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with contact"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"contact_id":    contactID,
		"datasource_id": req.DatasourceID,
		"message":       "Datasource associated with contact successfully",
	})
}

// removeDatasourceFromContact handles requests to remove a datasource association from a contact
func (server *Server) removeDatasourceFromContact(ctx *gin.Context) {
	// Get contact ID and datasource ID from URL params
	contactIDParam := ctx.Param("contact_id")
	datasourceIDParam := ctx.Param("datasource_id")

	contactID, err := strconv.Atoi(contactIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	datasourceID, err := strconv.Atoi(datasourceIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datasource ID format"})
		return
	}

	// Check if association exists
	_, err = server.store.GetContactDatasourceAssociation(ctx, db.GetContactDatasourceAssociationParams{
		ContactID:    int32(contactID),
		DatasourceID: int32(datasourceID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Association between contact and datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing association"})
		return
	}

	// Remove association
	err = server.store.RemoveDatasourceFromContact(ctx, db.RemoveDatasourceFromContactParams{
		ContactID:    int32(contactID),
		DatasourceID: int32(datasourceID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove datasource from contact"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Datasource removed from contact successfully"})
}

// listDatasourcesByContact handles requests to get all datasources for a specific contact
func (server *Server) listDatasourcesByContact(ctx *gin.Context) {
	// Get contact ID from URL param
	contactIDParam := ctx.Param("contact_id")
	contactID, err := strconv.Atoi(contactIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Check if contact exists
	_, err = server.store.GetContactByID(ctx, int32(contactID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
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

	// Get datasources for contact
	datasources, err := server.store.ListDatasourcesByContact(ctx, db.ListDatasourcesByContactParams{
		ContactID: int32(contactID),
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
		responses[i] = convertContactDatasourceToResponse(datasource)
	}

	ctx.JSON(http.StatusOK, responses)
}
