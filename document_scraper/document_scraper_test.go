package docscraper

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"
// 	"time"

// 	_ "github.com/lib/pq" // PostgreSQL driver
// 	db "github.com/mbaxamb3/nusli/db/sqlc"
// )

// // Define missing enum value
// const (
// 	// DatasourceTypeDocument is a type constant for document datasources
// 	DatasourceTypeDocument = "document"
// )

// // TestDocumentScraper demonstrates extracting content from a Word document
// func TestDocumentScraper(t *testing.T) {
// 	// Skip this test in automated test runs, run it manually
// 	if testing.Short() {
// 		t.Skip("Skipping integration test in short mode")
// 	}

// 	// Define path to the Word document
// 	filePath := "./samples/sample_document.docx"

// 	// Create document scraper
// 	fmt.Printf("📄 Creating document scraper for %s\n", filePath)
// 	scraper, err := NewDocumentScraper(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to create document scraper: %v", err)
// 	}

// 	// Run the scraper
// 	fmt.Println("🚀 Running document scraper to extract content...")
// 	err = scraper.Run()
// 	if err != nil {
// 		t.Fatalf("Document scraping failed: %v", err)
// 	}
// 	fmt.Printf("✅ Extracted %d unique content items\n", len(scraper.ContentItems))

// 	// Print the document tree structure
// 	fmt.Println("\nDocument Structure:")
// 	PrintDocTree(scraper.Tree, "")

// 	// Save a text representation of the extracted content
// 	outputPath := "./output/extracted_content.md"
// 	fmt.Printf("💾 Saving extracted content to %s\n", outputPath)
// 	err = scraper.SaveTextRepresentation(outputPath)
// 	if err != nil {
// 		t.Fatalf("Failed to save text representation: %v", err)
// 	}
// 	fmt.Println("✅ Saved text representation")

// 	// Get grouped content
// 	fmt.Println("\nGetting grouped content...")
// 	sections := scraper.GetGroupedContent()
// 	fmt.Printf("Found %d sections\n", len(sections))

// 	// Print section summaries
// 	fmt.Println("\nSection Summaries:")
// 	for _, section := range sections {
// 		fmt.Printf("- %s (Level %d): %d paragraphs, %d lists\n",
// 			section.Title, section.Level, len(section.Paragraphs), len(section.Lists))
// 	}

// 	// Get document overview
// 	fmt.Println("\nDocument Overview:")
// 	overview := getDocumentOverview(scraper)
// 	fmt.Printf("Title: %s\n", overview["title"])
// 	fmt.Printf("Sections: %d\n", overview["sectionCount"])
// 	fmt.Printf("Total Paragraphs: %d\n", overview["totalParagraphs"])
// 	fmt.Printf("Total Words: %d\n", overview["totalWords"])
// 	fmt.Printf("Main Topics: %v\n", overview["mainTopics"])

// 	// Export all sections with content
// 	allSectionsPath := "./output/all_sections.md"
// 	fmt.Printf("\n💾 Saving all sections to %s\n", allSectionsPath)
// 	file, err := os.Create(allSectionsPath)
// 	if err != nil {
// 		t.Fatalf("Failed to create output file: %v", err)
// 	}
// 	defer file.Close()
// 	err = scraper.ExportAllSectionsContent(file)
// 	if err != nil {
// 		t.Fatalf("Failed to export all sections: %v", err)
// 	}
// 	fmt.Println("✅ Exported all sections")

// 	// Try to extract thematic blocks from a section
// 	if len(sections) > 0 {
// 		fmt.Println("\nExtracting thematic blocks from first section...")
// 		blocks, err := groupParagraphsByTheme(scraper, sections[0].HeadingPath)
// 		if err != nil {
// 			t.Fatalf("Failed to extract thematic blocks: %v", err)
// 		}

