package tag

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"micron/commons"
)

// AddTagsForUser handles adding new tags for user
func (t *TagsService) AddTagsForUser(username string, newTagIds []string) bool {
	return t.withUsersDb(func(collection *mongo.Collection) interface{} {
		currentUserTags := t.GetUserTags(username)
		union := commons.Union(currentUserTags, newTagIds)
		return t.doUpdateTagsField(collection, username, union)
	}).(bool)
}

// RemoveTagsFromUser handles removing tags from user
func (t *TagsService) RemoveTagsFromUser(username string, removable []string) bool {
	return t.withUsersDb(func(collection *mongo.Collection) interface{} {
		currentUserTags := t.GetUserTags(username)
		diff := commons.Difference(currentUserTags, removable)
		return t.doUpdateTagsField(collection, username, diff)
	}).(bool)
}

func (t *TagsService) doUpdateTagsField(collection *mongo.Collection, username string, tags []string) bool {
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
		log.Errorf("Error occurred during updating of %s - %s", username, err)
	} else {
		log.Infof("Updating tags of user[%s] was successful - Result = %+v", username, updateResult)
	}
	return err == nil
}

// GetUserTags fetches the user tags
func (t *TagsService) GetUserTags(username string) []string {
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

		if err != nil && tags.Tags == nil {
			log.Errorf("Error occurred while decoding the user tags - %s", err)
		}
		return tags.Tags
	}).([]string)
}

type dbAction func(collection *mongo.Collection) interface{}

func (t *TagsService) withUsersDb(action dbAction) interface{} {
	collection := t.db.Database("micron").Collection("users")
	return action(collection)
}

func field(username string) bson.D {
	return bson.D{
		{"username", username},
	}
}
