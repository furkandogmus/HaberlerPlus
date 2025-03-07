package impl

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// RSSItem represents a single item in an RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        struct {
		Value       string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
}

// RSSChannel represents an RSS channel
type RSSChannel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

// RSS represents an RSS feed
type RSS struct {
	XMLName xml.Name    `xml:"rss"`
	Channel RSSChannel  `xml:"channel"`
}

// AtomLink represents a link in an Atom feed
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// AtomEntry represents an entry in an Atom feed
type AtomEntry struct {
	Title   string    `xml:"title"`
	Link    AtomLink  `xml:"link"`
	Content string    `xml:"content"`
	ID      string    `xml:"id"`
}

// Atom represents an Atom feed
type Atom struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Link    AtomLink    `xml:"link"`
	Entries []AtomEntry `xml:"entry"`
}

// RSSSource implements a news source for RSS feeds
type RSSSource struct {
	name       string
	categories []string
	feedURLs   map[string]string
}

// NewRSSSource creates a new RSSSource instance
func NewRSSSource(name string, categories []string, feedURLs map[string]string) *RSSSource {
	return &RSSSource{
		name:       name,
		categories: categories,
		feedURLs:   feedURLs,
	}
}

// Name returns the name of the news source
func (r *RSSSource) Name() string {
	return r.name
}

// Categories returns the available categories for this news source
func (r *RSSSource) Categories() []string {
	return r.categories
}

// FetchNews fetches news items for the specified category
func (r *RSSSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	// Adjust for 0-based indexing
	if categoryIndex >= 0 && categoryIndex < len(r.categories) {
		// This is 0-based indexing, convert to 1-based for our internal use
		categoryIndex++
	}

	if categoryIndex < 1 || categoryIndex > len(r.categories) {
		fmt.Printf("Debug: Invalid category index: %d, returning empty results\n", categoryIndex)
		return []NewsItem{}, nil
	}

	category := r.categories[categoryIndex-1]
	feedURL, ok := r.feedURLs[category]
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

	// Try to parse as RSS first
	var rss RSS
	err1 := xml.Unmarshal(body, &rss)
	if err1 == nil && len(rss.Channel.Items) > 0 {
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
			if item.Title == "" || item.Link == "" {
				continue
			}
			
			// Clean up the title
			title := strings.TrimSpace(item.Title)
			title = strings.TrimPrefix(title, "<![CDATA[")
			title = strings.TrimSuffix(title, "]]>")
			
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

	// If RSS parsing failed, try Atom
	var atom Atom
	err2 := xml.Unmarshal(body, &atom)
	if err2 == nil && len(atom.Entries) > 0 {
		fmt.Printf("Debug: Successfully parsed as Atom, found %d entries\n", len(atom.Entries))
		
		// Process Atom entries
		maxItems := 30
		if len(atom.Entries) < maxItems {
			maxItems = len(atom.Entries)
		}
		
		newsItems := make([]NewsItem, 0, maxItems)
		for i := 0; i < maxItems; i++ {
			entry := atom.Entries[i]
			
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

	// If we get here, both parsing attempts failed
	fmt.Printf("Debug: Failed to parse as RSS: %v\n", err1)
	fmt.Printf("Debug: Failed to parse as Atom: %v\n", err2)
	fmt.Printf("Debug: Kategori: %s, Kategori Numarası: %d\n", category, categoryIndex)
	fmt.Printf("Debug: Haber bulunamadı.\n")
	return []NewsItem{}, nil
}

