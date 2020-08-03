package commons

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ParsedJwtResult struct {
	Username    string
	parsedToken *jwt.Token
	err         error
}

func (p *ParsedJwtResult) isOk() bool {
	return p.err == nil
}

func (p *ParsedJwtResult) IsJwtValid() error {
	claims, ok := p.parsedToken.Claims.(jwt.MapClaims)

	if !p.isOk() {
		return errors.New("jwt error occurred : " + p.err.Error())
	}

	if !(ok && p.parsedToken.Valid) {
		return errors.New("jwt error occurred : parsed token is invalid")
	}

	if username, usernameOk := claims["username"]; usernameOk {
		p.Username = username.(string)
		return nil
	}

	return errors.New("jwt error occurred : username is invalid")
}

// JwtService deals with jwt related actions
type JwtService struct {
	db *mongo.Client
}

// NewJwtService creates a new instance of the service
func NewJwtService(mongo *mongo.Client) *JwtService {
	return &JwtService{
		db: mongo,
	}
}

// ParseJwt parses and validates the given jwt token
func (j *JwtService) ParseJwt(rawJwt string) (*ParsedJwtResult, error) {
	parseJwt := j.doParseJwt(rawJwt)

	if err := parseJwt.IsJwtValid(); err != nil {
		return nil, err
	}
	return parseJwt, nil
}

func (j *JwtService) doParseJwt(rawJwt string) *ParsedJwtResult {
	parsedToken, err := jwt.Parse(rawJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt error: signing method is wrong")
		}
		return []byte(Config.JwtSecret), nil
	})

	if err == nil && !j.DoesJwtExist(rawJwt) {
		return &ParsedJwtResult{
			parsedToken: parsedToken,
			err:         errors.New("token expired"),
		}
	}

	return &ParsedJwtResult{
		parsedToken: parsedToken,
		err:         err,
	}
}

// SignedToken creates a signed token
func (j *JwtService) SignedToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(Config.JwtTokenExpiresIn).Unix(),
		"iat":      time.Now().Unix(),
	})
	return token.SignedString([]byte(Config.JwtSecret))
}

// SaveJwt saves the token to the database
func (j *JwtService) SaveJwt(jwt string) bool {
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

// DoesJwtExist checks the existence of the token in the database
func (j *JwtService) DoesJwtExist(jwt string) bool {
	return j.withJwtDb(func(collection *mongo.Collection) interface{} {
		findOne := collection.FindOne(context.Background(),
			bson.D{
				{"jwt", jwt},
			})

		return findOne.Err() == nil
	}).(bool)
}

// DeleteJwt deletes the jwt token from the database, hence invalidating it
func (j *JwtService) DeleteJwt(jwt string) bool {
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

func (j *JwtService) withJwtDb(action dbAction) interface{} {
	collection := j.db.Database("micron").Collection("jwts")
	return action(collection)
}
