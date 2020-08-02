package tag

import "micron/scraper"

type Tag struct {
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
	Size int   `json:"size"`
}

type Service interface {
	AddTagsForUser(username string, newTagIds []string) bool
	RemoveTagsFromUser(username string, removable []string) bool
	GetUserTags(username string) []string
	GetTagById(name string) *Tag
}

type Repository struct {
	database scraper.Scraper
}

func NewRepository(database scraper.Scraper) *Repository {
	return &Repository{
		database: database,
	}
}
