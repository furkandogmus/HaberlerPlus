package impl

import (
	"encoding/xml"
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
		return []NewsItem{}, nil
	}

	category := n.categories[categoryIndex-1]
	feedURL, ok := n.feedURLs[category]
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

	
	// Try to parse as Atom first (NTV uses Atom format)
	var atom Atom
	err1 := xml.Unmarshal(body, &atom)
	if err1 == nil && len(atom.Entries) > 0 {
		
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
		}
		
		return newsItems, nil
	}

	// If Atom parsing failed, try RSS
	var rss RSS
	err2 := xml.Unmarshal(body, &rss)
	if err2 == nil && len(rss.Channel.Items) > 0 {
		
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
			
		}
		
		return newsItems, nil
	}

	return []NewsItem{}, nil
}


