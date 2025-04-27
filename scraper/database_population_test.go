package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
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
	// Replace these with your actual database credentials
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

	// 2. Create scraper for DruidAI
	websiteURL := "https://druidai.com"
	maxDepth := 2 // Keep depth limited for testing
	fmt.Printf("🔍 Creating scraper for %s (max depth: %d)\n", websiteURL, maxDepth)

	scraper, err := NewScraper(websiteURL, maxDepth)
	if err != nil {
		t.Fatalf("Failed to create scraper: %v", err)
	}

	// 3. Run the scraper
	fmt.Println("🚀 Running scraper and extracting site structure...")
	err = scraper.Run()
	if err != nil {
		t.Fatalf("Scraping failed: %v", err)
	}
	fmt.Printf("✅ Scraped %d pages\n", len(scraper.Data))

	// 4. Build site tree
	fmt.Println("🌳 Building site tree...")
	_, err = scraper.BuildSiteTree()
	if err != nil {
		t.Fatalf("Failed to build site tree: %v", err)
	}

	// 5. Create a company record for DruidAI
	fmt.Println("💾 Creating company record for DruidAI...")
	companyParams := db.CreateCompanyParams{
		Name: "DruidAI",
		Website: sql.NullString{
			String: websiteURL,
			Valid:  true,
		},
		Industry: sql.NullString{
			String: "Artificial Intelligence",
			Valid:  true,
		},
		Description: sql.NullString{
			String: "DruidAI is a company specializing in conversational AI and intelligent virtual assistants for enterprise solutions.",
			Valid:  true,
		},
		HeadquartersLocation: sql.NullString{
			String: "Bucharest, Romania",
			Valid:  true,
		},
		FoundedYear: sql.NullInt32{
			Int32: 2018,
			Valid: true,
		},
		IsPublic: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
	}

	company, err := queries.CreateCompany(ctx, companyParams)
	if err != nil {
		t.Fatalf("Failed to create company: %v", err)
	}
	fmt.Printf("✅ Created company with ID: %d\n", company.CompanyID)

	// 6. Create a fictional project
	fmt.Println("🏗️ Creating project for AI agent integration...")
	projectParams := db.CreateProjectParams{
		ProjectName: "AI Agent Integration Platform",
		Description: sql.NullString{
			String: "A project to evaluate and implement DruidAI's agent technology for customer service automation and enterprise knowledge base integration.",
			Valid:  true,
		},
	}

	project, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}
	fmt.Printf("✅ Created project with ID: %d\n", project.ProjectID)

	// 7. Associate company with project
	fmt.Println("🔗 Associating company with project...")
	associationParams := db.AssociateCompanyWithProjectParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
		AssociationNotes: sql.NullString{
			String: "Potential vendor for AI agent solution with strong capabilities in enterprise integration.",
			Valid:  true,
		},
		MatchingScore: sql.NullString{
			String: "0.85",
			Valid:  true,
		},
		ApproachStrategy: sql.NullString{
			String: "Request demo of their platform focusing on knowledge base integration and Microsoft ecosystem compatibility.",
			Valid:  true,
		},
	}

	_, err = queries.AssociateCompanyWithProject(ctx, associationParams)
	if err != nil {
		t.Fatalf("Failed to associate company with project: %v", err)
	}
	fmt.Println("✅ Associated company with project")

	// 8. Create company website record
	fmt.Println("🌐 Adding company website record...")
	websiteParams := db.CreateCompanyWebsiteParams{
		CompanyID: company.CompanyID,
		BaseUrl:   websiteURL,
		SiteTitle: sql.NullString{
			String: getWebsiteTitle(scraper, websiteURL),
			Valid:  true,
		},
		ScrapeFrequencyDays: sql.NullInt32{
			Int32: 30, // Check monthly
			Valid: true,
		},
		IsActive: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}

	website, err := queries.CreateCompanyWebsite(ctx, websiteParams)
	if err != nil {
		t.Fatalf("Failed to create company website: %v", err)
	}
	fmt.Printf("✅ Created company website with ID: %d\n", website.WebsiteID)

	// 9. Create datasource record
	fmt.Println("📊 Creating datasource record...")
	datasourceParams := db.CreateDatasourceParams{
		SourceType: "website",
	}

	datasource, err := queries.CreateDatasource(ctx, datasourceParams)
	if err != nil {
		t.Fatalf("Failed to create datasource: %v", err)
	}
	fmt.Printf("✅ Created datasource with ID: %d\n", datasource.DatasourceID)

	// Update website with datasource ID
	updateWebsiteParams := db.UpdateCompanyWebsiteParams{
		WebsiteID: website.WebsiteID,
		BaseUrl:   website.BaseUrl,
		SiteTitle: website.SiteTitle,
		LastScrapedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ScrapeFrequencyDays: website.ScrapeFrequencyDays,
		IsActive:            website.IsActive,
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
	}

	_, err = queries.UpdateCompanyWebsite(ctx, updateWebsiteParams)
	if err != nil {
		t.Fatalf("Failed to update company website with datasource ID: %v", err)
	}

	// 10. Create website pages
	fmt.Println("📄 Adding website pages to database...")
	addedPages := 0
	for pageURL, pageData := range scraper.Data {
		// Skip external URLs
		parsedURL, err := url.Parse(pageURL)
		if err != nil || !isSameHostname(parsedURL.Hostname(), websiteURL) {
			continue
		}

		// Get page path
		pagePath := parsedURL.Path
		if pagePath == "" {
			pagePath = "/"
		}

		// Determine page depth
		depth := countPathSegments(pagePath)

		// Create page record
		pageParams := db.CreateWebsitePageParams{
			WebsiteID: website.WebsiteID,
			Url:       pageURL,
			Path:      pagePath,
			Title: sql.NullString{
				String: pageData.Title,
				Valid:  pageData.Title != "",
			},
			Depth: int32(depth),
			ExtractStatus: sql.NullString{
				String: "completed",
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: datasource.DatasourceID,
				Valid: true,
			},
		}

		_, err = queries.CreateWebsitePage(ctx, pageParams)
		if err != nil {
			// Log error but continue with other pages
			log.Printf("Error creating page %s: %v", pageURL, err)
			continue
		}
		addedPages++
	}
	fmt.Printf("✅ Added %d website pages to database\n", addedPages)

	// 11. Create some paragraphs from content
	fmt.Println("📝 Extracting paragraphs from content...")
	addedParagraphs := 0

	// Select a few key pages to extract paragraphs from
	for pageURL, pageData := range scraper.Data {
		if strings.Contains(pageURL, "about") || strings.Contains(pageURL, "platform") ||
			pageURL == websiteURL || strings.Contains(pageURL, "solutions") {

			// Basic content splitting - in a real implementation, you'd want more sophisticated extraction
			contentParts := splitIntoParagraphs(pageData.Content)

			for _, paragraph := range contentParts {
				// Skip very short paragraphs
				if len(paragraph) < 50 {
					continue
				}

				paragraphParams := db.CreateParagraphParams{
					DatasourceID: sql.NullInt32{
						Int32: datasource.DatasourceID,
						Valid: true,
					},
					Content: paragraph,
					MainIdea: sql.NullString{
						String: extractMainIdea(paragraph),
						Valid:  true,
					},
					Classification: sql.NullString{
						String: classifyParagraph(paragraph),
						Valid:  true,
					},
					ConfidenceScore: sql.NullString{
						String: "0.75", // Fictional score
						Valid:  true,
					},
				}

				_, err = queries.CreateParagraph(ctx, paragraphParams)
				if err != nil {
					// Log error but continue with other paragraphs
					log.Printf("Error creating paragraph: %v", err)
					continue
				}
				addedParagraphs++

				// Limit paragraphs for testing purposes
				if addedParagraphs >= 10 {
					break
				}
			}
		}

		// Limit to just a few pages for testing purposes
		if addedParagraphs >= 10 {
			break
		}
	}
	fmt.Printf("✅ Added %d paragraphs to database\n", addedParagraphs)

	fmt.Println("✅ Successfully populated database with DruidAI data!")
}

