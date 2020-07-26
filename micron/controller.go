package micron

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"micron/commons"
	"micron/model"
	"micron/user"
)

func HandleMicronRetrieval(c *gin.Context) {
	commons.WithUsername(c, func(username string) {
		micronRef := getRandomMicron(username)
		if micronRef != nil {
			c.JSON(200, *micronRef)
		} else {
			c.JSON(404, model.DefaultResponse{
				Message: "Could not find a micron",
			})
		}
	})
}

func getRandomMicron(username string) *Micron {
	tags := user.GetUserTags(username)
	if len(tags) > 0 {
		chosen := rand.Intn(len(tags))
		micronRef := GetARandomMicronForTag(tags[chosen])
		return micronRef
	} else {
		return nil
	}
}
