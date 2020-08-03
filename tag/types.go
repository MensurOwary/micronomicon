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

type TagsService struct {
	db    *mongo.Client
	tagDb *Repository
}

func NewService(mongo *mongo.Client, repository *Repository) *TagsService {
	return &TagsService{
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