// Helper functions

func getWebsiteTitle(s *Scraper, url string) string {
	if data, exists := s.Data[url]; exists && data.Title != "" {
		return data.Title
	}
	return "DruidAI Website" // Default title
}

func isSameHostname(hostname, baseURL string) bool {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	baseHost := strings.TrimPrefix(parsed.Hostname(), "www.")
	hostname = strings.TrimPrefix(hostname, "www.")

	return hostname == baseHost
}

func countPathSegments(path string) int {
	segments := strings.Split(path, "/")
	count := 0
	for _, segment := range segments {
		if segment != "" {
			count++
		}
	}
	return count
}

func splitIntoParagraphs(content string) []string {
	// Basic implementation - in a real scenario, you'd want more sophisticated text processing
	rawParagraphs := strings.Split(content, "\n\n")
	var paragraphs []string

	for _, p := range rawParagraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			paragraphs = append(paragraphs, p)
		}
	}

	return paragraphs
}

func extractMainIdea(paragraph string) string {
	// Simplified implementation - in a real scenario, you'd use NLP techniques
	words := strings.Fields(paragraph)
	if len(words) <= 10 {
		return paragraph
	}

	// Just take the first 10-15 words as a simple summarization
	end := 15
	if len(words) < end {
		end = len(words)
	}

	return strings.Join(words[:end], " ") + "..."
}

func classifyParagraph(paragraph string) string {
	// Simplified classification based on keyword presence
	paragraph = strings.ToLower(paragraph)

	if strings.Contains(paragraph, "ai") || strings.Contains(paragraph, "artificial intelligence") {
		return "AI Technology"
	} else if strings.Contains(paragraph, "customer") || strings.Contains(paragraph, "service") {
		return "Customer Service"
	} else if strings.Contains(paragraph, "platform") || strings.Contains(paragraph, "solution") {
		return "Product"
	} else if strings.Contains(paragraph, "integration") || strings.Contains(paragraph, "api") {
		return "Technical Integration"
	}

	return "General Information"
}
