package tag

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"micron/commons"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleTagsRetrieval(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mockService)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	HandleTagsRetrieval(c, service)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, commons.ToJSON(Tags{
		Tags: []Tag{
			{Name: "react"},
			{Name: "ruby"},
		},
		Size: 2,
	}), recorder.Body.String())
}

type mockService struct{}

func (m *mockService) GetAvailableTags() Tags {
	return Tags{
		Tags: []Tag{
			{Name: "react"},
			{Name: "ruby"},
		},
		Size: 2,
	}
}