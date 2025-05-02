package docscraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/unidoc/unioffice/document"
)

// ContentType represents the type of content extracted
type ContentType int

const (
	ContentTypeParagraph ContentType = iota
	ContentTypeOrderedList
	ContentTypeUnorderedList
	ContentTypeTable
)

// ContentItem represents a piece of content extracted from the document
type ContentItem struct {
	Type         ContentType
	Heading      string   // Parent heading
	HeadingPath  []string // Full path of headings (for hierarchy)
	HeadingLevel int      // Level of the parent heading (1-based)
	Content      string   // The actual content
	ListItems    []string // For list types only
	Hash         string   // Unique hash to identify content
}

// DocTree represents the document structure as a tree
type DocTree struct {
	Title    string
	Level    int
	Content  []ContentItem
	Children map[string]*DocTree
}

// Section represents a complete document section with heading and all its content
type Section struct {
	Title       string
	Level       int
	HeadingPath []string
	Paragraphs  []ContentItem
	Lists       []ListGroup
	Tables      []TableContent
	Subsections []*Section
}

// ListGroup represents a group of related list items
type ListGroup struct {
	Type  ContentType // OrderedList or UnorderedList
	Items []string
	Hash  string
}

// TableContent represents a table extracted from the document
type TableContent struct {
	Headers []string
	Rows    [][]string
	Hash    string
}

// DocumentScraper handles extraction from Word documents
type DocumentScraper struct {
	FilePath     string
	ContentItems []ContentItem
	Tree         *DocTree
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
		Tree: &DocTree{
			Title:    "Root",
			Level:    0,
			Content:  []ContentItem{},
			Children: make(map[string]*DocTree),
		},
		seenContent: make(map[string]bool),
	}, nil
}

// Run executes the complete document scraping process
func (ds *DocumentScraper) Run() error {
	// Open the document
	doc, err := document.Open(ds.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	// Extract document title from filename if not available in metadata
	docTitle := doc.CoreProperties.Title()
	if docTitle == "" {
		docTitle = strings.TrimSuffix(strings.Split(ds.FilePath, "/")[len(strings.Split(ds.FilePath, "/"))-1], ".docx")
	}
	ds.Tree.Title = docTitle

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

			// Update document tree
			ds.updateDocTree(currentHeadingPath, headingText, headingLevel)
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

			// Detect if it's a list item
			isListItem, listType := ds.detectListItem(para, paraText)

			if isListItem {
				// Process as part of a list - this would need to group consecutive list items
				// For simplicity in this example, we're treating each as a separate item
				// A more complex implementation would group consecutive list items into a single ContentItem
				ds.addContentItem(ContentItem{
					Type:         listType,
					Heading:      headingText,
					HeadingPath:  append([]string{}, currentHeadingPath...),
					HeadingLevel: currentHeadingLevel,
					Content:      content,
					ListItems:    []string{content},
				})
			} else {
				// Process as normal paragraph
				ds.addContentItem(ContentItem{
					Type:         ContentTypeParagraph,
					Heading:      headingText,
					HeadingPath:  append([]string{}, currentHeadingPath...),
					HeadingLevel: currentHeadingLevel,
					Content:      content,
				})
			}
		}
	}

	return nil
}

// detectListItem determines if a paragraph is a list item
func (ds *DocumentScraper) detectListItem(para document.Paragraph, text string) (bool, ContentType) {
	// Check for numbering property - this needs to be adapted to the actual API
	// The exact implementation will depend on how numbering is stored in the document package
	if para.Properties().X() != nil && para.Properties().X().NumPr != nil {
		// If we have numbering properties, it's likely a list item
		// We would need to check the numbering definition to determine if ordered or unordered
		// For simplicity, assume ordered list if we find numbering properties
		return true, ContentTypeOrderedList
	}

	// Alternative detection: bullet symbols at beginning of text
	if strings.HasPrefix(text, "•") || strings.HasPrefix(text, "○") ||
		strings.HasPrefix(text, "▪") || strings.HasPrefix(text, "-") {
		return true, ContentTypeUnorderedList
	}

	// Check for numbered format like "1.", "a.", "(i)", etc.
	numPattern := regexp.MustCompile(`^(\d+\.|[a-z]\.|[ivxlcdm]+\.|[IVXLCDM]+\.|\([0-9a-zA-Z]+\))\s`)
	if numPattern.MatchString(text) {
		return true, ContentTypeOrderedList
	}

	return false, ContentTypeParagraph
}

