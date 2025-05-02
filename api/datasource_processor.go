// api/datasource_processor.go

package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
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
	// Get datasource ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datasource ID format"})
		return
	}

	// Get datasource from database
	datasource, err := server.store.GetDatasourceByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Datasource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasource"})
		return
	}

	// Process based on datasource type
	var paragraphCount int
	var message string

	switch datasource.SourceType {
	case db.DatasourceTypeWebsite:
		// Process website using web scraper
		if !datasource.Link.Valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Website datasource has no link"})
			return
		}

		count, msg, err := processWebsiteDatasource(ctx, server.store, datasource)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to process website",
				"details": err.Error(),
			})
			return
		}
		paragraphCount = count
		message = msg

	case db.DatasourceTypeWordDocument, db.DatasourceTypePdf:
		// For future implementation
		ctx.JSON(http.StatusNotImplemented, gin.H{
			"error": fmt.Sprintf("Processing %s datasources is not yet implemented", datasource.SourceType),
		})
		return

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Processing for datasource type %s is not supported", datasource.SourceType),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, processDatasourceResponse{
		DatasourceID:   datasource.DatasourceID,
		SourceType:     string(datasource.SourceType),
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
