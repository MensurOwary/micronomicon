package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (r *repository) SaveUser(user *User) bool {
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
			log.Printf("Error occurred during insertion of %s - %s", user, err)
		} else {
			log.Printf("Insertion successful - InsertionId = %s", result.InsertedID)
		}
		return err == nil
	}).(bool)
}

func (r *repository) FindUser(username string) *User {
	return r.databaseAccess.withUsersDb(func(collection *mongo.Collection) interface{} {
		cursor := collection.FindOne(
			context.Background(),
			field(username),
		)
		user := User{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Printf("Error occurred during retrieval of %s - %s", username, err)
			return nil
		}
		return &user
	}).(*User)
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
