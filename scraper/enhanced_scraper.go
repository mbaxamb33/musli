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
func extractImprovedContentSections(e *colly.HTMLElement, pageURL string, es *EnhancedScraper) {
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
	e.ForEach("h1, h2, h3, h4, h5, h6, p, div", func(_ int, el *colly.HTMLElement) {
		elements = append(elements, el)
	})

	// Process elements in document order
	for _, el := range elements {
		tagName := el.Name

		// Handle heading elements
		if len(tagName) == 2 && tagName[0] == 'h' && tagName[1] >= '1' && tagName[1] <= '6' {
			headingLevel := int(tagName[1] - '0')
			headingText := strings.TrimSpace(el.Text)

			if headingText == "" {
				continue // Skip empty headings
			}

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

			// Update the section hierarchy
			sectionHierarchy[headingLevel] = currentIndex

			// Clear any lower level sections (they're no longer active)
			for i := headingLevel + 1; i < len(sectionHierarchy); i++ {
				sectionHierarchy[i] = -1
			}

		} else if tagName == "p" || (tagName == "div" && !hasNestedBlockElements(el)) {
			// Handle content elements
			content := strings.TrimSpace(el.Text)

			// Skip if content is too short or empty
			if len(content) < 50 || countWords(content) < 10 {
				continue
			}

			// Skip if content looks like HTML
			if strings.Contains(content, "<") && strings.Contains(content, ">") {
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
				// Add content to the active section
				sections[activeSection].Content = append(sections[activeSection].Content, content)
			}
		}
	}

	// Process each section to create ContentItems
	for _, section := range sections {
		// Skip sections with no content
		if len(section.Content) == 0 {
			continue
		}

		// Combine all content for this section
		combinedContent := strings.Join(section.Content, "\n\n")

		// Create hash to detect duplicates
		hash := generateContentHash(section.Title, combinedContent)

		// Add to content items if not a duplicate
		es.mu.Lock()
		if !es.seenContent[hash] {
			es.seenContent[hash] = true

			// Create the content item
			es.ContentItems = append(es.ContentItems, ContentItem{
				URL:       pageURL,
				Title:     section.Title,
				Paragraph: combinedContent,
				Hash:      hash,
			})
		}
		es.mu.Unlock()
	}
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

// truncateString cuts a string to a maximum length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
