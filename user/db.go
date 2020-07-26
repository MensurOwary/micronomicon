package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"micron/commons"
)

var MongoClient *mongo.Client

func DoSaveUser(user *User) bool {
	return withUsersDb(func(collection *mongo.Collection) interface{} {
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

func DoFindUser(username string) *User {
	return withUsersDb(func(collection *mongo.Collection) interface{} {
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

func DoAddTagsForUser(username string, newTagIds []string) bool {
	return withUsersDb(func(collection *mongo.Collection) interface{} {
		currentUserTags := DoGetUserTags(username)
		union := commons.Union(currentUserTags, newTagIds)
		return doUpdateTagsField(collection, username, union)
	}).(bool)
}

func DoRemoveTagsFromUser(username string, removable []string) bool {
	return withUsersDb(func(collection *mongo.Collection) interface{} {
		currentUserTags := DoGetUserTags(username)
		diff := commons.Difference(currentUserTags, removable)
		return doUpdateTagsField(collection, username, diff)
	}).(bool)
}

func doUpdateTagsField(collection *mongo.Collection, username string, tags []string) bool {
	updateResult, err := collection.UpdateOne(
		context.Background(),
		field(username),
		bson.D{
			{"$set", bson.D{
				{"tags", tags},
			}},
		},
		(&options.UpdateOptions{}).SetUpsert(true),
	)
	if err != nil {
		log.Printf("Error occurred during updating of %s - %s\n", username, err)
		return false
	} else {
		log.Printf("Update successful - Result = %+v\n", updateResult)
		return true
	}
}

func DoGetUserTags(username string) []string {
	return withUsersDb(func(collection *mongo.Collection) interface{} {
		singleResult := collection.FindOne(context.Background(),
			field(username),
			(&options.FindOneOptions{}).SetProjection(bson.D{
				{"_id", 0},
				{"tags", 1},
			}))
		var tags struct {
			Tags []string
		}
		err := singleResult.Decode(&tags)

		if err != nil && tags.Tags != nil {
			log.Printf("Error occurred while decoding the response - %s", err)
		}
		return tags.Tags
	}).([]string)
}

func field(username string) bson.D {
	return bson.D{
		{"username", username},
	}
}

type DbAction func(collection *mongo.Collection) interface{}

func withUsersDb(action DbAction) interface{} {
	collection := MongoClient.Database("micron").Collection("users")
	return action(collection)
}
