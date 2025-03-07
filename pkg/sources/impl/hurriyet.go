package impl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HurriyetSource implements a news source for hurriyet.com.tr
type HurriyetSource struct{}

// NewHurriyetSource creates a new HurriyetSource instance
func NewHurriyetSource() *HurriyetSource {
	return &HurriyetSource{}
}

// Name returns the name of the news source
func (h *HurriyetSource) Name() string {
	return "Hurriyet.com.tr"
}

// Categories returns the available categories for this news source
func (h *HurriyetSource) Categories() []string {
	return []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "TEKNOLOJİ", "SAĞLIK", "YAŞAM"}
}

// FetchNews fetches news items for the specified category
func (h *HurriyetSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	categories := h.Categories()
	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return nil, fmt.Errorf("invalid category index: %d", categoryIndex)
	}

	categoryMap := map[string]string{
		"GÜNDEM":    "gundem",
		"DÜNYA":     "dunya",
		"EKONOMİ":   "ekonomi",
		"SPOR":      "spor",
		"TEKNOLOJİ": "teknoloji",
		"SAĞLIK":    "saglik",
		"YAŞAM":     "yasam",
	}

	category := categories[categoryIndex]
	categoryPath := categoryMap[category]
	url := fmt.Sprintf("https://www.hurriyet.com.tr/%s/", categoryPath)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var newsItems []NewsItem

	// Extract news items from category list
	doc.Find("div.category__list__item").Each(func(i int, s *goquery.Selection) {
		// Find the title element
		titleElement := s.Find("h2")
		title := strings.TrimSpace(titleElement.Text())
		
		// Find the link element
		linkElement := s.Find("a[href]").First()
		href, exists := linkElement.Attr("href")
		
		if exists && title != "" {
			fullURL := href
			if !strings.HasPrefix(href, "http") {
				fullURL = fmt.Sprintf("https://www.hurriyet.com.tr%s", href)
			}
			
			// Add the news item to our list
			newsItems = append(newsItems, NewsItem{
				Title: title,
				URL:   fullURL,
			})
		}
	})

	// If no news items were found, try an alternative approach
	if len(newsItems) == 0 {
		fmt.Printf("Debug: Kategori: %s, Kategori Numarası: %d\n", category, categoryIndex)
		fmt.Printf("Debug: Haber bulunamadı.\n")
	}

	return newsItems, nil
}