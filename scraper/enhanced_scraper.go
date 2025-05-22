package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

// ContentItem represents a piece of content with title and paragraph
type ContentItem struct {
	URL       string
	Title     string
	Paragraph string
	Hash      string // Unique hash to identify content
}

// EnhancedScraper extends the basic Scraper functionality for paragraph extraction
type EnhancedScraper struct {
	*Scraper
	ContentItems []ContentItem
	seenContent  map[string]bool // Track already seen content by hash
	mu           sync.Mutex
}

// NewEnhancedScraper creates a new enhanced scraper instance
func NewEnhancedScraper(baseURL string, maxDepth int) (*EnhancedScraper, error) {
	scraper, err := NewScraper(baseURL, maxDepth)
	if err != nil {
		return nil, err
	}

	return &EnhancedScraper{
		Scraper:      scraper,
		ContentItems: []ContentItem{},
		seenContent:  make(map[string]bool),
	}, nil
}

// ExtractTitleParagraphPairs extracts structured content from visited pages
func (es *EnhancedScraper) ExtractTitleParagraphPairs() error {
	// First make sure we've gathered all links
	if len(es.LinksToVisit) == 0 {
		if err := es.GatherLinks(); err != nil {
			return fmt.Errorf("error gathering links: %w", err)
		}
	}

	// Now extract content from all visited links
	baseDomain := getDomain(es.BaseURL)
	wwwDomain := "www." + baseDomain

	c := colly.NewCollector(
		colly.AllowedDomains(baseDomain, wwwDomain),
		colly.MaxDepth(es.MaxDepth),
	)

	// Track what URLs have been processed to avoid re-processing
	processedURLs := make(map[string]bool)

	// Define extractors for different types of content sections
	c.OnHTML("article, section, div.content, div.main, .content-area", func(e *colly.HTMLElement) {
		// Get the page URL
		pageURL := e.Request.URL.String()

		// Skip if this specific selector on this URL has already been processed
		selectorPath := pageURL + "#" + e.Name + "-" + e.Attr("class") + "-" + e.Attr("id")
		if processedURLs[selectorPath] {
			return
		}
		processedURLs[selectorPath] = true

		// Use the improved section extraction
		extractImprovedContentSections(e, pageURL, es)
	})

	// Visit each page in our link tree
	visitCount := 0
	for link := range es.LinksToVisit {
		err := c.Visit(link)
		if err != nil {
			fmt.Printf("Error visiting %s: %v\n", link, err)
			// Continue with other links
		}
		visitCount++

		// Debug info
		if visitCount%5 == 0 {
			fmt.Printf("Visited %d links, found %d content items so far\n",
				visitCount, len(es.ContentItems))
		}
	}

	c.Wait()

	// Apply additional de-duplication at the end
	es.removeDuplicateContent()

	return nil
}

// Section represents a logical section from a webpage
type Section struct {
	Title     string
	Level     int
	Content   []string
	StartTime time.Time
}

