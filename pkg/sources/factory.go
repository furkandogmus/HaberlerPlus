package sources

import (
	"github.com/furkandogmus/HaberlerPlus/pkg/sources/impl"
)

// Factory functions for creating news sources
func NewGztSource() NewsSource {
	return impl.NewGztSource()
}

func NewHurriyetSource() NewsSource {
	return impl.NewHurriyetSource()
}

func NewSozcuSource() NewsSource {
	return impl.NewSozcuSource()
}

func NewMilliyetSource() NewsSource {
	return impl.NewMilliyetSource()
}

func NewHaberlerComSource() NewsSource {
	return impl.NewHaberlerComSource()
}

// Factory functions for RSS-based sources
func NewCNNTurkSource() NewsSource {
	return impl.NewCNNTurkSource()
}

func NewNTVSource() NewsSource {
	return impl.NewNTVSource()
}

func NewHaberturkSource() NewsSource {
	return impl.NewHaberturkSource()
} 