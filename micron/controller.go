package micron

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"micron/commons"
	"micron/model"
	"micron/tag"
	"net/http"
)

type service interface {
	GetARandomMicronForTag(tag tag.Tag) Micron
}

type userService interface {
	GetUserTags(username string) []tag.Tag
}

// HandleMicronRetrieval handles the retrieval of a micron
func HandleMicronRetrieval(c *gin.Context, s service, user userService) {
	commons.WithUsername(c, func(username string) {
		micron := getRandomMicron(username, s, user)
		if micron == EmptyMicron {
			c.JSON(http.StatusNotFound, model.Response("Could not find a micron"))
		} else {
			c.JSON(http.StatusOK, micron)
		}
	})
}

func getRandomMicron(username string, s service, user userService) Micron {
	tags := user.GetUserTags(username)
	if len(tags) > 0 {
		chosen := rand.Intn(len(tags))
		return s.GetARandomMicronForTag(tags[chosen])
	}
	return EmptyMicron
}
