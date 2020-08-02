package micron

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"micron/commons"
	"micron/model"
	"micron/user"
)

func HandleMicronRetrieval(c *gin.Context, s Service, user *user.Service) {
	commons.WithUsername(c, func(username string) {
		micron := getRandomMicron(username, s, user)
		if micron != EmptyMicron {
			c.JSON(200, micron)
		} else {
			c.JSON(404, model.Response("Could not find a micron"))
		}
	})
}

func getRandomMicron(username string, s Service, user *user.Service) Micron {
	tags := user.GetUserTags(username)
	if len(tags) > 0 {
		chosen := rand.Intn(len(tags))
		return s.GetARandomMicronForTag(tags[chosen])
	} else {
		return EmptyMicron
	}
}
