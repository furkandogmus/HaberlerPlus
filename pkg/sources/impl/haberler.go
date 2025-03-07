package impl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HaberlerComSource implements a news source for haberler.com
type HaberlerComSource struct{}

// NewHaberlerComSource creates a new HaberlerComSource instance
func NewHaberlerComSource() *HaberlerComSource {
	return &HaberlerComSource{}
}

// Name returns the name of the news source
func (h *HaberlerComSource) Name() string {
	return "Haberler.com"
}

// Categories returns the available categories for this news source
func (h *HaberlerComSource) Categories() []string {
	return []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "SAĞLIK", "TEKNOLOJİ"}
}

// FetchNews fetches news items for the specified category
func (h *HaberlerComSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	categories := h.Categories()
	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return nil, fmt.Errorf("invalid category index: %d", categoryIndex)
	}

	categoryMap := map[string]string{
		"GÜNDEM":    "guncel",
		"DÜNYA":     "dunya",
		"EKONOMİ":   "ekonomi",
		"SPOR":      "spor",
		"SAĞLIK":    "saglik",
		"TEKNOLOJİ": "teknoloji",
	}

	category := categories[categoryIndex]
	categoryPath := categoryMap[category]
	url := fmt.Sprintf("https://www.haberler.com/%s/", categoryPath)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var newsItems []NewsItem
	// Map to track URLs we've already added to avoid duplicates
	seenURLs := make(map[string]bool)
	// Map to track titles we've already added to avoid similar content
	seenTitles := make(map[string]bool)

	// Helper function to add a news item if it's not a duplicate
	addNewsItem := func(title, href string) {
		if title == "" || href == "" {
			return
		}

		// Normalize URL
		fullURL := href
		if !strings.HasPrefix(href, "http") {
			fullURL = fmt.Sprintf("https://www.haberler.com%s", href)
		}

		// Skip if we've already seen this URL
		if seenURLs[fullURL] {
			return
		}

		// Skip if we've already seen this title or a very similar one
		if seenTitles[title] {
			return
		}

		// Check if this is from the correct category
		if !strings.Contains(fullURL, categoryPath) {
			return
		}

		// Add the news item
		newsItems = append(newsItems, NewsItem{
			Title: title,
			URL:   fullURL,
		})

		// Mark as seen
		seenURLs[fullURL] = true
		seenTitles[title] = true
	}

	// First try to get news from the main slider (featured news)
	doc.Find(".new3slide").Each(func(i int, s *goquery.Selection) {
		linkElement := s.Find("a")
		href, exists := linkElement.Attr("href")
		
		if exists {
			// Get the title from h2 inside the caption
			titleElement := s.Find(".new3caption h2")
			title := strings.TrimSpace(titleElement.Text())
			
			addNewsItem(title, href)
		}
	})

	// Then get news from the card grid (main content area)
	doc.Find(".new3card").Each(func(i int, s *goquery.Selection) {
		linkElement := s.Find("a")
		href, exists := linkElement.Attr("href")
		
		if exists {
			// Get the title from h3 inside the card body
			titleElement := s.Find(".new3card-body h3")
			title := strings.TrimSpace(titleElement.Text())
			
			addNewsItem(title, href)
		}
	})

	// If we still don't have enough news, try the older format
	if len(newsItems) < 10 {
		doc.Find("div.hblnBox, article.box, .news-item").Each(func(i int, s *goquery.Selection) {
			linkElement := s.Find("a")
			href, exists := linkElement.Attr("href")
			
			if exists {
				titleElement := s.Find("a.hblnTitle, h3, .news-title")
				title := strings.TrimSpace(titleElement.Text())
				
				addNewsItem(title, href)
			}
		})
	}

	// If we still don't have enough news, try a more targeted approach
	if len(newsItems) < 10 {
		// Look specifically for news items with the category in the URL
		doc.Find("a[href*='" + categoryPath + "']").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				title := strings.TrimSpace(s.Text())
				// Filter out very short or very long titles
				if title != "" && len(title) > 10 && len(title) < 200 {
					addNewsItem(title, href)
				}
			}
		})
	}

	// Add debug information if no news items were found
	if len(newsItems) == 0 {
		fmt.Printf("Debug: Kategori: %s, Kategori Numarası: %d\n", category, categoryIndex)
		fmt.Printf("Debug: Haber bulunamadı.\n")
	}

	return newsItems, nil
}