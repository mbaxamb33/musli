// api/paragraphs.go

package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// paragraphResponse represents the API response structure for paragraph data
type paragraphResponse struct {
	ParagraphID  int32  `json:"paragraph_id"`
	DatasourceID int32  `json:"datasource_id"`
	Title        string `json:"title,omitempty"`
	MainIdea     string `json:"main_idea,omitempty"`
	Content      string `json:"content"`
	CreatedAt    string `json:"created_at,omitempty"`
}

// createParagraphRequest represents the request to create a new paragraph
type createParagraphRequest struct {
	DatasourceID int32  `json:"datasource_id" binding:"required"`
	Title        string `json:"title,omitempty"`
	MainIdea     string `json:"main_idea,omitempty"`
	Content      string `json:"content" binding:"required"`
}

// updateParagraphRequest represents the request to update a paragraph
type updateParagraphRequest struct {
	Title    string `json:"title,omitempty"`
	MainIdea string `json:"main_idea,omitempty"`
	Content  string `json:"content" binding:"required"`
}

// convertParagraphToResponse converts a database paragraph to an API response
func convertParagraphToResponse(paragraph db.Paragraph) paragraphResponse {
	title := ""
	if paragraph.Title.Valid {
		title = paragraph.Title.String
	}

	mainIdea := ""
	if paragraph.MainIdea.Valid {
		mainIdea = paragraph.MainIdea.String
	}

	createdAt := ""
	if paragraph.CreatedAt.Valid {
		createdAt = paragraph.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return paragraphResponse{
		ParagraphID:  paragraph.ParagraphID,
		DatasourceID: paragraph.DatasourceID,
		Title:        title,
		MainIdea:     mainIdea,
		Content:      paragraph.Content,
		CreatedAt:    createdAt,
	}
}

// convertCompanyParagraphToResponse converts a company paragraph to an API response
func convertCompanyParagraphToResponse(paragraph db.GetCompanyParagraphsRow) paragraphResponse {
	title := ""
	if paragraph.Title.Valid {
		title = paragraph.Title.String
	}

	mainIdea := ""
	if paragraph.MainIdea.Valid {
		mainIdea = paragraph.MainIdea.String
	}

	createdAt := ""
	if paragraph.CreatedAt.Valid {
		createdAt = paragraph.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return paragraphResponse{
		ParagraphID:  paragraph.ParagraphID,
		DatasourceID: paragraph.DatasourceID,
		Title:        title,
		MainIdea:     mainIdea,
		Content:      paragraph.Content,
		CreatedAt:    createdAt,
	}
}

// convertContactParagraphToResponse converts a contact paragraph to an API response
func convertContactParagraphToResponse(paragraph db.GetContactParagraphsRow) paragraphResponse {
	title := ""
	if paragraph.Title.Valid {
		title = paragraph.Title.String
	}

	mainIdea := ""
	if paragraph.MainIdea.Valid {
		mainIdea = paragraph.MainIdea.String
	}

	createdAt := ""
	if paragraph.CreatedAt.Valid {
		createdAt = paragraph.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return paragraphResponse{
		ParagraphID:  paragraph.ParagraphID,
		DatasourceID: paragraph.DatasourceID,
		Title:        title,
		MainIdea:     mainIdea,
		Content:      paragraph.Content,
		CreatedAt:    createdAt,
	}
}

// createParagraph handles requests to create a new paragraph
func (server *Server) createParagraph(ctx *gin.Context) {
	// Parse request body
	var req createParagraphRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if datasource exists
	_, err := server.store.GetDatasourceByID(ctx, req.DatasourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasource"})
		return
	}

	// Create paragraph in database
	arg := db.CreateParagraphParams{
		DatasourceID: req.DatasourceID,
		Title:        sql.NullString{String: req.Title, Valid: req.Title != ""},
		MainIdea:     sql.NullString{String: req.MainIdea, Valid: req.MainIdea != ""},
		Content:      req.Content,
	}

	paragraph, err := server.store.CreateParagraph(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create paragraph"})
		return
	}

	// Return created paragraph as response
	ctx.JSON(http.StatusCreated, convertParagraphToResponse(paragraph))
}

// getParagraphByID handles requests to get a specific paragraph
func (server *Server) getParagraphByID(ctx *gin.Context) {
	// Get paragraph ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paragraph ID format"})
		return
	}

	// Get paragraph from database
	paragraph, err := server.store.GetParagraphByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Paragraph not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paragraph"})
		return
	}

	// Return paragraph response
	ctx.JSON(http.StatusOK, convertParagraphToResponse(paragraph))
}

// updateParagraph handles requests to update an existing paragraph
func (server *Server) updateParagraph(ctx *gin.Context) {
	// Get paragraph ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paragraph ID format"})
		return
	}

	// Check if paragraph exists
	_, err = server.store.GetParagraphByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Paragraph not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paragraph"})
		return
	}

	// Parse request body
	var req updateParagraphRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update paragraph in database
	arg := db.UpdateParagraphParams{
		ParagraphID: int32(id),
		Title:       sql.NullString{String: req.Title, Valid: req.Title != ""},
		MainIdea:    sql.NullString{String: req.MainIdea, Valid: req.MainIdea != ""},
		Content:     req.Content,
	}

	updatedParagraph, err := server.store.UpdateParagraph(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update paragraph"})
		return
	}

	// Return updated paragraph
	ctx.JSON(http.StatusOK, convertParagraphToResponse(updatedParagraph))
}

