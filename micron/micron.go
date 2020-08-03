package micron

import (
	"math/rand"
	"micron/scraper"
	"micron/tag"
)

// Represents a micron - a resource
type Micron struct {
	id    string
	URL   string
	Title string
	tag   tag.Tag
}

// Deals with Micron related interactions
type Service struct {
	microns scraper.Scraper
}

// Creates a new instance of micron service
func NewService(scraper scraper.Scraper) *Service {
	return &Service{
		microns: scraper,
	}
}

// Represents a non-existent resource
var EmptyMicron = Micron{}

// Fetches a random micron given the tag
func (s *Service) GetARandomMicronForTag(tag tag.Tag) Micron {
	rows := s.microns.Database()[tag.Name]

	if rows == nil {
		return EmptyMicron
	}

	chosen := rand.Intn(len(rows))
	row := &rows[chosen]
	return Micron{
		URL:   row.Link,
		Title: row.Title,
		tag:   tag,
	}
}
