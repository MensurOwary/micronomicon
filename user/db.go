package user

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SaveUser saves the user to database
func (r *repository) SaveUser(user User) bool {
	return r.databaseAccess.withUsersDb(func(collection *mongo.Collection) interface{} {
		result, err := collection.InsertOne(
			context.Background(),
			bson.D{
				{"username", user.Username},
				{"password", user.Password},
				{"email", user.Email},
				{"name", user.Name},
			},
		)
		if err != nil {
			log.Errorf("Saving user failed [%s] : %s", user, err)
		} else {
			log.Infof("Saving user successful : InsertionId : %s", result.InsertedID)
		}
		return err == nil
	}).(bool)
}

// DoesNotExist represents a non-existent resource
var DoesNotExist = User{}

// FindUser finds the user from database based on the username
func (r *repository) FindUser(username string) User {
	return r.databaseAccess.withUsersDb(func(collection *mongo.Collection) interface{} {
		cursor := collection.FindOne(
			context.Background(),
			field(username),
		)
		user := User{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Errorf("Unmarshalling the user object from response failed [%s] : %s", username, err)
			return DoesNotExist
		}
		return user
	}).(User)
}

func field(username string) bson.D {
	return bson.D{
		{"username", username},
	}
}

type dbAction func(collection *mongo.Collection) interface{}

func (d *databaseAccess) withUsersDb(action dbAction) interface{} {
	collection := d.db.Database("micron").Collection("users")
	return action(collection)
}
