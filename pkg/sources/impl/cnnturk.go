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
		return []NewsItem{}, nil
	}

	category := c.categories[categoryIndex-1]
	feedURL, ok := c.feedURLs[category]
	if !ok {
		return []NewsItem{}, nil
	}


	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Fetch the feed
	resp, err := client.Get(feedURL)
	if err != nil {
		return []NewsItem{}, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []NewsItem{}, nil
	}

	
	// Try to parse as RSS
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return []NewsItem{}, nil
	}

	if len(rss.Channel.Items) == 0 {
		return []NewsItem{}, nil
	}

	
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
		
	}
	
	return newsItems, nil
}

