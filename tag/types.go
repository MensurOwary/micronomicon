package tag

import (
	"go.mongodb.org/mongo-driver/mongo"
	"micron/scraper"
)

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
	GetTagById(name string) Tag
}

type tagsService struct {
	db    *mongo.Client
	tagDb *Repository
}

func NewService(mongo *mongo.Client, repository *Repository) Service {
	return &tagsService{
		db:    mongo,
		tagDb: repository,
	}
}


type Repository struct {
	database scraper.Scraper
}

func NewRepository(database scraper.Scraper) *Repository {
	return &Repository{
		database: database,
	}
}
