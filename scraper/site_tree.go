package scraper

import (
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
)

// TreeNode represents a node in the site tree
type TreeNode struct {
	URL      string
	Path     string
	Title    string
	Children map[string]*TreeNode
}

// BuildSiteTree constructs a tree representation of the site from scraped data
func (s *Scraper) BuildSiteTree() (*TreeNode, error) {
	// Create the root node
	baseURL, err := url.Parse(s.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	// Clean the base hostname (remove www if present)
	baseHostname := strings.TrimPrefix(baseURL.Hostname(), "www.")

	// Create root node
	root := &TreeNode{
		URL:      s.BaseURL,
		Path:     "/",
		Title:    "Root",
		Children: make(map[string]*TreeNode),
	}

	// Find the title for the root node if available
	if data, exists := s.Data[s.BaseURL]; exists {
		root.Title = data.Title
	}

	// Add all pages to the tree
	for pageURL, data := range s.Data {
		// Skip if it's the root
		if pageURL == s.BaseURL {
			continue
		}

		// Parse URL
		parsedURL, err := url.Parse(pageURL)
		if err != nil {
			fmt.Printf("Error parsing URL %s: %v\n", pageURL, err)
			continue
		}

		// Skip URLs with different hostnames
		pageHostname := strings.TrimPrefix(parsedURL.Hostname(), "www.")
		if pageHostname != baseHostname {
			continue
		}

		// Get path segments
		pathSegments := getPathSegments(parsedURL.Path)

		// Insert node into the tree
		currentNode := root
		currentPath := "/"

		for i, segment := range pathSegments {
			if segment == "" {
				continue
			}

			currentPath = path.Join(currentPath, segment)

			// Check if this path segment already exists as a child
			if child, exists := currentNode.Children[segment]; exists {
				currentNode = child
			} else {
				// Create a new node for this path segment
				newNode := &TreeNode{
					URL:      pageURL,
					Path:     currentPath,
					Title:    segment, // Default title is the path segment
					Children: make(map[string]*TreeNode),
				}

				// If this is the last segment, use the page title
				if i == len(pathSegments)-1 {
					newNode.Title = data.Title
					if newNode.Title == "" {
						newNode.Title = segment
					}
				}

				currentNode.Children[segment] = newNode
				currentNode = newNode
			}
		}
	}

	return root, nil
}

// getPathSegments splits a URL path into segments
func getPathSegments(urlPath string) []string {
	// Clean the path first
	urlPath = path.Clean(urlPath)
	// Split by /
	return strings.Split(urlPath, "/")
}

// PrintSiteTree prints the site tree with proper indentation
func PrintSiteTree(node *TreeNode, indent string) {
	if node == nil {
		return
	}

	// Print current node
	fmt.Printf("%s%s (%s)\n", indent, node.Title, node.Path)

	// Get sorted child keys for consistent output
	var childKeys []string
	for key := range node.Children {
		childKeys = append(childKeys, key)
	}
	sort.Strings(childKeys)

	// Print children
	for _, key := range childKeys {
		PrintSiteTree(node.Children[key], indent+"  ")
	}
}

// ExportSiteTreeDOT exports the site tree in DOT format for visualization
func ExportSiteTreeDOT(node *TreeNode) string {
	var sb strings.Builder

	sb.WriteString("digraph SiteMap {\n")
	sb.WriteString("  node [shape=box, style=filled, fillcolor=lightblue];\n")
	sb.WriteString("  rankdir=LR;\n\n")

	// Use a map to track node IDs
	nodeIDs := make(map[string]int)
	counter := 0

	// Generate unique IDs for each node
	nodeIDs[node.Path] = counter
	counter++

	// Add nodes and edges
	exportDOTNodes(&sb, node, nodeIDs, &counter)
	exportDOTEdges(&sb, node, nodeIDs)

	sb.WriteString("}\n")
	return sb.String()
}

// exportDOTNodes adds nodes to DOT representation
func exportDOTNodes(sb *strings.Builder, node *TreeNode, nodeIDs map[string]int, counter *int) {
	if node == nil {
		return
	}

	// Add current node
	nodeTitle := strings.Replace(node.Title, "\"", "\\\"", -1)
	sb.WriteString(fmt.Sprintf("  node%d [label=\"%s\"];\n", nodeIDs[node.Path], nodeTitle))

	// Add children
	for _, child := range node.Children {
		if _, exists := nodeIDs[child.Path]; !exists {
			nodeIDs[child.Path] = *counter
			*counter++
		}
		exportDOTNodes(sb, child, nodeIDs, counter)
	}
}

// exportDOTEdges adds edges to DOT representation
func exportDOTEdges(sb *strings.Builder, node *TreeNode, nodeIDs map[string]int) {
	if node == nil {
		return
	}

	// Add edges to children
	for _, child := range node.Children {
		sb.WriteString(fmt.Sprintf("  node%d -> node%d;\n", nodeIDs[node.Path], nodeIDs[child.Path]))
		exportDOTEdges(sb, child, nodeIDs)
	}
}
