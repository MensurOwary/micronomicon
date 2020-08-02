package user

import (
	"go.mongodb.org/mongo-driver/mongo"
	"micron/commons"
	"micron/tag"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Account struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Tags     tag.Tags `json:"tags"`
}

type Service struct {
	store   Repository
	tags    tag.Service
	encrypt commons.EncryptService
	jwt     commons.JwtService
}

type ServiceConfig struct {
	Store   Repository
	Tags    tag.Service
	Encrypt commons.EncryptService
	Jwt     commons.JwtService
}

type databaseAccess struct {
	db *mongo.Client
}

type repository struct {
	databaseAccess
}

func NewRepository(MongoClient *mongo.Client) Repository {
	return &repository{
		databaseAccess{
			db: MongoClient,
		},
	}
}

func NewService(config ServiceConfig) *Service {
	return &Service{
		store:   config.Store,
		tags:    config.Tags,
		encrypt: config.Encrypt,
		jwt:     config.Jwt,
	}
}

type Repository interface {
	SaveUser(user User) bool
	FindUser(username string) User
}