// 		fmt.Printf("Found %d thematic blocks\n", len(blocks))
// 		for i, block := range blocks {
// 			fmt.Printf("Block %d: Theme: %s\n", i+1, block.Theme)
// 			fmt.Printf("  Summary: %s\n", block.Summary)
// 			fmt.Printf("  Key Terms: %v\n", block.KeyTerms)
// 			fmt.Printf("  Paragraphs: %d\n", len(block.Paragraphs))
// 		}
// 	}
// }

// type TextBlock struct {
// 	Paragraphs []ContentItem
// 	Theme      string
// 	KeyTerms   []string
// 	Summary    string
// }

// // TestPopulateDocumentDB demonstrates extracting content from a Word document and adding it to the database
// func TestPopulateDocumentDB(t *testing.T) {
// 	// Skip this test in automated test runs, run it manually
// 	if testing.Short() {
// 		t.Skip("Skipping integration test in short mode")
// 	}

// 	// 1. Set up database connection
// 	connStr := "postgresql://root:secret@localhost:5432/musli?sslmode=disable"
// 	dbConn, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to database: %v", err)
// 	}
// 	defer dbConn.Close()

// 	// Test the connection
// 	err = dbConn.Ping()
// 	if err != nil {
// 		t.Fatalf("Failed to ping database: %v", err)
// 	}
// 	fmt.Println("✅ Successfully connected to database")

// 	// Create a queries object
// 	queries := db.New(dbConn)
// 	ctx := context.Background()

// 	// 2. Create document scraper
// 	filePath := "./samples/company_report.docx"
// 	fmt.Printf("📄 Creating document scraper for %s\n", filePath)

// 	docScraper, err := NewDocumentScraper(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to create document scraper: %v", err)
// 	}

// 	// 3. Run the document scraper
// 	fmt.Println("🚀 Running document scraper to extract content...")
// 	err = docScraper.Run()
// 	if err != nil {
// 		t.Fatalf("Document scraping failed: %v", err)
// 	}
// 	fmt.Printf("✅ Extracted %d unique content items\n", len(docScraper.ContentItems))

// 	// 4. Print the document structure
// 	fmt.Println("\nDocument Structure:")
// 	PrintDocTree(docScraper.Tree, "")

// 	// 5. Create a user
// 	fmt.Println("👤 Creating user record...")
// 	user, err := queries.CreateUser(ctx, db.CreateUserParams{
// 		Username: "analyst_" + randomString(6),
// 		Password: "secure_" + randomString(10),
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to create user: %v", err)
// 	}
// 	fmt.Printf("✅ Created user with ID: %d\n", user.UserID)

// 	// 6. Create a company record
// 	fmt.Println("💾 Creating company record...")
// 	companyName := "Unknown"

// 	// Try to extract company name from document title
// 	if docScraper.Tree.Title != "" {
// 		companyName = docScraper.Tree.Title
// 	}

// 	companyParams := db.CreateCompanyParams{
// 		UserID:      user.UserID,
// 		CompanyName: companyName,
// 		Industry:    sql.NullString{String: "Extracted from document", Valid: true},
// 		Website:     sql.NullString{String: "", Valid: false},
// 		Description: sql.NullString{String: "Company extracted from document content", Valid: true},
// 	}

// 	company, err := queries.CreateCompany(ctx, companyParams)
// 	if err != nil {
// 		t.Fatalf("Failed to create company: %v", err)
// 	}
// 	fmt.Printf("✅ Created company with ID: %d\n", company.CompanyID)

// 	// 7. Create datasource record for the document
// 	fmt.Println("📊 Creating datasource record...")

// 	// Read file contents
// 	fileData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to read file: %v", err)
// 	}

// 	datasourceParams := db.CreateDatasourceParams{
// 		SourceType: DatasourceTypeDocument,
// 		Link:       sql.NullString{String: "", Valid: false},
// 		FileName:   sql.NullString{String: filepath.Base(filePath), Valid: true},
// 		FileData:   fileData,
// 	}

