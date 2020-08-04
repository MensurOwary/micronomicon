package tag

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type tagsService interface {
	GetAvailableTags() Tags
}

// HandleTagsRetrieval deals with the retrieval of tags
func HandleTagsRetrieval(c *gin.Context, service tagsService) {
	tagCollection := service.GetAvailableTags()
	c.JSON(http.StatusOK, tagCollection)
}
