package tag

import (
	"micron/scraper"
)

type Tag struct {
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
	Size int   `json:"size"`
}

func GetAvailableTags() Tags {
	database := scraper.Database()
	keys := make([]Tag, 0, len(database))
	for tag, _ := range database {
		keys = append(keys, Tag{
			Name: tag,
		})
	}
	return Tags{
		Tags: keys,
		Size: len(keys),
	}
}

func GetTagById(name string) *Tag {
	for _, tag := range GetAvailableTags().Tags {
		if tag.Name == name {
			return &tag
		}
	}
	return nil
}
