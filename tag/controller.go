package tag

import (
	"github.com/gin-gonic/gin"
)

func HandleTagsRetrieval(c *gin.Context, service *Repository) {
	tagCollection := service.GetAvailableTags()
	c.JSON(200, tagCollection)
}
