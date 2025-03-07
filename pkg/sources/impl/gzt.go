package impl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GztSource implements the original gzt.com news source
type GztSource struct{}

// NewGztSource creates a new GztSource instance
func NewGztSource() *GztSource {
	return &GztSource{}
}

// Name returns the name of the news source
func (g *GztSource) Name() string {
	return "GZT.com"
}

// Categories returns the available categories for this news source
func (g *GztSource) Categories() []string {
	return []string{"POLITIKA", "DUNYA", "EKONOMI", "BILIM", "GUNCEL", "AKTUEL-KULTUR", "SAGLIK"}
}

// FetchNews fetches news items for the specified category
func (g *GztSource) FetchNews(categoryIndex int) ([]NewsItem, error) {
	categories := g.Categories()
	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return nil, fmt.Errorf("invalid category index: %d", categoryIndex)
	}

	category := categories[categoryIndex]
	url := fmt.Sprintf("https://www.gzt.com/%s", strings.ToLower(category))

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var newsItems []NewsItem

	doc.Find(".feed-card-content.news-card-content").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		text = strings.Split(text, "…..devamı")[0]
		firstLink := s.Find("a").First()
		href, exists := firstLink.Attr("href")
		if exists {
			newsItems = append(newsItems, NewsItem{
				Title: text,
				URL:   fmt.Sprintf("https://www.gzt.com%s", href),
			})
		}
	})

	return newsItems, nil
}

// NewsItem represents a single news item
type NewsItem struct {
	Title string
	URL   string
}