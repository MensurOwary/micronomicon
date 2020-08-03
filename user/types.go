package user

import (
	"go.mongodb.org/mongo-driver/mongo"
	"micron/commons"
	"micron/tag"
)

// Represents a user payload
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Represents a domain user
type Account struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Tags     tag.Tags `json:"tags"`
}

// Represents a means that deals with Jwt related actions
type JwtService interface {
	SignedToken(username string) (string, error)
	SaveJwt(jwt string) bool
	DeleteJwt(jwt string) bool
}

// Represents a means that deals with Tags related actions
type TagsService interface {
	AddTagsForUser(username string, newTagIds []string) bool
	RemoveTagsFromUser(username string, removable []string) bool
	GetUserTags(username string) []string
	GetTagByID(name string) tag.Tag
}

// Represents the user service entity
type Service struct {
	store   Repository
	tags    TagsService
	encrypt commons.EncryptService
	jwt     JwtService
}

// Configuration for the user service entity
type ServiceConfig struct {
	Store   Repository
	Tags    TagsService
	Encrypt commons.EncryptService
	Jwt     JwtService
}

type databaseAccess struct {
	db *mongo.Client
}

type repository struct {
	databaseAccess
}

// Creates a new instance of user repository
func NewRepository(MongoClient *mongo.Client) Repository {
	return &repository{
		databaseAccess{
			db: MongoClient,
		},
	}
}

// Creates a new instance of user service
func NewService(config ServiceConfig) *Service {
	return &Service{
		store:   config.Store,
		tags:    config.Tags,
		encrypt: config.Encrypt,
		jwt:     config.Jwt,
	}
}

// Represents a means that deals with database related interactions
type Repository interface {
	SaveUser(user User) bool
	FindUser(username string) User
}