// updateDocTree updates the document tree structure with new heading
func (ds *DocumentScraper) updateDocTree(headingPath []string, headingText string, level int) {
	// Start at the root
	current := ds.Tree

	// Navigate/create path to the current heading level
	for i, heading := range headingPath[:level] {
		if i == len(headingPath)-1 {
			// This is the heading we're trying to add
			if _, exists := current.Children[heading]; !exists {
				current.Children[heading] = &DocTree{
					Title:    heading,
					Level:    i + 1,
					Content:  []ContentItem{},
					Children: make(map[string]*DocTree),
				}
			}
		} else {
			// This is a parent heading, navigate to it
			if next, exists := current.Children[heading]; exists {
				current = next
			} else {
				// Create missing parent
				current.Children[heading] = &DocTree{
					Title:    heading,
					Level:    i + 1,
					Content:  []ContentItem{},
					Children: make(map[string]*DocTree),
				}
				current = current.Children[heading]
			}
		}
	}
}

// addContentItem adds a content item to the document scraper and updates the tree
func (ds *DocumentScraper) addContentItem(item ContentItem) {
	// Generate hash for deduplication
	item.Hash = generateContentHash(item.Heading, item.Content)

	// Skip if we've seen this content before
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if !ds.seenContent[item.Hash] {
		ds.seenContent[item.Hash] = true
		ds.ContentItems = append(ds.ContentItems, item)

		// Also update the document tree
		current := ds.Tree
		for i, heading := range item.HeadingPath {
			if next, exists := current.Children[heading]; exists {
				current = next
			} else {
				// Create missing heading node
				current.Children[heading] = &DocTree{
					Title:    heading,
					Level:    i + 1,
					Content:  []ContentItem{},
					Children: make(map[string]*DocTree),
				}
				current = current.Children[heading]
			}
		}

		// Add content to the appropriate section
		current.Content = append(current.Content, item)
	}
}

