package tag

import (
	"github.com/gin-gonic/gin"
)

type tagsService interface {
	GetAvailableTags() Tags
}

func HandleTagsRetrieval(c *gin.Context, service tagsService) {
	tagCollection := service.GetAvailableTags()
	c.JSON(200, tagCollection)
}
