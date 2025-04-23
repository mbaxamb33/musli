package scraper

import (
	"fmt"
	"testing"
)

func TestRunScraper(t *testing.T) {
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
	for u, data := range s.Data {
		fmt.Printf("- %s: %s (Found %d links)\n", u, data.Title, len(data.Links))
	}
}

func TestPrint(t *testing.T) {
	fmt.Println("🔍 This should print if you use -v")
}
