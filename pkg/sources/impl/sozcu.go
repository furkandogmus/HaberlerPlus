package impl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SozcuSource implements a news source for sozcu.com.tr
type SozcuSource struct{}

// NewSozcuSource creates a new SozcuSource instance
func NewSozcuSource() *SozcuSource {
	return &SozcuSource{}
}

// Name returns the name of the news source
func (s *SozcuSource) Name() string {
	return "Sozcu.com.tr"
}

// Categories returns the available categories for this news source
func (s *SozcuSource) Categories() []string {
	return []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "FİNANS", "YAŞAM", "SAĞLIK"}
}

// FetchNews fetches news items for the specified category
func (s *SozcuSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	categories := s.Categories()
	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return nil, fmt.Errorf("invalid category index: %d", categoryIndex)
	}

	categoryMap := map[string]string{
		"GÜNDEM":  "gundem",
		"DÜNYA":   "dunya",
		"EKONOMİ": "ekonomi",
		"SPOR":    "spor",
		"FİNANS":  "finans",
		"YAŞAM":   "yasam",
		"SAĞLIK":  "saglik",
	}

	category := categories[categoryIndex]
	categoryPath := categoryMap[category]
	url := fmt.Sprintf("https://www.sozcu.com.tr/%s/", categoryPath)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var newsItems []NewsItem

	// Extract news items from the list-content section
	doc.Find(".list-content .row").Each(func(i int, s *goquery.Selection) {
		// Find the link element
		linkElement := s.Find("a")
		href, exists := linkElement.Attr("href")
		
		// Find the title element
		titleElement := s.Find("span.d-block.fs-5.fw-semibold")
		title := strings.TrimSpace(titleElement.Text())
		
		if exists && title != "" {
			// Add the news item to our list
			newsItems = append(newsItems, NewsItem{
				Title: title,
				URL:   href,
			})
		}
	})

	// If no news items were found, try an alternative approach
	if len(newsItems) == 0 {
		// Try to find news items in other sections
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.Contains(href, categoryPath) {
				title := strings.TrimSpace(s.Text())
				if title != "" && len(title) > 10 && len(title) < 200 {
					newsItems = append(newsItems, NewsItem{
						Title: title,
						URL:   href,
					})
				}
			}
		})
	}

	// If still no news items were found, log debug info
	if len(newsItems) == 0 {
		fmt.Printf("Debug: Kategori: %s, Kategori Numarası: %d\n", category, categoryIndex)
		fmt.Printf("Debug: Haber bulunamadı.\n")
	}

	return newsItems, nil
}