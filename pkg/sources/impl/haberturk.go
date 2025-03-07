package impl

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HaberturkSource implements a news source for Habertürk RSS feeds
type HaberturkSource struct {
	name       string
	categories []string
	feedURLs   map[string]string
}

// NewHaberturkSource creates a new HaberturkSource instance
func NewHaberturkSource() *HaberturkSource {
	categories := []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "SAĞLIK", "TEKNOLOJİ"}
	
	feedURLs := map[string]string{
		"GÜNDEM":    "https://www.haberturk.com/rss/kategori/gundem.xml",
		"DÜNYA":     "https://www.haberturk.com/rss/kategori/dunya.xml",
		"EKONOMİ":   "https://www.haberturk.com/rss/kategori/ekonomi.xml",
		"SPOR":      "https://www.haberturk.com/rss/kategori/spor.xml",
		"SAĞLIK":    "https://www.haberturk.com/rss/kategori/saglik.xml",
		"TEKNOLOJİ": "https://www.haberturk.com/rss/kategori/teknoloji.xml",
	}
	
	return &HaberturkSource{
		name:       "Habertürk",
		categories: categories,
		feedURLs:   feedURLs,
	}
}

// Name returns the name of the news source
func (h *HaberturkSource) Name() string {
	return h.name
}

// Categories returns the available categories for this news source
func (h *HaberturkSource) Categories() []string {
	return h.categories
}

// FetchNews fetches news items for the specified category
func (h *HaberturkSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	// Adjust for 0-based indexing
	if categoryIndex >= 0 && categoryIndex < len(h.categories) {
		// This is 0-based indexing, convert to 1-based for our internal use
		categoryIndex++
	}
	
	if categoryIndex < 1 || categoryIndex > len(h.categories) {
		fmt.Printf("Debug: Invalid category index: %d, returning empty results\n", categoryIndex)
		return []NewsItem{}, nil
	}

	category := h.categories[categoryIndex-1]
	feedURL, ok := h.feedURLs[category]
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

	// Try to parse as RSS
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Printf("Debug: Failed to parse as RSS: %v\n", err)
		return []NewsItem{}, nil
	}

	if len(rss.Channel.Items) == 0 {
		fmt.Printf("Debug: No items found in RSS feed\n")
		return []NewsItem{}, nil
	}

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
