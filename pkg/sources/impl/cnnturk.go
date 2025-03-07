package impl

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// CNNTurkSource implements a news source for CNN Türk RSS feeds
type CNNTurkSource struct {
	name       string
	categories []string
	feedURLs   map[string]string
}

// NewCNNTurkSource creates a new CNNTurkSource instance
func NewCNNTurkSource() *CNNTurkSource {
	categories := []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "SAĞLIK", "TEKNOLOJİ"}
	
	feedURLs := map[string]string{
		"GÜNDEM":    "https://www.cnnturk.com/feed/rss/turkiye/news",
		"DÜNYA":     "https://www.cnnturk.com/feed/rss/dunya/news",
		"EKONOMİ":   "https://www.cnnturk.com/feed/rss/ekonomi/news",
		"SPOR":      "https://www.cnnturk.com/feed/rss/spor/news",
		"SAĞLIK":    "https://www.cnnturk.com/feed/rss/saglik/news",
		"TEKNOLOJİ": "https://www.cnnturk.com/feed/rss/bilim-teknoloji/news",
	}
	
	return &CNNTurkSource{
		name:       "CNN Türk",
		categories: categories,
		feedURLs:   feedURLs,
	}
}

// Name returns the name of the news source
func (c *CNNTurkSource) Name() string {
	return c.name
}

// Categories returns the available categories for this news source
func (c *CNNTurkSource) Categories() []string {
	return c.categories
}

// FetchNews fetches news items for the specified category
func (c *CNNTurkSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	// Adjust for 0-based indexing
	if categoryIndex >= 0 && categoryIndex < len(c.categories) {
		// This is 0-based indexing, convert to 1-based for our internal use
		categoryIndex++
	}

	if categoryIndex < 1 || categoryIndex > len(c.categories) {
		fmt.Printf("Debug: Invalid category index: %d, returning empty results\n", categoryIndex)
		return []NewsItem{}, nil
	}

	category := c.categories[categoryIndex-1]
	feedURL, ok := c.feedURLs[category]
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
	maxItems := 30
	if len(rss.Channel.Items) < maxItems {
		maxItems = len(rss.Channel.Items)
	}
	
	newsItems := make([]NewsItem, 0, maxItems)
	for i := 0; i < maxItems; i++ {
		item := rss.Channel.Items[i]
		
		// Skip items with empty titles or links
		if item.Title == "" {
			continue
		}
		
		// Clean up the title
		title := strings.TrimSpace(item.Title)
		title = strings.TrimPrefix(title, "<![CDATA[")
		title = strings.TrimSuffix(title, "]]>")
		
		// Get the URL from either the link or the GUID
		url := item.Link
		if url == "" && item.GUID.IsPermaLink == "true" {
			url = item.GUID.Value
		}
		
		// If URL is still empty or doesn't start with http, construct it
		if url == "" || !strings.HasPrefix(url, "http") {
			// Use the GUID value as the ID
			id := item.GUID.Value
			
			// Determine the category path
			categoryPath := ""
			switch category {
			case "GÜNDEM":
				categoryPath = "turkiye"
			case "DÜNYA":
				categoryPath = "dunya"
			case "EKONOMİ":
				categoryPath = "ekonomi"
			case "SPOR":
				categoryPath = "spor"
			case "SAĞLIK":
				categoryPath = "saglik"
			case "TEKNOLOJİ":
				categoryPath = "bilim-teknoloji"
			}
			
			url = fmt.Sprintf("https://www.cnnturk.com/%s/%s", categoryPath, id)
		}
		
		// Add the item to the list
		newsItems = append(newsItems, NewsItem{
			Title: title,
			URL:   url,
		})
		
		if i < 3 {
			fmt.Printf("Debug: Item %d: Title=%s, URL=%s\n", i, title, url)
		}
	}
	
	fmt.Printf("Debug: Found %d news items\n", len(newsItems))
	return newsItems, nil
}

