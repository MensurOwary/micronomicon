package tag

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"micron/commons"
)

func (t *Repository) GetAvailableTags() Tags {
	database := t.database.Database()
	keys := make([]Tag, 0, len(database))
	for tag := range database {
		keys = append(keys, Tag{
			Name: tag,
		})
	}
	return Tags{
		Tags: keys,
		Size: len(keys),
	}
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

func (t *tagsService) GetTagById(name string) *Tag {
	for _, tag := range t.tagDb.GetAvailableTags().Tags {
		if tag.Name == name {
			return &tag
		}
	}
	return nil
}

func (t *tagsService) AddTagsForUser(username string, newTagIds []string) bool {
	return t.withUsersDb(func(collection *mongo.Collection) interface{} {
		currentUserTags := t.GetUserTags(username)
		union := commons.Union(currentUserTags, newTagIds)
		return t.doUpdateTagsField(collection, username, union)
	}).(bool)
}

func (t *tagsService) RemoveTagsFromUser(username string, removable []string) bool {
	return t.withUsersDb(func(collection *mongo.Collection) interface{} {
		currentUserTags := t.GetUserTags(username)
		diff := commons.Difference(currentUserTags, removable)
		return t.doUpdateTagsField(collection, username, diff)
	}).(bool)
}

func (t *tagsService) doUpdateTagsField(collection *mongo.Collection, username string, tags []string) bool {
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

func (t *tagsService) GetUserTags(username string) []string {
	return t.withUsersDb(func(collection *mongo.Collection) interface{} {
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

type dbAction func(collection *mongo.Collection) interface{}

func (t *tagsService) withUsersDb(action dbAction) interface{} {
	collection := t.db.Database("micron").Collection("users")
	return action(collection)
}

func field(username string) bson.D {
	return bson.D{
		{"username", username},
	}
}
