package micron

import (
	"math/rand"
	"micron/scraper"
	"micron/tag"
)

type Micron struct {
	id    string
	Url   string
	Title string
	tag   tag.Tag
}

type Service interface {
	GetARandomMicronForTag(tag tag.Tag) *Micron
}

type service struct {
	microns scraper.Scraper
}

func NewService(scraper scraper.Scraper) Service {
	return &service{
		microns: scraper,
	}
}

func (s *service) GetARandomMicronForTag(tag tag.Tag) *Micron {
	rows := s.microns.Database()[tag.Name]

	if rows != nil {
		chosen := rand.Intn(len(rows))
		row := &rows[chosen]
		return &Micron{
			Url:   row.Link,
			Title: row.Title,
			tag:   tag,
		}
	}

	return nil
}