// deleteParagraph handles requests to delete a paragraph
func (server *Server) deleteParagraph(ctx *gin.Context) {
	// Get paragraph ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paragraph ID format"})
		return
	}

	// Check if paragraph exists
	_, err = server.store.GetParagraphByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Paragraph not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paragraph"})
		return
	}

	// Delete paragraph
	err = server.store.DeleteParagraph(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete paragraph"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Paragraph deleted successfully"})
}

// listParagraphsByDatasource handles requests to get all paragraphs for a specific datasource
func (server *Server) listParagraphsByDatasource(ctx *gin.Context) {
	// Get datasource ID from URL param
	datasourceIDParam := ctx.Param("datasource_id")
	datasourceID, err := strconv.Atoi(datasourceIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datasource ID format"})
		return
	}

	// Check if datasource exists
	_, err = server.store.GetDatasourceByID(ctx, int32(datasourceID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasource"})
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

	// Get paragraphs from database
	paragraphs, err := server.store.ListParagraphsByDatasource(ctx, db.ListParagraphsByDatasourceParams{
		DatasourceID: int32(datasourceID),
		Limit:        int32(limit),
		Offset:       int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paragraphs"})
		return
	}

	// Convert paragraphs to response format
	responses := make([]paragraphResponse, len(paragraphs))
	for i, paragraph := range paragraphs {
		responses[i] = convertParagraphToResponse(paragraph)
	}

	ctx.JSON(http.StatusOK, responses)
}

// listCompanyParagraphs handles requests to get all paragraphs for a specific company
func (server *Server) listCompanyParagraphs(ctx *gin.Context) {
	// Get company ID from URL param
	companyIDParam := ctx.Param("id")
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

	// Get paragraphs for company
	paragraphs, err := server.store.GetCompanyParagraphs(ctx, db.GetCompanyParagraphsParams{
		CompanyID: int32(companyID),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paragraphs"})
		return
	}

	// Convert paragraphs to response format
	responses := make([]paragraphResponse, len(paragraphs))
	for i, paragraph := range paragraphs {
		responses[i] = convertCompanyParagraphToResponse(paragraph)
	}

	ctx.JSON(http.StatusOK, responses)
}

// listContactParagraphs handles requests to get all paragraphs for a specific contact
func (server *Server) listContactParagraphs(ctx *gin.Context) {
	// Get contact ID from URL param
	contactIDParam := ctx.Param("id")
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

	// Get paragraphs for contact
	paragraphs, err := server.store.GetContactParagraphs(ctx, db.GetContactParagraphsParams{
		ContactID: int32(contactID),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paragraphs"})
		return
	}

	// Convert paragraphs to response format
	responses := make([]paragraphResponse, len(paragraphs))
	for i, paragraph := range paragraphs {
		responses[i] = convertContactParagraphToResponse(paragraph)
	}

	ctx.JSON(http.StatusOK, responses)
}

// searchCompanyParagraphs handles requests to search paragraphs for a specific company
func (server *Server) searchCompanyParagraphs(ctx *gin.Context) {
	// Get company ID from URL param
	companyIDParam := ctx.Param("id")
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

	// Get search query from URL param
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
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

	// Search paragraphs for company
	paragraphs, err := server.store.SearchCompanyParagraphs(ctx, db.SearchCompanyParagraphsParams{
		CompanyID: int32(companyID),
		Column2:   sql.NullString{String: query, Valid: true},
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search paragraphs"})
		return
	}

	// Prepare response
	type searchParagraphResponse struct {
		ParagraphID  int32  `json:"paragraph_id"`
		CompanyID    int32  `json:"company_id"`
		CompanyName  string `json:"company_name"`
		DatasourceID int32  `json:"datasource_id"`
		SourceType   string `json:"source_type"`
		Title        string `json:"title,omitempty"`
		MainIdea     string `json:"main_idea,omitempty"`
		Content      string `json:"content"`
	}

	// Convert paragraphs to response format
	responses := make([]searchParagraphResponse, len(paragraphs))
	for i, paragraph := range paragraphs {
		title := ""
		if paragraph.Title.Valid {
			title = paragraph.Title.String
		}

		mainIdea := ""
		if paragraph.MainIdea.Valid {
			mainIdea = paragraph.MainIdea.String
		}

		responses[i] = searchParagraphResponse{
			ParagraphID:  paragraph.ParagraphID,
			CompanyID:    paragraph.CompanyID,
			CompanyName:  paragraph.CompanyName,
			DatasourceID: paragraph.DatasourceID,
			SourceType:   string(paragraph.SourceType),
			Title:        title,
			MainIdea:     mainIdea,
			Content:      paragraph.Content,
		}
	}

	ctx.JSON(http.StatusOK, responses)
}

// searchContactParagraphs handles requests to search paragraphs for a specific contact
func (server *Server) searchContactParagraphs(ctx *gin.Context) {
	// Get contact ID from URL param
	contactIDParam := ctx.Param("id")
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

	// Get search query from URL param
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
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

	// Search paragraphs for contact
	paragraphs, err := server.store.SearchContactParagraphs(ctx, db.SearchContactParagraphsParams{
		ContactID: int32(contactID),
		Column2:   sql.NullString{String: query, Valid: true},
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search paragraphs"})
		return
	}

	// Prepare response
	type searchParagraphResponse struct {
		ParagraphID  int32  `json:"paragraph_id"`
		ContactID    int32  `json:"contact_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		DatasourceID int32  `json:"datasource_id"`
		SourceType   string `json:"source_type"`
		Title        string `json:"title,omitempty"`
		MainIdea     string `json:"main_idea,omitempty"`
		Content      string `json:"content"`
	}

	// Convert paragraphs to response format
	responses := make([]searchParagraphResponse, len(paragraphs))
	for i, paragraph := range paragraphs {
		title := ""
		if paragraph.Title.Valid {
			title = paragraph.Title.String
		}

		mainIdea := ""
		if paragraph.MainIdea.Valid {
			mainIdea = paragraph.MainIdea.String
		}

		responses[i] = searchParagraphResponse{
			ParagraphID:  paragraph.ParagraphID,
			ContactID:    paragraph.ContactID,
			FirstName:    paragraph.FirstName,
			LastName:     paragraph.LastName,
			DatasourceID: paragraph.DatasourceID,
			SourceType:   string(paragraph.SourceType),
			Title:        title,
			MainIdea:     mainIdea,
			Content:      paragraph.Content,
		}
	}

	ctx.JSON(http.StatusOK, responses)
}
