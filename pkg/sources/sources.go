package sources

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// NewsSource defines the interface for all news sources
type NewsSource interface {
	Name() string
	Categories() []string
	FetchNews(categoryIndex int) ([]NewsItem, error)
}

// NewsItem represents a single news item
type NewsItem struct {
	Title string
	URL   string
}

// GztSource implements the original gzt.com news source
type GztSource struct{}

func (g *GztSource) Name() string {
	return "GZT.com"
}

func (g *GztSource) Categories() []string {
	return []string{"POLITIKA", "DUNYA", "EKONOMI", "BILIM", "GUNCEL", "AKTUEL-KULTUR", "SAGLIK"}
}

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

// HurriyetSource implements a news source for hurriyet.com.tr
type HurriyetSource struct{}

func (h *HurriyetSource) Name() string {
	return "Hurriyet.com.tr"
}

func (h *HurriyetSource) Categories() []string {
	return []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "TEKNOLOJİ", "SAĞLIK", "YAŞAM"}
}

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

	doc.Find("div.news-card").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("h3.news-card__title")
		title := strings.TrimSpace(titleElement.Text())
		
		linkElement := s.Find("a.news-card__link")
		href, exists := linkElement.Attr("href")
		if exists && title != "" {
			fullURL := href
			if !strings.HasPrefix(href, "http") {
				fullURL = fmt.Sprintf("https://www.hurriyet.com.tr%s", href)
			}
			newsItems = append(newsItems, NewsItem{
				Title: title,
				URL:   fullURL,
			})
		}
	})

	return newsItems, nil
}

// SozcuSource implements a news source for sozcu.com.tr
type SozcuSource struct{}

func (s *SozcuSource) Name() string {
	return "Sozcu.com.tr"
}

func (s *SozcuSource) Categories() []string {
	return []string{"GÜNDEM", "DÜNYA", "EKONOMİ", "SPOR", "FİNANS", "YAŞAM", "SAĞLIK"}
}

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

	doc.Find("div.news-item").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("h3.news-item__title")
		title := strings.TrimSpace(titleElement.Text())
		
		linkElement := s.Find("a.news-item__link")
		href, exists := linkElement.Attr("href")
		if exists && title != "" {
			newsItems = append(newsItems, NewsItem{
				Title: title,
				URL:   href,
			})
		}
	})

	return newsItems, nil
}

// GetAllSources returns all available news sources
func GetAllSources() []NewsSource {
	return []NewsSource{
		&GztSource{},
		&HurriyetSource{},
		&SozcuSource{},
	}
} 