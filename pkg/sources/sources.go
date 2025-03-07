package sources

import (
	"github.com/furkandogmus/HaberlerPlus/pkg/sources/impl"
)

// NewsSource defines the interface for all news sources
type NewsSource interface {
	Name() string
	Categories() []string
	FetchNews(categoryIndex int) ([]NewsItem, error)
}

// NewsItem is an alias for impl.NewsItem
type NewsItem = impl.NewsItem

// GetAllSources returns all available news sources
func GetAllSources() []NewsSource {
	return []NewsSource{
		NewGztSource(),
		NewHurriyetSource(),
		NewSozcuSource(),
		NewMilliyetSource(),
		NewHaberlerComSource(),
	}
}