// PrintDocTree prints the document tree structure
func PrintDocTree(node *DocTree, indent string) {
	if node == nil {
		return
	}

	// Print current node
	fmt.Printf("%s%s (Level %d)\n", indent, node.Title, node.Level)

	// Print content summary
	if len(node.Content) > 0 {
		fmt.Printf("%s  Content: %d items\n", indent, len(node.Content))
	}

	// Print children
	for _, child := range node.Children {
		PrintDocTree(child, indent+"  ")
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

// SaveTextRepresentation saves a plain text representation of the extracted content
func (ds *DocumentScraper) SaveTextRepresentation(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write document title
	fmt.Fprintf(file, "# %s\n\n", ds.Tree.Title)

	// Write structured content
	ds.writeDocTreeContent(file, ds.Tree, 0)

	return nil
}

// writeDocTreeContent writes the document tree content recursively
func (ds *DocumentScraper) writeDocTreeContent(w io.Writer, node *DocTree, level int) {
	if node == nil {
		return
	}

	// Write heading with appropriate level markers
	if level > 0 {
		fmt.Fprintf(w, "%s %s\n\n", strings.Repeat("#", level), node.Title)
	}

	// Write content items
	for _, item := range node.Content {
		switch item.Type {
		case ContentTypeParagraph:
			fmt.Fprintf(w, "%s\n\n", item.Content)

		case ContentTypeOrderedList:
			fmt.Fprintln(w, "Ordered list:")
			for i, listItem := range item.ListItems {
				fmt.Fprintf(w, "%d. %s\n", i+1, listItem)
			}
			fmt.Fprintln(w)

		case ContentTypeUnorderedList:
			fmt.Fprintln(w, "Unordered list:")
			for _, listItem := range item.ListItems {
				fmt.Fprintf(w, "* %s\n", listItem)
			}
			fmt.Fprintln(w)
		}
	}

	// Process children
	for _, child := range node.Children {
		ds.writeDocTreeContent(w, child, level+1)
	}
}

// GetGroupedContent returns the document content grouped by sections
func (ds *DocumentScraper) GetGroupedContent() []*Section {
	var sections []*Section
	ds.buildSectionTree(ds.Tree, nil, &sections)
	return sections
}

// buildSectionTree converts the DocTree to a Section-based structure
func (ds *DocumentScraper) buildSectionTree(node *DocTree, parentPath []string, sections *[]*Section) {
	if node == nil {
		return
	}

	// Skip the root node
	if node.Level > 0 {
		// Create a new section
		section := &Section{
			Title:       node.Title,
			Level:       node.Level,
			HeadingPath: append(append([]string{}, parentPath...), node.Title),
			Subsections: []*Section{},
		}

		// Group content by type
		for _, item := range node.Content {
			switch item.Type {
			case ContentTypeParagraph:
				section.Paragraphs = append(section.Paragraphs, item)

			case ContentTypeOrderedList, ContentTypeUnorderedList:
				// Check if this list item should be part of an existing list group
				foundExisting := false
				for i, list := range section.Lists {
					if list.Type == item.Type {
						// This is a potential match - if it's part of the same logical list
						// In a real implementation, you'd need a better way to identify related list items
						section.Lists[i].Items = append(section.Lists[i].Items, item.Content)
						foundExisting = true
						break
					}
				}

				if !foundExisting {
					// Create a new list group
					section.Lists = append(section.Lists, ListGroup{
						Type:  item.Type,
						Items: []string{item.Content},
						Hash:  item.Hash,
					})
				}

			case ContentTypeTable:
				// You would implement table grouping here
			}
		}

		// Add to sections list
		*sections = append(*sections, section)

		// Update parent path for children
		parentPath = append(parentPath, node.Title)
	}

	// Process children
	for _, child := range node.Children {
		ds.buildSectionTree(child, parentPath, sections)
	}
}

// GetContentByHeading returns all content associated with a specific heading
func (ds *DocumentScraper) GetContentByHeading(heading string) ([]ContentItem, error) {
	var result []ContentItem
	found := false

	// Helper function to search the tree
	var searchTree func(node *DocTree)
	searchTree = func(node *DocTree) {
		if node == nil {
			return
		}

		// Check if this node matches
		if node.Title == heading {
			result = append(result, node.Content...)
			found = true
			return
		}

		// Check children
		for _, child := range node.Children {
			searchTree(child)
		}
	}

	searchTree(ds.Tree)

	if !found {
		return nil, fmt.Errorf("heading not found: %s", heading)
	}

	return result, nil
}

// GetSectionWithContext returns a section with its content and contextual information
func (ds *DocumentScraper) GetSectionWithContext(headingPath []string) (*Section, error) {
	if len(headingPath) == 0 {
		return nil, fmt.Errorf("empty heading path")
	}

	// Find the node that matches the path
	current := ds.Tree
	for _, heading := range headingPath {
		found := false
		for title, child := range current.Children {
			if title == heading {
				current = child
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("heading path not found: %v", headingPath)
		}
	}

	// Create a section from this node
	section := &Section{
		Title:       current.Title,
		Level:       current.Level,
		HeadingPath: headingPath,
		Subsections: []*Section{},
	}

	// Group content
	for _, item := range current.Content {
		switch item.Type {
		case ContentTypeParagraph:
			section.Paragraphs = append(section.Paragraphs, item)
		case ContentTypeOrderedList, ContentTypeUnorderedList:
			// Check if this list item belongs to an existing list
			foundExisting := false
			for i, list := range section.Lists {
				if list.Type == item.Type {
					section.Lists[i].Items = append(section.Lists[i].Items, item.Content)
					foundExisting = true
					break
				}
			}

			if !foundExisting {
				section.Lists = append(section.Lists, ListGroup{
					Type:  item.Type,
					Items: []string{item.Content},
					Hash:  item.Hash,
				})
			}
		}
	}

	// Add subsections
	for title, child := range current.Children {
		subsection := &Section{
			Title:       title,
			Level:       child.Level,
			HeadingPath: append(append([]string{}, headingPath...), title),
		}
		section.Subsections = append(section.Subsections, subsection)
	}

	return section, nil
}

// GetSummary returns a concise summary of the document
func (ds *DocumentScraper) GetSummary() map[string]interface{} {
	// Create a summary structure
	summary := map[string]interface{}{
		"title":            ds.Tree.Title,
		"sections":         len(ds.getHeadingCount()),
		"totalParagraphs":  ds.countContentByType(ContentTypeParagraph),
		"totalLists":       ds.countContentByType(ContentTypeOrderedList) + ds.countContentByType(ContentTypeUnorderedList),
		"totalTables":      ds.countContentByType(ContentTypeTable),
		"headingStructure": ds.getHeadingCount(),
		"topLevelSections": ds.getTopLevelSectionNames(),
	}

	return summary
}

// getHeadingCount returns a map with counts of headings at each level
func (ds *DocumentScraper) getHeadingCount() map[int]int {
	result := make(map[int]int)

	// Helper function to count headings
	var countHeadings func(node *DocTree)
	countHeadings = func(node *DocTree) {
		if node == nil {
			return
		}

		if node.Level > 0 {
			result[node.Level]++
		}

		for _, child := range node.Children {
			countHeadings(child)
		}
	}

	countHeadings(ds.Tree)
	return result
}

// getTopLevelSectionNames returns the names of top-level sections
func (ds *DocumentScraper) getTopLevelSectionNames() []string {
	var names []string

	for title := range ds.Tree.Children {
		names = append(names, title)
	}

	return names
}

// countContentByType counts content items of a specific type
func (ds *DocumentScraper) countContentByType(contentType ContentType) int {
	count := 0
	for _, item := range ds.ContentItems {
		if item.Type == contentType {
			count++
		}
	}
	return count
}

// ExportSectionContent exports a section and its content as formatted text
func (ds *DocumentScraper) ExportSectionContent(w io.Writer, headingPath []string) error {
	section, err := ds.GetSectionWithContext(headingPath)
	if err != nil {
		return err
	}

	// Write section title
	fmt.Fprintf(w, "%s %s\n\n", strings.Repeat("#", section.Level), section.Title)

	// Write paragraph content
	for _, para := range section.Paragraphs {
		fmt.Fprintf(w, "%s\n\n", para.Content)
	}

	// Write lists
	for _, list := range section.Lists {
		if list.Type == ContentTypeOrderedList {
			fmt.Fprintln(w, "Ordered list:")
			for i, item := range list.Items {
				fmt.Fprintf(w, "%d. %s\n", i+1, item)
			}
			fmt.Fprintln(w)
		} else {
			fmt.Fprintln(w, "Unordered list:")
			for _, item := range list.Items {
				fmt.Fprintf(w, "* %s\n", item)
			}
			fmt.Fprintln(w)
		}
	}

	// Write subsection titles
	if len(section.Subsections) > 0 {
		fmt.Fprintln(w, "Subsections:")
		for _, sub := range section.Subsections {
			fmt.Fprintf(w, "- %s\n", sub.Title)
		}
		fmt.Fprintln(w)
	}

	return nil
}

// ExportAllSectionsContent exports all sections with their content
func (ds *DocumentScraper) ExportAllSectionsContent(w io.Writer) error {
	sections := ds.GetGroupedContent()

	// Write document title
	fmt.Fprintf(w, "# %s\n\n", ds.Tree.Title)

	// Write sections in order
	for _, section := range sections {
		// Write section heading
		fmt.Fprintf(w, "%s %s\n\n",
			strings.Repeat("#", section.Level+1),
			section.Title)

		// Write paragraph content
		for _, para := range section.Paragraphs {
			fmt.Fprintf(w, "%s\n\n", para.Content)
		}

		// Write lists
		for _, list := range section.Lists {
			if list.Type == ContentTypeOrderedList {
				for i, item := range list.Items {
					fmt.Fprintf(w, "%d. %s\n", i+1, item)
				}
				fmt.Fprintln(w)
			} else {
				for _, item := range list.Items {
					fmt.Fprintf(w, "* %s\n", item)
				}
				fmt.Fprintln(w)
			}
		}
	}

	return nil
}

// The remaining functions for content analysis remain the same...

// splitIntoSentences splits text into sentences
func splitIntoSentences(text string) []string {
	// Simple sentence splitting - in a real implementation,
	// this would be more sophisticated to handle abbreviations, etc.
	regex := regexp.MustCompile(`[.!?]\s+`)
	sentences := regex.Split(text, -1)

	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			result = append(result, sentence)
		}
	}

	return result
}

// extractKeywords extracts important keywords from text
func extractKeywords(text string) []string {
	// In a real implementation, this would use TF-IDF or other techniques
	// For this example, we'll use a simple approach

	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove common punctuation
	text = regexp.MustCompile(`[.,;:!?()]`).ReplaceAllString(text, "")

	// Split into words
	words := strings.Fields(text)

	// Count word frequencies
	wordCounts := make(map[string]int)
	for _, word := range words {
		// Skip very short words and common English stopwords
		if len(word) <= 2 || isStopword(word) {
			continue
		}
		wordCounts[word]++
	}

	// Sort by frequency
	type wordFreq struct {
		word  string
		count int
	}

	var wordFreqs []wordFreq
	for word, count := range wordCounts {
		wordFreqs = append(wordFreqs, wordFreq{word, count})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count > wordFreqs[j].count
	})

	// Take top 10 keywords
	count := min(10, len(wordFreqs))
	keywords := make([]string, count)
	for i := 0; i < count; i++ {
		keywords[i] = wordFreqs[i].word
	}

	return keywords
}