// 	datasource, err := queries.CreateDatasource(ctx, datasourceParams)
// 	if err != nil {
// 		t.Fatalf("Failed to create datasource: %v", err)
// 	}
// 	fmt.Printf("✅ Created datasource with ID: %d\n", datasource.DatasourceID)

// 	// 8. Associate datasource with company
// 	fmt.Println("🔗 Associating datasource with company...")
// 	err = queries.AssociateDatasourceWithCompany(ctx, db.AssociateDatasourceWithCompanyParams{
// 		CompanyID:    company.CompanyID,
// 		DatasourceID: datasource.DatasourceID,
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to associate datasource with company: %v", err)
// 	}
// 	fmt.Println("✅ Associated datasource with company")

// 	// 9. Create paragraphs from extracted content
// 	fmt.Println("📝 Creating paragraphs from extracted content...")
// 	addedParagraphs := 0

// 	// Get grouped content for better organization
// 	sections := docScraper.GetGroupedContent()
// 	fmt.Printf("Found %d sections to process\n", len(sections))

// 	// Keep track of content we've already added to the database
// 	seenHashes := make(map[string]bool)

// 	// Process each section
// 	for _, section := range sections {
// 		fmt.Printf("Processing section: %s\n", section.Title)

// 		// Add paragraphs
// 		for _, para := range section.Paragraphs {
// 			// Skip if we've already added this exact content
// 			if seenHashes[para.Hash] {
// 				continue
// 			}
// 			seenHashes[para.Hash] = true

// 			// Extract main idea (could be implemented with NLP in a real system)
// 			mainIdea := ""
// 			if len(para.Content) > 0 && countWords(para.Content) >= 10 {
// 				// For demo, just use first sentence as main idea
// 				sentences := splitIntoSentences(para.Content)
// 				if len(sentences) > 0 {
// 					mainIdea = strings.TrimSpace(sentences[0])
// 				}
// 			}

// 			paragraphParams := db.CreateParagraphParams{
// 				DatasourceID: datasource.DatasourceID,
// 				Title:        sql.NullString{String: cleanText(section.Title), Valid: true},
// 				MainIdea:     sql.NullString{String: mainIdea, Valid: mainIdea != ""},
// 				Content:      cleanText(para.Content),
// 			}

// 			_, err = queries.CreateParagraph(ctx, paragraphParams)
// 			if err != nil {
// 				log.Printf("Error creating paragraph: %v", err)
// 				continue
// 			}
// 			addedParagraphs++

// 			// Print sample of what we're adding
// 			if addedParagraphs <= 5 {
// 				fmt.Printf("  Added: Title: %s\n  Paragraph: %s\n\n",
// 					truncateString(section.Title, 40),
// 					truncateString(para.Content, 100))
// 			}
// 		}

// 		// Add lists as paragraphs too
// 		for _, list := range section.Lists {
// 			// Create a single paragraph from the list items
// 			listContent := ""
// 			for i, item := range list.Items {
// 				if list.Type == ContentTypeOrderedList {
// 					listContent += fmt.Sprintf("%d. %s\n", i+1, item)
// 				} else {
// 					listContent += fmt.Sprintf("• %s\n", item)
// 				}
// 			}

// 			if listContent == "" {
// 				continue
// 			}

// 			// Generate a hash for the list content
// 			listHash := generateContentHash(section.Title, listContent)

// 			// Skip if we've already added this content
// 			if seenHashes[listHash] {
// 				continue
// 			}
// 			seenHashes[listHash] = true

// 			// Add to database
// 			paragraphParams := db.CreateParagraphParams{
// 				DatasourceID: datasource.DatasourceID,
// 				Title:        sql.NullString{String: cleanText(section.Title) + " - List", Valid: true},
// 				MainIdea:     sql.NullString{String: "List items", Valid: true},
// 				Content:      cleanText(listContent),
// 			}

