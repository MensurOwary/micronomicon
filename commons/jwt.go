package commons

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type JwtService interface {
	SignedToken(username string) (string, error)
	SaveJwt(jwt string) bool
	DoesJwtExist(jwt string) bool
	DeleteJwt(jwt string) bool
}

type jwtService struct {
	db *mongo.Client
}

func NewJwtService(mongo *mongo.Client) JwtService {
	return &jwtService{
		db: mongo,
	}
}

func (j *jwtService) SignedToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	return token.SignedString([]byte(Config.JwtSecret))
}

func (j *jwtService) SaveJwt(jwt string) bool {
	return j.withJwtDb(func(collection *mongo.Collection) interface{} {
		insertOneResult, err := collection.InsertOne(context.Background(), bson.D{
			{"jwt", jwt},
		})
		if err != nil {
			log.Printf("Error while saving Jwt to database - %s", err.Error())
		} else {
			log.Printf("Jwt saved to database - %s", insertOneResult)
		}
		return err == nil
	}).(bool)
}

func (j *jwtService) DoesJwtExist(jwt string) bool {
	return j.withJwtDb(func(collection *mongo.Collection) interface{} {
		findOne := collection.FindOne(context.Background(),
			bson.D{
				{"jwt", jwt},
			})

		return findOne.Err() == nil
	}).(bool)
}

func (j *jwtService) DeleteJwt(jwt string) bool {
	return j.withJwtDb(func(collection *mongo.Collection) interface{} {
		_, err := collection.DeleteOne(context.Background(), bson.D{
			{"jwt", jwt},
		})

		if err != nil {
			log.Printf("Error while deleting a Jwt token - " + err.Error())
		}
		return err == nil
	}).(bool)
}

type dbAction func(collection *mongo.Collection) interface{}

func (j *jwtService) withJwtDb(action dbAction) interface{} {
	collection := j.db.Database("micron").Collection("jwts")
	return action(collection)
}