// extractImprovedContentSections implements the header + associated content approach
// extractImprovedContentSections implements the header + associated content approach
func extractImprovedContentSections(e *colly.HTMLElement, pageURL string, es *EnhancedScraper) {
	fmt.Println("======= STARTING CONTENT EXTRACTION FOR:", pageURL, "=======")

	// Track all sections we find
	var sections []Section

	// Keep track of the current active section hierarchy
	// The array index represents the heading level (0-6)
	// and the value is the index in the sections array
	sectionHierarchy := make([]int, 7)
	for i := range sectionHierarchy {
		sectionHierarchy[i] = -1 // -1 means no section at this level yet
	}

	// Default section for content before any headings
	defaultSection := Section{
		Title:     "Page Content",
		Level:     0,
		Content:   []string{},
		StartTime: time.Now(),
	}
	sections = append(sections, defaultSection)
	sectionHierarchy[0] = 0 // The default section is at index 0

	// Create an ordered list of all elements to process
	var elements []*colly.HTMLElement

	// First collect all relevant elements to ensure proper ordering
	fmt.Println("Collecting elements for processing...")
	elementCount := 0
	e.ForEach("h1, h2, h3, h4, h5, h6, p, div", func(_ int, el *colly.HTMLElement) {
		elements = append(elements, el)
		elementCount++
	})
	fmt.Printf("Collected %d elements to process\n", elementCount)

	// Process elements in document order
	for i, el := range elements {
		tagName := el.Name

		// Log element info
		rawText := el.Text
		trimmedText := strings.TrimSpace(rawText)
		fmt.Printf("\nElement #%d: Tag=%s, Length=%d\n", i, tagName, len(trimmedText))

		if len(trimmedText) > 20 {
			fmt.Printf("  Preview: '%s...'\n", trimmedText[:20])
		} else if len(trimmedText) > 0 {
			fmt.Printf("  Preview: '%s'\n", trimmedText)
		}

		// Handle heading elements
		if len(tagName) == 2 && tagName[0] == 'h' && tagName[1] >= '1' && tagName[1] <= '6' {
			headingLevel := int(tagName[1] - '0')
			headingText := strings.TrimSpace(el.Text)

			if headingText == "" {
				fmt.Println("  Skipping empty heading")
				continue // Skip empty headings
			}

			fmt.Printf("  Processing heading L%d: '%s'\n", headingLevel, headingText)

			// Create a new section for this heading
			newSection := Section{
				Title:     headingText,
				Level:     headingLevel,
				Content:   []string{},
				StartTime: time.Now(),
			}

			// Add section to our collection
			sections = append(sections, newSection)
			currentIndex := len(sections) - 1
			fmt.Printf("  Created section #%d with title: '%s'\n", currentIndex, headingText)

			// Update the section hierarchy
			sectionHierarchy[headingLevel] = currentIndex

			// Clear any lower level sections (they're no longer active)
			for i := headingLevel + 1; i < len(sectionHierarchy); i++ {
				sectionHierarchy[i] = -1
			}

		} else if tagName == "p" || (tagName == "div" && !hasNestedBlockElements(el)) {
			// Handle content elements
			content := strings.TrimSpace(el.Text)
			fmt.Printf("  Content element with %d chars, %d words\n",
				len(content), countWords(content))

			// Log full content for debugging
			if len(content) > 0 {
				// Print first 50 chars as preview
				preview := content
				if len(content) > 50 {
					preview = content[:50] + "..."
				}
				fmt.Printf("  Content preview: '%s'\n", preview)

				// Print last 20 chars to check for truncation
				if len(content) > 20 {
					lastChars := content[len(content)-20:]
					fmt.Printf("  Content ending with: '%s'\n", lastChars)

					// Check for ellipsis or other truncation signs
					if strings.HasSuffix(content, "...") {
						fmt.Printf("  WARNING: Content appears to be truncated\n")
					}
				}
			}

			// Skip if content is too short or empty
			if len(content) < 50 || countWords(content) < 10 {
				fmt.Println("  Skipping: Content too short")
				continue
			}

			// Skip if content looks like HTML
			if strings.Contains(content, "<") && strings.Contains(content, ">") {
				fmt.Println("  Skipping: Content appears to contain HTML")
				continue
			}

			// Find the active section to add this content to
			// Start from the highest heading level and work down
			activeSection := -1
			for level := 6; level >= 0; level-- {
				if sectionHierarchy[level] != -1 {
					activeSection = sectionHierarchy[level]
					break
				}
			}

			if activeSection != -1 {
				// Print which section we're adding to
				fmt.Printf("  Adding content to section #%d '%s'\n",
					activeSection, sections[activeSection].Title)

				// Check original HTML for debugging
				// Check original HTML for debugging
				htmlContent, htmlErr := el.DOM.Html()
				if htmlErr != nil {
					fmt.Printf("  Error getting HTML: %v\n", htmlErr)
				} else {
					htmlPreview := htmlContent
					if len(htmlContent) > 100 {
						htmlPreview = htmlContent[:100] + "..."
					}
					fmt.Printf("  Original HTML preview: %s\n", htmlPreview)
				}

				// Add content to the active section
				sections[activeSection].Content = append(sections[activeSection].Content, content)
			} else {
				fmt.Println("  WARNING: No active section found for content")
			}
		} else {
			fmt.Printf("  Skipping element: Not a heading, paragraph, or simple div\n")
		}
	}

	// Process each section to create ContentItems
	fmt.Printf("\nProcessing %d sections to create content items...\n", len(sections))
	for i, section := range sections {
		// Skip sections with no content
		if len(section.Content) == 0 {
			fmt.Printf("Section #%d '%s': Skipping - No content\n", i, section.Title)
			continue
		}

		// Combine all content for this section
		combinedContent := strings.Join(section.Content, "\n\n")
		fmt.Printf("Section #%d '%s': Combined %d content chunks, total length: %d\n",
			i, section.Title, len(section.Content), len(combinedContent))

		// Log first and last part of combined content
		if len(combinedContent) > 0 {
			firstPart := combinedContent
			if len(combinedContent) > 50 {
				firstPart = combinedContent[:50] + "..."
			}
			fmt.Printf("  Content starts with: '%s'\n", firstPart)

			if len(combinedContent) > 20 {
				lastPart := combinedContent[len(combinedContent)-20:]
				fmt.Printf("  Content ends with: '%s'\n", lastPart)

				if strings.HasSuffix(combinedContent, "...") {
					fmt.Printf("  WARNING: Final combined content appears to be truncated!\n")
				}
			}
		}

		// Create hash to detect duplicates
		hash := generateContentHash(section.Title, combinedContent)

		// Add to content items if not a duplicate
		es.mu.Lock()
		if !es.seenContent[hash] {
			es.seenContent[hash] = true

			// Create the content item
			newItem := ContentItem{
				URL:       pageURL,
				Title:     section.Title,
				Paragraph: combinedContent,
				Hash:      hash,
			}

			// Log the item being created
			fmt.Printf("  Adding content item: Title='%s', Length=%d\n",
				newItem.Title, len(newItem.Paragraph))

			// Check for potential truncation in final output
			if strings.HasSuffix(newItem.Paragraph, "...") {
				fmt.Printf("  CRITICAL: Final content item ends with '...', confirming truncation!\n")
			}

			es.ContentItems = append(es.ContentItems, newItem)
		} else {
			fmt.Printf("  Skipping duplicate content with hash: %s\n", hash)
		}
		es.mu.Unlock()
	}

	fmt.Printf("======= COMPLETED CONTENT EXTRACTION FOR %s: Added %d content items =======\n\n",
		pageURL, len(es.ContentItems))
}