// 			_, err = queries.CreateParagraph(ctx, paragraphParams)
// 			if err != nil {
// 				log.Printf("Error creating list paragraph: %v", err)
// 				continue
// 			}
// 			addedParagraphs++
// 		}
// 	}

// 	fmt.Printf("✅ Added %d unique paragraphs to database\n", addedParagraphs)

// 	fmt.Println("✅ Successfully populated database with document content!")
// }

// // Helper function to generate random strings
// func randomString(n int) string {
// 	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
// 		time.Sleep(1 * time.Nanosecond) // To ensure uniqueness
// 	}
// 	return string(b)
// }

// // truncateString cuts a string to a maximum length with ellipsis
// func truncateString(s string, maxLen int) string {
// 	if len(s) <= maxLen {
// 		return s
// 	}
// 	return s[:maxLen-3] + "..."
// }

// // Implementation of GetDocumentOverview
// func getDocumentOverview(ds *DocumentScraper) map[string]interface{} {
// 	// Count total words
// 	totalWords := 0
// 	for _, item := range ds.ContentItems {
// 		if item.Type == ContentTypeParagraph {
// 			totalWords += countWords(item.Content)
// 		}
// 	}

// 	// Get section summaries
// 	var sectionSummaries []map[string]interface{}
// 	for title, child := range ds.Tree.Children {
// 		summary := map[string]interface{}{
// 			"title":           title,
// 			"paragraphCount":  len(child.Content),
// 			"subsectionCount": len(child.Children),
// 		}
// 		sectionSummaries = append(sectionSummaries, summary)
// 	}

// 	// Identify main topics
// 	allText := ""
// 	for _, item := range ds.ContentItems {
// 		if item.Type == ContentTypeParagraph {
// 			allText += item.Content + " "
// 		}
// 	}
// 	mainTopics := extractKeyTerms(allText, 7)

// 	return map[string]interface{}{
// 		"title":                ds.Tree.Title,
// 		"sectionCount":         len(ds.Tree.Children),
// 		"totalParagraphs":      countContentByType(ds, ContentTypeParagraph),
// 		"totalWords":           totalWords,
// 		"averageSectionLength": totalWords / max(1, len(ds.Tree.Children)),
// 		"mainTopics":           mainTopics,
// 		"sections":             sectionSummaries,
// 	}
// }

// // Implementation of GroupParagraphsByTheme
// func groupParagraphsByTheme(ds *DocumentScraper, headingPath []string) ([]TextBlock, error) {
// 	section, err := ds.GetSectionWithContext(headingPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// TextBlock definition if not already defined
// 	type TextBlock struct {
// 		Paragraphs []ContentItem
// 		Theme      string
// 		KeyTerms   []string
// 		Summary    string
// 	}

// 	// If fewer than 2 paragraphs, just return a single block
// 	if len(section.Paragraphs) < 2 {
// 		if len(section.Paragraphs) == 0 {
// 			return []TextBlock{}, nil
// 		}

// 		return []TextBlock{
// 			{
// 				Paragraphs: section.Paragraphs,
// 				Theme:      section.Title,
// 				KeyTerms:   extractKeyTerms(section.Paragraphs[0].Content, 3),
// 				Summary:    truncateString(section.Paragraphs[0].Content, 100),
// 			},
// 		}, nil
// 	}

// 	// In a real implementation, this would use more sophisticated NLP techniques
// 	// For this example, we'll use a simple approach based on keyword overlap

// 	var blocks []TextBlock
// 	var currentBlock TextBlock
// 	var currentBlockText string

// 	for i, para := range section.Paragraphs {
// 		// For first paragraph, start a new block
// 		if i == 0 {
// 			currentBlock = TextBlock{
// 				Paragraphs: []ContentItem{para},
// 				KeyTerms:   extractKeyTerms(para.Content, 3),
// 			}
// 			currentBlockText = para.Content
// 			continue
// 		}

