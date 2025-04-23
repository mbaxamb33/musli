package scraper

import (
	"fmt"
	"os"
	"testing"
)

func TestBuildSiteTree(t *testing.T) {
	url := "https://druidai.com"
	depth := 2

	fmt.Printf("Creating scraper for %s (depth %d)\n", url, depth)
	s, err := NewScraper(url, depth)
	if err != nil {
		t.Fatalf("Failed to create scraper: %v", err)
	}

	fmt.Println("Running scraper...")
	err = s.Run()
	if err != nil {
		t.Fatalf("Scraping failed: %v", err)
	}

	fmt.Printf("Scraped %d pages\n", len(s.Data))

	// Build the site tree
	fmt.Println("Building site tree...")
	root, err := s.BuildSiteTree()
	if err != nil {
		t.Fatalf("Failed to build site tree: %v", err)
	}

	// Print the tree
	fmt.Println("\nSite Tree:")
	PrintSiteTree(root, "")

	// Export as DOT graph
	fmt.Println("\nExporting DOT graph...")
	dotGraph := ExportSiteTreeDOT(root)

	// Save DOT file
	dotFile := "sitemap.dot"
	err = os.WriteFile(dotFile, []byte(dotGraph), 0644)
	if err != nil {
		t.Fatalf("Failed to write DOT file: %v", err)
	}

	fmt.Printf("DOT graph saved to %s\n", dotFile)
	fmt.Println("To visualize the graph, install Graphviz and run: dot -Tpng sitemap.dot -o sitemap.png")
}
