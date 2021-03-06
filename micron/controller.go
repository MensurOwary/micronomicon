package micron

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"micron/commons"
	"micron/tag"
	"net/http"
)

type service interface {
	GetARandomMicronForTag(tag tag.Tag) Micron
}

type userService interface {
	GetUserTags(username string) tag.Tags
}

// HandleMicronRetrieval handles the retrieval of a micron
func HandleMicronRetrieval(c *gin.Context, s service, user userService) {
	commons.WithUsername(c, func(username string) {
		micron := getRandomMicron(username, s, user)
		if micron == EmptyMicron {
			log.Warnf("Micron could not be found for user[%s]", username)
			c.JSON(http.StatusNotFound, commons.Response("Could not find a micron"))
		} else {
			c.JSON(http.StatusOK, micron)
		}
	})
}

func getRandomMicron(username string, s service, user userService) Micron {
	tags := user.GetUserTags(username)
	if tags.Size > 0 {
		chosen := rand.Intn(tags.Size)
		return s.GetARandomMicronForTag(tags.Tags[chosen])
	}
	return EmptyMicron
}