// extractKeyTerms extracts important terms (which may be multi-word)
func extractKeyTerms(text string, count int) []string {
	// In a real implementation, this would use techniques like noun phrase extraction
	// For this example, we'll just use the most frequent words

	keywords := extractKeywords(text)
	if len(keywords) <= count {
		return keywords
	}
	return keywords[:count]
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// isStopword checks if a word is a common English stopword
func isStopword(word string) bool {
	stopwords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "but": true,
		"if": true, "or": true, "because": true, "as": true, "until": true,
		"while": true, "of": true, "at": true, "by": true, "for": true,
		"with": true, "about": true, "against": true, "between": true,
		"into": true, "through": true, "during": true, "before": true,
		"after": true, "above": true, "below": true, "to": true,
		"from": true, "up": true, "down": true, "in": true, "out": true,
		"on": true, "off": true, "over": true, "under": true, "again": true,
		"further": true, "then": true, "once": true, "here": true,
		"there": true, "when": true, "where": true, "why": true, "how": true,
		"all": true, "any": true, "both": true, "each": true, "few": true,
		"more": true, "most": true, "other": true, "some": true, "such": true,
		"no": true, "nor": true, "not": true, "only": true, "own": true,
		"same": true, "so": true, "than": true, "too": true, "very": true,
		"s": true, "t": true, "can": true, "will": true, "just": true,
		"don": true, "should": true, "now": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "could": true, "would": true,
		"shall": true, "may": true, "might": true, "must": true,
	}

	return stopwords[word]
}
