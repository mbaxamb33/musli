// api/datasources.go

package api

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// convertCompanyDatasourceToResponse converts a database datasource to an API response
func convertCompanyDatasourceToResponse(datasource db.ListDatasourcesByCompanyRow) datasourceResponse {
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

// Helper function to check if a user has access to a contact through company ownership
func (server *Server) userHasAccessToContact(ctx *gin.Context, contactID int32, cognitoSub string) (bool, error) {
	contact, err := server.store.GetContactByID(ctx, contactID)
	if err != nil {
		return false, err
	}

	// Check if the user owns the company that this contact belongs to
	return server.userHasAccessToCompany(ctx, contact.CompanyID, cognitoSub)
}

// createCompanyDatasource handles creating datasources for companies
func (server *Server) createCompanyDatasource(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get company ID from URL param
	companyIDParam := ctx.Param("id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if company exists and belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, int32(companyID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add datasources to this company"})
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
		FileData:   []byte{}, // Empty for non-file datasources
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

	ctx.JSON(http.StatusCreated, gin.H{
		"datasource_id": datasource.DatasourceID,
		"company_id":    companyID,
		"source_type":   string(datasource.SourceType),
		"link":          req.Link,
		"file_name":     req.FileName,
		"message":       "Datasource created and associated with company successfully",
	})
}

// createContactDatasource handles creating datasources for contacts
func (server *Server) createContactDatasource(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get contact ID from URL param
	contactIDParam := ctx.Param("id")
	contactID, err := strconv.Atoi(contactIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Check if contact exists and belongs to a company owned by the authenticated user
	hasAccess, err := server.userHasAccessToContact(ctx, int32(contactID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add datasources to this contact"})
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
		FileData:   []byte{}, // Empty for non-file datasources
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

	ctx.JSON(http.StatusCreated, gin.H{
		"datasource_id": datasource.DatasourceID,
		"contact_id":    contactID,
		"source_type":   string(datasource.SourceType),
		"link":          req.Link,
		"file_name":     req.FileName,
		"message":       "Datasource created and associated with contact successfully",
	})
}

// uploadDatasource handles file uploads for both companies and contacts
func (server *Server) uploadDatasource(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Determine if this is for a company or contact
	entityType := ctx.Param("entity_type") // "companies" or "contacts"
	entityID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Validate entity exists and user has access
	var hasAccess bool
	if entityType == "companies" {
		hasAccess, err = server.userHasAccessToCompany(ctx, int32(entityID), cognitoSub.(string))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
			return
		}
		if !hasAccess {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to upload files to this company"})
			return
		}
	} else if entityType == "contacts" {
		hasAccess, err = server.userHasAccessToContact(ctx, int32(entityID), cognitoSub.(string))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
			return
		}
		if !hasAccess {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to upload files to this contact"})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type, must be 'companies' or 'contacts'"})
		return
	}

	// Parse form data
	sourceType := ctx.PostForm("source_type")
	link := ctx.PostForm("link")

	// Determine datasource type
	var datasourceType db.DatasourceType
	if sourceType == "website" {
		datasourceType = db.DatasourceTypeWebsite
	} else if sourceType == "pdf" {
		datasourceType = db.DatasourceTypePdf
	} else if sourceType == "word_document" {
		datasourceType = db.DatasourceTypeWordDocument
	} else if sourceType == "excel" {
		datasourceType = db.DatasourceTypeExcel
	} else if sourceType == "powerpoint" {
		datasourceType = db.DatasourceTypePowerpoint
	} else if sourceType == "mp3" {
		datasourceType = db.DatasourceTypeMp3
	} else if sourceType == "plain_text" {
		datasourceType = db.DatasourceTypePlainText
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source type"})
		return
	}

	// Handle file upload if provided
	var fileData []byte
	var fileName string
	file, header, err := ctx.Request.FormFile("file")
	if err == nil {
		// File was provided
		defer file.Close()
		fileName = header.Filename
		fileData, err = io.ReadAll(file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
	} else if link == "" {
		// If no file and no link provided
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either file or link must be provided"})
		return
	}

	// Create datasource
	datasourceArg := db.CreateDatasourceParams{
		SourceType: datasourceType,
		Link:       sql.NullString{String: link, Valid: link != ""},
		FileData:   fileData,
		FileName:   sql.NullString{String: fileName, Valid: fileName != ""},
	}

	datasource, err := server.store.CreateDatasource(ctx, datasourceArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create datasource"})
		return
	}

	// Associate datasource with entity
	if entityType == "companies" {
		associateArg := db.AssociateDatasourceWithCompanyParams{
			CompanyID:    int32(entityID),
			DatasourceID: datasource.DatasourceID,
		}

		err = server.store.AssociateDatasourceWithCompany(ctx, associateArg)
		if err != nil {
			// Rollback datasource creation if association fails
			_ = server.store.DeleteDatasource(ctx, datasource.DatasourceID)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with company"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"datasource_id": datasource.DatasourceID,
			"company_id":    entityID,
			"source_type":   string(datasource.SourceType),
			"link":          link,
			"file_name":     fileName,
			"message":       "Datasource created and associated with company successfully",
		})
	} else {
		associateArg := db.AssociateDatasourceWithContactParams{
			ContactID:    int32(entityID),
			DatasourceID: datasource.DatasourceID,
		}

		err = server.store.AssociateDatasourceWithContact(ctx, associateArg)
		if err != nil {
			// Rollback datasource creation if association fails
			_ = server.store.DeleteDatasource(ctx, datasource.DatasourceID)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate datasource with contact"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"datasource_id": datasource.DatasourceID,
			"contact_id":    entityID,
			"source_type":   string(datasource.SourceType),
			"link":          link,
			"file_name":     fileName,
			"message":       "Datasource created and associated with contact successfully",
		})
	}
}

// listCompanyDatasources handles listing datasources for a company
func (server *Server) listCompanyDatasources(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get company ID from URL param
	companyIDParam := ctx.Param("id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if company exists and belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, int32(companyID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view datasources for this company"})
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
		responses[i] = convertCompanyDatasourceToResponse(datasource)
	}

	ctx.JSON(http.StatusOK, responses)
}

// listContactDatasources handles listing datasources for a contact
func (server *Server) listContactDatasources(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get contact ID from URL param
	contactIDParam := ctx.Param("id")
	contactID, err := strconv.Atoi(contactIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Check if contact exists and belongs to a company owned by the authenticated user
	hasAccess, err := server.userHasAccessToContact(ctx, int32(contactID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view datasources for this contact"})
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

// deleteCompanyDatasource handles removing a datasource from a company
func (server *Server) deleteCompanyDatasource(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get company ID and datasource ID from URL params
	companyIDParam := ctx.Param("id")
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

	// Check if company exists and belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, int32(companyID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete datasources from this company"})
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

// deleteContactDatasource handles removing a datasource from a contact
func (server *Server) deleteContactDatasource(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get contact ID and datasource ID from URL params
	contactIDParam := ctx.Param("id")
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

	// Check if contact exists and belongs to a company owned by the authenticated user
	hasAccess, err := server.userHasAccessToContact(ctx, int32(contactID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete datasources from this contact"})
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
