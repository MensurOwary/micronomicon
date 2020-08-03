package user

import (
	"github.com/gin-gonic/gin"
	"micron/commons"
	"micron/model"
	"micron/tag"
	"net/http"
)

// Responses to actions
var (
	TagsAddedSuccessfully   = model.Response("Successfully added the tags")
	NewTagsCouldNotBeAdded  = model.Response("Could not add the new tag(s)")
	TagsRemovedSuccessfully = model.Response("Successfully removed the tags")
	TagsCouldNotBeRemoved   = model.Response("Could not remove the tag(s)")
)

type usersInteractionService interface {
	GetUser(username string) (Account, error)
}

type tagsInteractionService interface {
	GetUserTags(username string) []tag.Tag
	AddTagsForUser(username string, newTagIds []string) bool
	RemoveTagsFromUser(username string, tagIdsToRemove []string) bool
}

// Deals with getting the user related data based on the bearer token
func HandleUserByTokenRetrieval(c *gin.Context, service usersInteractionService) {
	commons.WithUsername(c, func(username string) {
		retrievedUser, err := service.GetUser(username)
		if err != nil {
			c.JSON(http.StatusNotFound, model.Response(err.Error()))
		} else {
			c.JSON(http.StatusOK, retrievedUser)
		}
	})
}

// Deals with getting the user tags
func HandleUserTagsRetrieval(c *gin.Context, service tagsInteractionService) {
	commons.WithUsername(c, func(username string) {
		tags := service.GetUserTags(username)
		c.JSON(http.StatusOK, tags)
	})
}

// Deals with adding new tags for the user
func HandleUserTagsAddition(c *gin.Context, service tagsInteractionService) {
	commons.WithUsername(c, func(username string) {
		handleTagAddRemove(
			c,
			func(tagIds []string) bool {
				return service.AddTagsForUser(username, tagIds)
			},
			TagsAddedSuccessfully, NewTagsCouldNotBeAdded,
		)
	})
}

// Deals with removing some tags from the user
func HandleUserTagsRemoval(c *gin.Context, service tagsInteractionService) {
	commons.WithUsername(c, func(username string) {
		handleTagAddRemove(
			c,
			func(tagIds []string) bool {
				return service.RemoveTagsFromUser(username, tagIds)
			},
			TagsRemovedSuccessfully, TagsCouldNotBeRemoved,
		)
	})
}

type addOrRemoveTagsAction func(tagIds []string) bool

func handleTagAddRemove(c *gin.Context, action addOrRemoveTagsAction, successObj interface{}, failObj interface{}) {
	body := new(map[string][]string)
	_ = c.BindJSON(body)
	tagIds, ok := (*body)["ids"]
	if ok && action(tagIds) {
		c.JSON(http.StatusOK, successObj)
	} else {
		c.JSON(http.StatusBadRequest, failObj)
	}
}
