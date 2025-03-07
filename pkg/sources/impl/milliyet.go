package impl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// MilliyetSource implements a news source for milliyet.com.tr
type MilliyetSource struct{}

// NewMilliyetSource creates a new MilliyetSource instance
func NewMilliyetSource() *MilliyetSource {
	return &MilliyetSource{}
}

// Name returns the name of the news source
func (m *MilliyetSource) Name() string {
	return "Milliyet.com.tr"
}

// Categories returns the available categories for this news source
func (m *MilliyetSource) Categories() []string {
	return []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "TEKNOLOJİ"}
}

// FetchNews fetches news items for the specified category
func (m *MilliyetSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	categories := m.Categories()
	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return nil, fmt.Errorf("invalid category index: %d", categoryIndex)
	}

	categoryMap := map[string]string{
		"GÜNDEM":    "gundem",
		"DÜNYA":     "dunya",
		"EKONOMİ":   "ekonomi",
		"SPOR":      "spor",
		"TEKNOLOJİ": "teknoloji",
	}

	category := categories[categoryIndex]
	categoryPath := categoryMap[category]
	url := fmt.Sprintf("https://www.milliyet.com.tr/%s/", categoryPath)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var newsItems []NewsItem

	// First try to get news from the slider (featured news)
	doc.Find(".cat-slider__link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			titleElement := s.Find(".cat-slider__title")
			title := strings.TrimSpace(titleElement.Text())
			
			if title != "" {
				fullURL := href
				if !strings.HasPrefix(href, "http") {
					fullURL = fmt.Sprintf("https://www.milliyet.com.tr%s", href)
				}
				newsItems = append(newsItems, NewsItem{
					Title: title,
					URL:   fullURL,
				})
			}
		}
	})

	// Then get news from the category cards
	doc.Find(".category-card").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			titleElement := s.Find(".category-card__head")
			title := strings.TrimSpace(titleElement.Text())
			
			if title != "" {
				fullURL := href
				if !strings.HasPrefix(href, "http") {
					fullURL = fmt.Sprintf("https://www.milliyet.com.tr%s", href)
				}
				newsItems = append(newsItems, NewsItem{
					Title: title,
					URL:   fullURL,
				})
			}
		}
	})

	// If we still don't have any news, try the cat-list-card items
	if len(newsItems) == 0 {
		doc.Find(".cat-list-card__link").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				titleElement := s.Find(".cat-list-card__title")
				title := strings.TrimSpace(titleElement.Text())
				
				if title != "" {
					fullURL := href
					if !strings.HasPrefix(href, "http") {
						fullURL = fmt.Sprintf("https://www.milliyet.com.tr%s", href)
					}
					newsItems = append(newsItems, NewsItem{
						Title: title,
						URL:   fullURL,
					})
				}
			}
		})
	}

	return newsItems, nil
}