// 		// Decide whether to add to current block or start a new one
// 		currentKeyTerms := extractKeyTerms(currentBlockText, 5)
// 		paraKeyTerms := extractKeyTerms(para.Content, 5)

// 		// Check overlap
// 		overlap := keyTermOverlap(currentKeyTerms, paraKeyTerms)

// 		if overlap >= 0.3 || i == len(section.Paragraphs)-1 {
// 			// Add to current block
// 			currentBlock.Paragraphs = append(currentBlock.Paragraphs, para)
// 			currentBlockText += " " + para.Content

// 			// If it's the last paragraph, finalize the block
// 			if i == len(section.Paragraphs)-1 {
// 				currentBlock.KeyTerms = extractKeyTerms(currentBlockText, 4)
// 				currentBlock.Theme = inferTheme(currentBlock.KeyTerms, section.Title)
// 				currentBlock.Summary = generateBlockSummary(currentBlock.Paragraphs)
// 				blocks = append(blocks, currentBlock)
// 			}
// 		} else {
// 			// Finalize current block and start a new one
// 			currentBlock.KeyTerms = extractKeyTerms(currentBlockText, 4)
// 			currentBlock.Theme = inferTheme(currentBlock.KeyTerms, section.Title)
// 			currentBlock.Summary = generateBlockSummary(currentBlock.Paragraphs)
// 			blocks = append(blocks, currentBlock)

// 			// Start new block
// 			currentBlock = TextBlock{
// 				Paragraphs: []ContentItem{para},
// 			}
// 			currentBlockText = para.Content
// 		}
// 	}

// 	return blocks, nil
// }

// // Helper function to count content by type
// func countContentByType(ds *DocumentScraper, contentType ContentType) int {
// 	count := 0
// 	for _, item := range ds.ContentItems {
// 		if item.Type == contentType {
// 			count++
// 		}
// 	}
// 	return count
// }

// // keyTermOverlap calculates the overlap between two sets of key terms
// func keyTermOverlap(terms1, terms2 []string) float64 {
// 	// Convert to sets
// 	set1 := make(map[string]bool)
// 	for _, term := range terms1 {
// 		set1[term] = true
// 	}

// 	set2 := make(map[string]bool)
// 	for _, term := range terms2 {
// 		set2[term] = true
// 	}

// 	// Count overlaps
// 	overlap := 0
// 	for term := range set1 {
// 		if set2[term] {
// 			overlap++
// 		}
// 	}

// 	// Calculate Jaccard similarity
// 	union := len(set1) + len(set2) - overlap
// 	if union == 0 {
// 		return 0.0
// 	}

// 	return float64(overlap) / float64(union)
// }

// // inferTheme generates a theme based on key terms and heading
// func inferTheme(keyTerms []string, heading string) string {
// 	if len(keyTerms) == 0 {
// 		return heading
// 	}

// 	// In a real implementation, this would generate a more natural theme description
// 	// For this example, we'll combine the heading with key terms

// 	if len(keyTerms) <= 2 {
// 		return fmt.Sprintf("%s (%s)", heading, strings.Join(keyTerms, ", "))
// 	}

// 	return fmt.Sprintf("%s (%s, ...)", heading, strings.Join(keyTerms[:2], ", "))
// }

// // // generateBlockSummary creates a summary for a text block
// // func generateBlockSummary(paragraphs []ContentItem) string {
// // 	if len(paragraphs) == 0 {
// // 		return ""
// // 	}

// // 	// In a real implementation, this would use abstractive summarization
// // 	// For this example, just use the first sentence of the first paragraph

// // 	text := paragraphs[0].Content
// // 	sentences := splitIntoSentences(text)
// // 	if len(sentences) == 0 {
// // 		return ""
// // 	}

// // 	return truncateString(sentences[0], 100)
// // }
