package impl

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// NTVSource implements a news source for NTV RSS feeds
type NTVSource struct {
	name       string
	categories []string
	feedURLs   map[string]string
}

// NewNTVSource creates a new NTVSource instance
func NewNTVSource() *NTVSource {
	categories := []string{"SON DAKİKA", "GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "SAĞLIK", "TEKNOLOJİ"}
	
	feedURLs := map[string]string{
		"SON DAKİKA": "https://www.ntv.com.tr/son-dakika.rss",
		"GÜNDEM":     "https://www.ntv.com.tr/turkiye.rss",
		"DÜNYA":      "https://www.ntv.com.tr/dunya.rss",
		"EKONOMİ":    "https://www.ntv.com.tr/ekonomi.rss",
		"SPOR":       "https://www.ntv.com.tr/spor.rss",
		"SAĞLIK":     "https://www.ntv.com.tr/saglik.rss",
		"TEKNOLOJİ":  "https://www.ntv.com.tr/teknoloji.rss",
	}
	
	return &NTVSource{
		name:       "NTV",
		categories: categories,
		feedURLs:   feedURLs,
	}
}

// Name returns the name of the news source
func (n *NTVSource) Name() string {
	return n.name
}

// Categories returns the available categories for this news source
func (n *NTVSource) Categories() []string {
	return n.categories
}

// FetchNews fetches news items for the specified category
func (n *NTVSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	// Adjust for 0-based indexing
	if categoryIndex >= 0 && categoryIndex < len(n.categories) {
		// This is 0-based indexing, convert to 1-based for our internal use
		categoryIndex++
	}
	
	if categoryIndex < 1 || categoryIndex > len(n.categories) {
		fmt.Printf("Debug: Invalid category index: %d, returning empty results\n", categoryIndex)
		return []NewsItem{}, nil
	}

	category := n.categories[categoryIndex-1]
	feedURL, ok := n.feedURLs[category]
	if !ok {
		fmt.Printf("Debug: No feed URL for category: %s, returning empty results\n", category)
		return []NewsItem{}, nil
	}

	fmt.Printf("Debug: Fetching RSS feed from URL: %s\n", feedURL)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Fetch the feed
	resp, err := client.Get(feedURL)
	if err != nil {
		fmt.Printf("Debug: Error fetching feed: %v, returning empty results\n", err)
		return []NewsItem{}, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Debug: Error reading response body: %v, returning empty results\n", err)
		return []NewsItem{}, nil
	}

	fmt.Printf("Debug: Received %d bytes of data\n", len(body))

	// Try to parse as Atom first (NTV uses Atom format)
	var atom Atom
	err1 := xml.Unmarshal(body, &atom)
	if err1 == nil && len(atom.Entries) > 0 {
		fmt.Printf("Debug: Successfully parsed as Atom, found %d entries\n", len(atom.Entries))
		
		// Process Atom entries
		newsItems := make([]NewsItem, 0, min(len(atom.Entries), 30))
		for i, entry := range atom.Entries {
			if i >= 30 {
				break // Limit to 30 items
			}
			
			// Skip entries with empty titles or links
			if entry.Title == "" || entry.Link.Href == "" {
				continue
			}
			
			// Clean up the title
			title := strings.TrimSpace(entry.Title)
			
			// Add the entry to the list
			newsItems = append(newsItems, NewsItem{
				Title: title,
				URL:   entry.Link.Href,
			})
			
			if i < 3 {
				fmt.Printf("Debug: Item %d: Title=%s, URL=%s\n", i, title, entry.Link.Href)
			}
		}
		
		fmt.Printf("Debug: Found %d news items\n", len(newsItems))
		return newsItems, nil
	}

	// If Atom parsing failed, try RSS
	var rss RSS
	err2 := xml.Unmarshal(body, &rss)
	if err2 == nil && len(rss.Channel.Items) > 0 {
		fmt.Printf("Debug: Successfully parsed as RSS, found %d items\n", len(rss.Channel.Items))
		
		// Process RSS items
		newsItems := make([]NewsItem, 0, min(len(rss.Channel.Items), 30))
		for i, item := range rss.Channel.Items {
			if i >= 30 {
				break // Limit to 30 items
			}
			
			// Skip items with empty titles or links
			if item.Title == "" || item.Link == "" {
				continue
			}
			
			// Clean up the title
			title := strings.TrimSpace(item.Title)
			
			// Add the item to the list
			newsItems = append(newsItems, NewsItem{
				Title: title,
				URL:   item.Link,
			})
			
			if i < 3 {
				fmt.Printf("Debug: Item %d: Title=%s, URL=%s\n", i, title, item.Link)
			}
		}
		
		fmt.Printf("Debug: Found %d news items\n", len(newsItems))
		return newsItems, nil
	}

	// If we get here, both parsing attempts failed
	fmt.Printf("Debug: Failed to parse as Atom: %v\n", err1)
	fmt.Printf("Debug: Failed to parse as RSS: %v\n", err2)
	fmt.Printf("Debug: Kategori: %s, Kategori Numarası: %d\n", category, categoryIndex)
	fmt.Printf("Debug: Haber bulunamadı.\n")
	return []NewsItem{}, nil
}


