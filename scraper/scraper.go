package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

// Scraper represents a web scraper with configurable depth
type Scraper struct {
	MaxDepth     int
	BaseURL      string
	LinksToVisit map[string]int // maps URL to its depth
	VisitedLinks map[string]bool
	mu           sync.Mutex
	Data         map[string]PageData
}

// PageData stores information scraped from a page
type PageData struct {
	URL     string
	Title   string
	Content string
	Links   []string
}

// NewScraper creates a new scraper instance
func NewScraper(baseURL string, maxDepth int) (*Scraper, error) {
	// Validate URL
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	return &Scraper{
		MaxDepth:     maxDepth,
		BaseURL:      baseURL,
		LinksToVisit: make(map[string]int),
		VisitedLinks: make(map[string]bool),
		Data:         make(map[string]PageData),
	}, nil
}

// isSameDomain checks if the link belongs to the same domain as the base URL
func (s *Scraper) isSameDomain(link string) bool {
	baseURL, _ := url.Parse(s.BaseURL)
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	baseHost := baseURL.Hostname()
	linkHost := linkURL.Hostname()

	// Handle www subdomain variations
	baseHost = strings.TrimPrefix(baseHost, "www.")
	linkHost = strings.TrimPrefix(linkHost, "www.")

	return linkHost == baseHost
}

// GatherLinks collects all links up to the specified depth
func (s *Scraper) GatherLinks() error {
	baseDomain := getDomain(s.BaseURL)
	wwwDomain := "www." + baseDomain

	// Create a collector with options for handling redirects
	c := colly.NewCollector(
		colly.MaxDepth(s.MaxDepth),
		colly.AllowedDomains(baseDomain, wwwDomain),
		colly.AllowURLRevisit(),
	)

	// Add the base URL as the starting point
	s.LinksToVisit[s.BaseURL] = 0

	// Find and store all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if absoluteURL == "" {
			return
		}

		currentDepth := getDepthFromContext(e.Request.Ctx)
		if currentDepth >= s.MaxDepth {
			return
		}

		// Debug - print found links
		// fmt.Printf("Found link on %s: %s -> %s\n", e.Request.URL, link, absoluteURL)

		// Only process links that belong to the same domain
		if s.isSameDomain(absoluteURL) {
			s.mu.Lock()
			if _, exists := s.LinksToVisit[absoluteURL]; !exists && !s.VisitedLinks[absoluteURL] {
				s.LinksToVisit[absoluteURL] = currentDepth + 1
				fmt.Printf("Adding to visit queue: %s (depth: %d)\n", absoluteURL, currentDepth+1)
			}
			s.mu.Unlock()
		}
	})

	c.OnRequest(func(r *colly.Request) {
		s.mu.Lock()
		s.VisitedLinks[r.URL.String()] = true
		s.mu.Unlock()
		fmt.Printf("Gathering links: %s\n", r.URL.String())
	})

	// Debug - print response info
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Got response from %s: %d bytes, status: %d\n",
			r.Request.URL, len(r.Body), r.StatusCode)
	})

	// Debug - print errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on %s: %s\n", r.Request.URL, err)
	})

	// Start with the base URL
	err := c.Visit(s.BaseURL)
	if err != nil {
		return err
	}

	c.Wait()
	fmt.Printf("Found %d links to visit\n", len(s.LinksToVisit))
	return nil
}

// ScrapeLinks performs the actual scraping on the gathered links
func (s *Scraper) ScrapeLinks() error {
	baseDomain := getDomain(s.BaseURL)
	wwwDomain := "www." + baseDomain

	c := colly.NewCollector(
		colly.AllowedDomains(baseDomain, wwwDomain),
		colly.AllowURLRevisit(),
	)

	// Debug - print response info
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Scraping response from %s: %d bytes\n", r.Request.URL, len(r.Body))
	})

	// Extract page title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		s.mu.Lock()
		pageData := s.Data[url]
		pageData.Title = e.Text
		s.Data[url] = pageData
		s.mu.Unlock()
	})

	// Extract page content (example: body text)
	c.OnHTML("body", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		s.mu.Lock()
		pageData := s.Data[url]
		pageData.Content = strings.TrimSpace(e.Text)
		s.Data[url] = pageData
		s.mu.Unlock()
	})

	// Extract links - more comprehensive approach
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		// Debug - print all found links
		// fmt.Printf("Found link on %s: %s -> %s\n", url, link, absoluteURL)

		if absoluteURL != "" {
			// Store all links from the same domain
			if s.isSameDomain(absoluteURL) {
				s.mu.Lock()
				pageData := s.Data[url]
				// Avoid duplicate links
				isDuplicate := false
				for _, existingLink := range pageData.Links {
					if existingLink == absoluteURL {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					pageData.Links = append(pageData.Links, absoluteURL)
					s.Data[url] = pageData
				}
				s.mu.Unlock()
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Scraping: %s\n", r.URL.String())
		s.mu.Lock()
		url := r.URL.String()
		if _, exists := s.Data[url]; !exists {
			s.Data[url] = PageData{
				URL:   url,
				Links: []string{},
			}
		}
		s.mu.Unlock()
	})

	// Debug - print errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on %s: %s\n", r.Request.URL, err)
	})

	// Visit each link in LinksToVisit
	for link := range s.LinksToVisit {
		err := c.Visit(link)
		if err != nil {
			fmt.Printf("Error visiting %s: %v\n", link, err)
			// Continue with other links
		}
	}

	c.Wait()
	return nil
}

// Run executes the complete scraping process
func (s *Scraper) Run() error {
	// First gather all links
	err := s.GatherLinks()
	if err != nil {
		return fmt.Errorf("error gathering links: %v", err)
	}

	// Then scrape them
	err = s.ScrapeLinks()
	if err != nil {
		return fmt.Errorf("error scraping links: %v", err)
	}

	return nil
}

// Helper function to get depth from context
func getDepthFromContext(ctx *colly.Context) int {
	if ctx == nil {
		return 0
	}
	depth, ok := ctx.GetAny("depth").(int)
	if !ok {
		return 0
	}
	return depth
}

// Helper function to extract domain from URL
func getDomain(baseURL string) string {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}
