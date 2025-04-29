package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// TestPopulateDruidAI demonstrates scraping DruidAI website and adding the data to the database
func TestPopulateDruidAI(t *testing.T) {
	// Skip this test in automated test runs, run it manually
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. Set up database connection
	connStr := "postgresql://root:secret@localhost:5432/musli?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Test the connection
	err = dbConn.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✅ Successfully connected to database")

	// Create a queries object
	queries := db.New(dbConn)
	ctx := context.Background()

	// 2. Create enhanced scraper for the website
	websiteURL := "https://druidai.com"
	maxDepth := 2 // Keep depth limited for testing
	fmt.Printf("🔍 Creating enhanced scraper for %s (max depth: %d)\n", websiteURL, maxDepth)

	enhancedScraper, err := NewEnhancedScraper(websiteURL, maxDepth)
	if err != nil {
		t.Fatalf("Failed to create enhanced scraper: %v", err)
	}

	// 3. Run the enhanced scraper
	fmt.Println("🚀 Running enhanced scraper to extract content...")
	err = enhancedScraper.Run()
	if err != nil {
		t.Fatalf("Enhanced scraping failed: %v", err)
	}
	fmt.Printf("✅ Extracted %d unique content items\n", len(enhancedScraper.ContentItems))

	// 4. Build site tree (optional, for visualization)
	fmt.Println("🌳 Building site tree...")
	siteTree, err := enhancedScraper.BuildSiteTree()
	if err != nil {
		t.Fatalf("Failed to build site tree: %v", err)
	}
	fmt.Println("✅ Built site tree")

	// Optional: Print the site tree for debugging
	fmt.Println("\nSite Tree Structure:")
	PrintSiteTree(siteTree, "")

	// 5. Create a user
	fmt.Println("👤 Creating user record...")
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username: "analyst_" + randomString(6),
		Password: "secure_" + randomString(10),
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("✅ Created user with ID: %d\n", user.UserID)

	// 6. Create a company record
	fmt.Println("💾 Creating company record...")
	companyName := "Unknown"

	// Try to extract company name from title
	if len(enhancedScraper.Data) > 0 {
		for _, pageData := range enhancedScraper.Data {
			if strings.Contains(pageData.Title, "DRUID") || strings.Contains(pageData.Title, "Druid") {
				parts := strings.Split(pageData.Title, " - ")
				if len(parts) > 1 {
					companyName = parts[len(parts)-1]
				} else {
					companyName = "DRUID AI"
				}
				break
			}
		}
	}

	companyParams := db.CreateCompanyParams{
		UserID:      user.UserID,
		CompanyName: companyName,
		Industry:    sql.NullString{String: "Artificial Intelligence", Valid: true},
		Website:     sql.NullString{String: websiteURL, Valid: true},
		Description: sql.NullString{String: "Company extracted from website content", Valid: true},
	}

	company, err := queries.CreateCompany(ctx, companyParams)
	if err != nil {
		t.Fatalf("Failed to create company: %v", err)
	}
	fmt.Printf("✅ Created company with ID: %d\n", company.CompanyID)

	// 7. Create datasource record for the website
	fmt.Println("📊 Creating datasource record...")
	datasourceParams := db.CreateDatasourceParams{
		SourceType: db.DatasourceTypeWebsite,
		Link:       sql.NullString{String: websiteURL, Valid: true},
		FileName:   sql.NullString{String: "website_data.html", Valid: true},
		FileData:   []byte("Website content extracted via scraper"),
	}

	datasource, err := queries.CreateDatasource(ctx, datasourceParams)
	if err != nil {
		t.Fatalf("Failed to create datasource: %v", err)
	}
	fmt.Printf("✅ Created datasource with ID: %d\n", datasource.DatasourceID)

	// 8. Associate datasource with company
	fmt.Println("🔗 Associating datasource with company...")
	err = queries.AssociateDatasourceWithCompany(ctx, db.AssociateDatasourceWithCompanyParams{
		CompanyID:    company.CompanyID,
		DatasourceID: datasource.DatasourceID,
	})
	if err != nil {
		t.Fatalf("Failed to associate datasource with company: %v", err)
	}
	fmt.Println("✅ Associated datasource with company")

	// 9. Create paragraphs from extracted content
	fmt.Println("📝 Creating paragraphs from extracted content...")
	addedParagraphs := 0

	// Keep track of content we've already added to the database
	// This provides additional deduplication at the database level
	seenHashes := make(map[string]bool)

	for _, item := range enhancedScraper.ContentItems {
		// Skip if we've already added this exact content
		if seenHashes[item.Hash] {
			continue
		}
		seenHashes[item.Hash] = true

		// Skip if content is too short
		if countWords(item.Paragraph) < 10 {
			continue
		}

		paragraphParams := db.CreateParagraphParams{
			DatasourceID: datasource.DatasourceID,
			Title:        sql.NullString{String: cleanText(item.Title), Valid: true},
			MainIdea:     sql.NullString{String: "", Valid: false}, // Leave main idea empty as requested
			Content:      cleanText(item.Paragraph),
		}

		_, err = queries.CreateParagraph(ctx, paragraphParams)
		if err != nil {
			log.Printf("Error creating paragraph: %v", err)
			continue
		}
		addedParagraphs++

		// Print sample of what we're adding
		if addedParagraphs <= 5 {
			fmt.Printf("  Added: Title: %s\n  Paragraph: %s\n\n",
				truncateString(item.Title, 40),
				truncateString(item.Paragraph, 100))
		}
	}
	fmt.Printf("✅ Added %d unique paragraphs to database\n", addedParagraphs)

	// 10. Create a contact for the company (optional)
	fmt.Println("👥 Creating a contact for the company...")
	contactParams := db.CreateContactParams{
		CompanyID: company.CompanyID,
		FirstName: "Contact",
		LastName:  "Representative",
		Position:  sql.NullString{String: "Unknown Position", Valid: true},
		Email:     sql.NullString{String: "contact@" + getDomain(websiteURL), Valid: true},
		Phone:     sql.NullString{String: "Unknown", Valid: true},
		Notes:     sql.NullString{String: "Contact extracted from website", Valid: true},
	}

	contact, err := queries.CreateContact(ctx, contactParams)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}
	fmt.Printf("✅ Created contact with ID: %d\n", contact.ContactID)

	// 11. Associate datasource with contact
	fmt.Println("🔗 Associating datasource with contact...")
	err = queries.AssociateDatasourceWithContact(ctx, db.AssociateDatasourceWithContactParams{
		ContactID:    contact.ContactID,
		DatasourceID: datasource.DatasourceID,
	})
	if err != nil {
		t.Fatalf("Failed to associate datasource with contact: %v", err)
	}
	fmt.Println("✅ Associated datasource with contact")

	fmt.Println("✅ Successfully populated database with website content!")
}

// Helper function to generate random strings
func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
		time.Sleep(1 * time.Nanosecond) // To ensure uniqueness
	}
	return string(b)
}
