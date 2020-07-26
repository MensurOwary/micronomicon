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

var microns = scraper.Database()

func GetARandomMicronForTag(tag tag.Tag) *Micron {
	rows := microns[tag.Name]

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
