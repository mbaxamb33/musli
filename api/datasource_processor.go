// api/datasource_processor.go

package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
	docscraper "github.com/mbaxamb3/nusli/document_scraper"
	"github.com/mbaxamb3/nusli/scraper"
)

// processDatasourceRequest represents the request to process a datasource
type processDatasourceRequest struct {
	DatasourceID int32 `json:"datasource_id" binding:"required"`
}

// processDatasourceResponse represents the response after processing a datasource
type processDatasourceResponse struct {
	DatasourceID   int32  `json:"datasource_id"`
	SourceType     string `json:"source_type"`
	ParagraphCount int    `json:"paragraph_count"`
	Message        string `json:"message"`
}

// processDatasourceByID handles processing a specific datasource and generating paragraphs
func (server *Server) processDatasourceByID(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	_, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get datasource ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datasource ID format"})
		return
	}

	// Get basic datasource info first
	datasourceBasic, err := server.store.GetDatasourceByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasource"})
		return
	}

	// // Verify that the user has access to this datasource
	// hasAccess, err := server.userHasAccessToDatasource(ctx, int32(id), cognitoSub.(string))
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify datasource access"})
	// 	return
	// }
	// if !hasAccess {
	// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to process this datasource"})
	// 	return
	// }

	// Process based on datasource type
	var paragraphCount int
	var message string

	switch datasourceBasic.SourceType {
	case db.DatasourceTypeWebsite:
		// Process website using web scraper
		if !datasourceBasic.Link.Valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Website datasource has no link"})
			return
		}

		count, msg, err := processWebsiteDatasource(ctx, server.store, datasourceBasic)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to process website",
				"details": err.Error(),
			})
			return
		}
		paragraphCount = count
		message = msg

	case db.DatasourceTypeWordDocument:
		// Process Word document - need to get the full datasource with file data
		datasourceFull, err := server.store.GetFullDatasourceByID(ctx, int32(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch full datasource data"})
			return
		}

		if !datasourceFull.FileName.Valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Word document datasource has no file name"})
			return
		}

		count, msg, err := processWordDocumentDatasource(ctx, server.store, datasourceFull)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to process Word document",
				"details": err.Error(),
			})
			return
		}
		paragraphCount = count
		message = msg

	case db.DatasourceTypePdf:
		// For future implementation
		ctx.JSON(http.StatusNotImplemented, gin.H{
			"error": fmt.Sprintf("Processing %s datasources is not yet implemented", datasourceBasic.SourceType),
		})
		return

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Processing for datasource type %s is not supported", datasourceBasic.SourceType),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, processDatasourceResponse{
		DatasourceID:   datasourceBasic.DatasourceID,
		SourceType:     string(datasourceBasic.SourceType),
		ParagraphCount: paragraphCount,
		Message:        message,
	})
}

// processWebsiteDatasource processes a website datasource using the scraper
func processWebsiteDatasource(ctx *gin.Context, store *db.Store, datasource db.GetDatasourceByIDRow) (int, string, error) {
	// Create enhanced scraper with the link
	link := datasource.Link.String
	enhancedScraper, err := scraper.NewEnhancedScraper(link, 1) // Depth 1 to avoid going too deep
	if err != nil {
		return 0, "", fmt.Errorf("failed to create scraper: %w", err)
	}

	// Extract content
	err = enhancedScraper.Run()
	if err != nil {
		return 0, "", fmt.Errorf("failed to scrape website: %w", err)
	}

	// Create paragraphs from extracted content
	paragraphCount := 0
	for _, item := range enhancedScraper.ContentItems {
		// Skip items with very short paragraphs
		if len(item.Paragraph) < 100 {
			continue
		}

		// Create paragraph
		paragraphParams := db.CreateParagraphParams{
			DatasourceID: datasource.DatasourceID,
			Title:        sql.NullString{String: item.Title, Valid: item.Title != ""},
			MainIdea:     sql.NullString{String: "", Valid: false}, // Could implement a summarizer in the future
			Content:      item.Paragraph,
		}

		_, err := store.CreateParagraph(ctx, paragraphParams)
		if err != nil {
			return paragraphCount, "", fmt.Errorf("failed to create paragraph: %w", err)
		}
		paragraphCount++
	}

	message := fmt.Sprintf("Successfully extracted %d paragraphs from %s", paragraphCount, link)
	return paragraphCount, message, nil
}

// processWordDocumentDatasource processes a Word document datasource
func processWordDocumentDatasource(ctx *gin.Context, store *db.Store, datasource db.Datasource) (int, string, error) {
	// Create a temporary file to save the Word document
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("doc_%d_%s", datasource.DatasourceID, datasource.FileName.String))

	// Save file data to temporary file
	err := os.WriteFile(tempFile, datasource.FileData, 0644)
	if err != nil {
		return 0, "", fmt.Errorf("failed to save temporary file: %w", err)
	}
	defer os.Remove(tempFile) // Clean up

	// Create document scraper
	docScraper, err := docscraper.NewDocumentScraper(tempFile)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create document scraper: %w", err)
	}

	// Extract content
	err = docScraper.Run()
	if err != nil {
		return 0, "", fmt.Errorf("failed to scrape document: %w", err)
	}

	// Create paragraphs from extracted content
	paragraphCount := 0
	for _, item := range docScraper.ContentItems {
		// Skip items with very short paragraphs
		if len(item.Paragraph) < 100 {
			continue
		}

		// Create paragraph
		paragraphParams := db.CreateParagraphParams{
			DatasourceID: datasource.DatasourceID,
			Title:        sql.NullString{String: item.Title, Valid: item.Title != ""},
			MainIdea:     sql.NullString{String: "", Valid: false}, // Could implement a summarizer in the future
			Content:      item.Paragraph,
		}

		_, err := store.CreateParagraph(ctx, paragraphParams)
		if err != nil {
			return paragraphCount, "", fmt.Errorf("failed to create paragraph: %w", err)
		}
		paragraphCount++
	}

	message := fmt.Sprintf("Successfully extracted %d paragraphs from document %s", paragraphCount, datasource.FileName.String)
	return paragraphCount, message, nil
}
