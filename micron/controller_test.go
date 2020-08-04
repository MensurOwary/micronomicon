package micron

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"micron/commons"
	"micron/tag"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleMicronRetrieval(t *testing.T) {
	t.Run("Could not find micron", func(t *testing.T) {
		service := &mockService{}
		recorder, c := recorderAndContext(true)

		HandleMicronRetrieval(c, service, service)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, commons.ToJSON(commons.Response("Could not find a micron")), recorder.Body.String())
	})

	t.Run("Successfully gets a micron", func(t *testing.T) {
		service := &mockService{
			tags: []tag.Tag{
				{Name: "stuff"},
			},
		}
		recorder, c := recorderAndContext(true)

		HandleMicronRetrieval(c, service, service)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, commons.ToJSON(Micron{
			URL:   "www.hello.com",
			Title: "how to handle stuff",
		}), recorder.Body.String())
	})

	t.Run("When 'username' key is not present", func(t *testing.T) {
		service := &mockService{}
		recorder, c := recorderAndContext(false)

		HandleMicronRetrieval(c, service, service)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, commons.ToJSON(commons.Response("username was not found")), recorder.Body.String())
	})
}

func recorderAndContext(key bool) (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("POST", "http://localhost/users/me/microns", nil)
	c.Request = request
	if key {
		c.Set("username", "jane")
	}
	return recorder, c
}

type mockService struct {
	tags []tag.Tag
}

func (m *mockService) GetARandomMicronForTag(tag tag.Tag) Micron {
	if len(m.tags) > 0 {
		return Micron{
			id:    "demo-id",
			URL:   "www.hello.com",
			Title: "how to handle stuff",
			tag:   tag,
		}
	}
	return EmptyMicron
}
func (m *mockService) GetUserTags(_ string) []tag.Tag {
	return m.tags
}
