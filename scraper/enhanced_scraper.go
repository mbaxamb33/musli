package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"

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

		// Find headings and content sections
		extractContentSections(e, pageURL, es)
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

// extractContentSections finds title-paragraph pairs within an HTML element
func extractContentSections(e *colly.HTMLElement, pageURL string, es *EnhancedScraper) {
	// Track headings and their paragraphs at different levels
	type headingContent struct {
		title string
		depth int
	}

	var currentHeadings []headingContent

	// Process all child elements in sequence to maintain structure relationship
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		tagName := child.Name

		// Check if this is a heading element
		isHeading := false
		headingLevel := 0

		if len(tagName) == 2 && tagName[0] == 'h' && tagName[1] >= '1' && tagName[1] <= '6' {
			isHeading = true
			headingLevel = int(tagName[1] - '0')
		}

		if isHeading {
			// Found a heading, adjust current headings stack
			title := strings.TrimSpace(child.Text)
			if title == "" {
				return
			}

			// Remove any headings at same or deeper level
			for i := 0; i < len(currentHeadings); i++ {
				if currentHeadings[i].depth >= headingLevel {
					currentHeadings = currentHeadings[:i]
					break
				}
			}

			// Add this heading
			currentHeadings = append(currentHeadings, headingContent{
				title: title,
				depth: headingLevel,
			})
		} else if tagName == "p" ||
			tagName == "div" && !hasBlockElements(child) ||
			tagName == "span" && len(strings.TrimSpace(child.Text)) > 100 {
			// Found a paragraph or text-containing div
			paragraph := strings.TrimSpace(child.Text)

			// Skip if paragraph is too short or empty
			if len(paragraph) == 0 || countWords(paragraph) < 10 {
				return
			}

			// Skip if paragraph contains HTML or looks like a link
			if strings.Contains(paragraph, "<") && strings.Contains(paragraph, ">") {
				return
			}

			// Use the most specific (deepest) heading as the title
			title := "Untitled Content"
			if len(currentHeadings) > 0 {
				title = currentHeadings[len(currentHeadings)-1].title
			}

			// Create hash from both title and content to identify duplicates
			hash := generateContentHash(title, paragraph)

			// Skip if we've seen this exact content before
			es.mu.Lock()
			if !es.seenContent[hash] {
				es.seenContent[hash] = true

				// Add to content items
				es.ContentItems = append(es.ContentItems, ContentItem{
					URL:       pageURL,
					Title:     title,
					Paragraph: paragraph,
					Hash:      hash,
				})
			}
			es.mu.Unlock()
		}
	})
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
	fmt.Printf("Removed %d duplicate items\n", len(seen)-len(uniqueItems))
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

// Check if an element contains block-level elements
func hasBlockElements(e *colly.HTMLElement) bool {
	blockElements := []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "table", "section", "article"}

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

// truncateString cuts a string to a maximum length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
