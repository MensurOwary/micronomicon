package tag

import (
	"github.com/gin-gonic/gin"
)

func HandleTagsRetrieval(c *gin.Context) {
	tagCollection := GetAvailableTags()
	c.JSON(200, tagCollection)
}
