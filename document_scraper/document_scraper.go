// document_scraper/document_scraper.go

package docscraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/unidoc/unioffice/document"
)

// ContentItem represents a piece of content extracted from the document
type ContentItem struct {
	Heading      string   // Section heading
	HeadingPath  []string // Full path of headings (for hierarchy)
	HeadingLevel int      // Level of the parent heading (1-based)
	Title        string   // Cleaned up title for storage
	Paragraph    string   // The actual content - CHANGED from Content to Paragraph
	Hash         string   // Unique hash to identify content
}

// DocumentScraper handles extraction from Word documents
type DocumentScraper struct {
	FilePath     string
	ContentItems []ContentItem
	seenContent  map[string]bool // Track already seen content by hash
	mu           sync.Mutex
}

// NewDocumentScraper creates a new document scraper instance
func NewDocumentScraper(filePath string) (*DocumentScraper, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	return &DocumentScraper{
		FilePath:     filePath,
		ContentItems: []ContentItem{},
		seenContent:  make(map[string]bool),
	}, nil
}

// Run executes the complete document scraping process
func (ds *DocumentScraper) Run() error {
	// Open the document
	doc, err := document.Open(ds.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	// Process document content
	err = ds.extractStructuredContent(doc)
	if err != nil {
		return fmt.Errorf("error extracting content: %w", err)
	}

	// Apply additional deduplication
	ds.removeDuplicateContent()

	return nil
}

// extractStructuredContent processes the document to extract structured content
func (ds *DocumentScraper) extractStructuredContent(doc *document.Document) error {
	// Current heading context for associating content
	var currentHeadingPath []string
	currentHeadingLevel := 0

	// Process all paragraphs
	for _, para := range doc.Paragraphs() {
		// Get the text content of the paragraph
		paraText := ""
		for _, run := range para.Runs() {
			paraText += run.Text()
		}

		// Skip empty paragraphs
		if strings.TrimSpace(paraText) == "" {
			continue
		}

		// Check for heading style
		isHeading := false
		headingLevel := 0

		// Determine if this is a heading and what level
		if para.Properties().Style != nil {
			styleFunc := para.Properties().Style
			style := styleFunc()
			if strings.HasPrefix(style, "Heading") && len(style) > 7 {
				// Try to extract heading level
				levelStr := style[7:]
				fmt.Sscanf(levelStr, "%d", &headingLevel)
				if headingLevel > 0 {
					isHeading = true
				}
			}
		}

		// If it's a heading, update the heading context
		if isHeading {
			headingText := strings.TrimSpace(paraText)

			// Update heading path based on level
			if headingLevel <= len(currentHeadingPath) {
				// If it's at the same or higher level than current, pop back to appropriate level
				currentHeadingPath = currentHeadingPath[:headingLevel-1]
			}
			// Now add the new heading
			currentHeadingPath = append(currentHeadingPath, headingText)
			currentHeadingLevel = headingLevel
			continue
		}

		// Handle normal paragraph content
		content := strings.TrimSpace(paraText)
		if countWords(content) >= 10 {
			// Create content item
			headingText := "Untitled Section"
			if len(currentHeadingPath) > 0 {
				headingText = currentHeadingPath[len(currentHeadingPath)-1]
			}

			// Create clean title (for database storage)
			title := cleanText(headingText)
			if title == "" {
				title = "Section"
			}

			// Create content item
			item := ContentItem{
				Heading:      headingText,
				HeadingPath:  append([]string{}, currentHeadingPath...),
				HeadingLevel: currentHeadingLevel,
				Title:        title,
				Paragraph:    cleanText(content),
			}

			// Generate hash
			item.Hash = generateContentHash(item.Heading, item.Paragraph)

			// Add to collection
			ds.addContentItem(item)
		}
	}

	return nil
}

// addContentItem adds a content item to the document scraper
func (ds *DocumentScraper) addContentItem(item ContentItem) {
	// Skip if we've seen this content before
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if !ds.seenContent[item.Hash] {
		ds.seenContent[item.Hash] = true
		ds.ContentItems = append(ds.ContentItems, item)
	}
}

// removeDuplicateContent filters out any duplicate content
func (ds *DocumentScraper) removeDuplicateContent() {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Create a map to track unique content
	seen := make(map[string]bool)
	var uniqueItems []ContentItem

	// Only keep unique content based on hash
	for _, item := range ds.ContentItems {
		if !seen[item.Hash] {
			seen[item.Hash] = true
			uniqueItems = append(uniqueItems, item)
		}
	}

	// Replace with deduplicated list
	ds.ContentItems = uniqueItems
	fmt.Printf("Removed %d duplicate items\n", len(ds.seenContent)-len(uniqueItems))
}

// generateContentHash creates a unique hash for content
func generateContentHash(heading, content string) string {
	// Normalize content before hashing
	normalizedHeading := cleanText(heading)
	normalizedContent := cleanText(content)

	// Create hash from combined content
	h := sha256.New()
	h.Write([]byte(normalizedHeading + "|" + normalizedContent))
	return hex.EncodeToString(h.Sum(nil))
}

// cleanText normalizes text for better comparison
func cleanText(s string) string {
	// Replace newlines and multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// countWords counts words in a string
func countWords(s string) int {
	return len(strings.Fields(s))
}