// Run executes the complete enhanced scraping process
func (es *EnhancedScraper) Run() error {
	// First run the basic scraping
	err := es.Scraper.Run()
	if err != nil {
		return fmt.Errorf("error in basic scraping: %w", err)
	}

	// Then extract content
	err = es.ExtractTitleParagraphPairs()
	if err != nil {
		return fmt.Errorf("error extracting content: %w", err)
	}

	return nil
}

// RemoveDuplicateContent filters out any duplicate content
func (es *EnhancedScraper) removeDuplicateContent() {
	es.mu.Lock()
	defer es.mu.Unlock()

	// Create a map to track unique content
	seen := make(map[string]bool)
	var uniqueItems []ContentItem

	// Only keep unique content based on hash
	for _, item := range es.ContentItems {
		if !seen[item.Hash] {
			seen[item.Hash] = true
			uniqueItems = append(uniqueItems, item)
		}
	}

	// Replace with deduplicated list
	es.ContentItems = uniqueItems
	fmt.Printf("Removed %d duplicate items\n", len(es.ContentItems)-len(uniqueItems))
}

// Helper functions

// generateContentHash creates a unique hash for content to detect duplicates
func generateContentHash(title, content string) string {
	// Normalize content before hashing
	normalizedTitle := cleanText(title)
	normalizedContent := cleanText(content)

	// Create hash from combined content
	h := sha256.New()
	h.Write([]byte(normalizedTitle + "|" + normalizedContent))
	return hex.EncodeToString(h.Sum(nil))
}

// hasNestedBlockElements checks if an element contains common block-level elements
func hasNestedBlockElements(e *colly.HTMLElement) bool {
	blockElements := []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6",
		"ul", "ol", "table", "section", "article", "aside", "nav"}

	for _, tag := range blockElements {
		if e.DOM.Find(tag).Length() > 0 {
			return true
		}
	}

	return false
}

// countWords counts words in a string
func countWords(s string) int {
	return len(strings.Fields(s))
}

// Clean up text by removing excess whitespace
func cleanText(s string) string {
	// Replace newlines and multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// // truncateString cuts a string to a maximum length with ellipsis
// func truncateString(s string, maxLen int) string {
// 	if len(s) <= maxLen {
// 		return s
// 	}
// 	return s[:maxLen-3] + "..."
// }
