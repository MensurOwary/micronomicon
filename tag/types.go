package tag

import (
	"go.mongodb.org/mongo-driver/mongo"
	"micron/scraper"
)

// Tag represents a single tag
type Tag struct {
	Name string `json:"name"`
}

// Tags represents a collection of tags
type Tags struct {
	Tags []Tag `json:"tags"`
	Size int   `json:"size"`
}

// TagsService deals with tags related interactions
type TagsService struct {
	db    *mongo.Client
	tagDb *Repository
}

// NewService initializes a new TagsService
func NewService(mongo *mongo.Client, repository *Repository) *TagsService {
	return &TagsService{
		db:    mongo,
		tagDb: repository,
	}
}

// Repository represents means to interact with database
type Repository struct {
	database scraper.Scraper
}

// NewRepository initializes a new Repository
func NewRepository(database scraper.Scraper) *Repository {
	return &Repository{
		database: database,
	}
}
