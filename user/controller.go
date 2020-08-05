package user

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"micron/commons"
	"micron/tag"
	"net/http"
)

// Responses to actions
var (
	TagsAddedSuccessfully   = commons.Response("Successfully added the tags")
	NewTagsCouldNotBeAdded  = commons.Response("Could not add the new tag(s)")
	TagsRemovedSuccessfully = commons.Response("Successfully removed the tags")
	TagsCouldNotBeRemoved   = commons.Response("Could not remove the tag(s)")
)

type usersInteractionService interface {
	GetUser(username string) (Account, error)
}

type tagsInteractionService interface {
	GetUserTags(username string) tag.Tags
	AddTagsForUser(username string, newTagIds []string) bool
	RemoveTagsFromUser(username string, tagIdsToRemove []string) bool
}

// HandleUserByTokenRetrieval deals with getting the user related data based on the bearer token
func HandleUserByTokenRetrieval(c *gin.Context, service usersInteractionService) {
	commons.WithUsername(c, func(username string) {
		retrievedUser, err := service.GetUser(username)
		if err != nil {
			log.Errorf("Could not retrieve the user : %s", err.Error())
			c.JSON(http.StatusNotFound, commons.Response(err.Error()))
		} else {
			log.Infof("Successfully retrieved the user data for [%s]", username)
			c.JSON(http.StatusOK, retrievedUser)
		}
	})
}

// HandleUserTagsRetrieval deals with getting the user tags
func HandleUserTagsRetrieval(c *gin.Context, service tagsInteractionService) {
	commons.WithUsername(c, func(username string) {
		tags := service.GetUserTags(username)
		c.JSON(http.StatusOK, tags)
	})
}

// HandleUserTagsAddition deals with adding new tags for the user
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

// HandleUserTagsRemoval deals with removing some tags from the user
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
		log.Info("Tag adding/removing succeeded")
		c.JSON(http.StatusOK, successObj)
	} else {
		log.Error("Tag adding/removing failed")
		c.JSON(http.StatusBadRequest, failObj)
	}
}
