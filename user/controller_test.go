package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"micron/commons"
	"micron/model"
	"micron/tag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleUserByTokenRetrieval(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("When user is found for the token", func(t *testing.T) {
		service := &mockUsersInteractionService{}

		rec, c := recorderAndContext()

		HandleUserByTokenRetrieval(c, service)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, commons.ToJSON(Account{Username: "jane"}), rec.Body.String())
	})

	t.Run("When user is not found for the token", func(t *testing.T) {
		service := &mockUsersInteractionService{
			err: ErrNotFound,
		}

		rec, c := recorderAndContext()

		HandleUserByTokenRetrieval(c, service)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, commons.ToJSON(model.Response(ErrNotFound.Error())), rec.Body.String())
	})
}

func TestHandleUserTagsRetrieval(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Returns the tags of the user", func(t *testing.T) {
		service := &mockTagsInteractionService{}

		rec, c := recorderAndContext()

		HandleUserTagsRetrieval(c, service)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, commons.ToJSON([]tag.Tag{
			{Name: "react"},
			{Name: "ruby"},
		}), rec.Body.String())
	})

}

func TestHandleUserTagsAddition(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Tags added successfully", func(t *testing.T) {
		service := &mockTagsInteractionService{addTagsForUser: true}

		rec, c := recorderAndContextAndBody(true)

		HandleUserTagsAddition(c, service)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, commons.ToJSON(TagsAddedSuccessfully), rec.Body.String())
	})

	t.Run("Tags could not be added", func(t *testing.T) {
		service := &mockTagsInteractionService{addTagsForUser: false}

		rec, c := recorderAndContextAndBody(true)

		HandleUserTagsAddition(c, service)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, commons.ToJSON(NewTagsCouldNotBeAdded), rec.Body.String())
	})
}

func TestHandleUserTagsRemoval(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Tags removed successfully", func(t *testing.T) {
		service := &mockTagsInteractionService{removeTagsFromUser: true}

		rec, c := recorderAndContextAndBody(true)

		HandleUserTagsRemoval(c, service)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, commons.ToJSON(TagsRemovedSuccessfully), rec.Body.String())
	})

	t.Run("Tags could not be removed", func(t *testing.T) {
		service := &mockTagsInteractionService{removeTagsFromUser: false}

		rec, c := recorderAndContextAndBody(true)

		HandleUserTagsRemoval(c, service)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, commons.ToJSON(TagsCouldNotBeRemoved), rec.Body.String())
	})

}

func recorderAndContext() (*httptest.ResponseRecorder, *gin.Context) {
	return recorderAndContextAndBody(false)
}

func recorderAndContextAndBody(shouldHaveBody bool) (*httptest.ResponseRecorder, *gin.Context) {
	var body io.Reader
	if shouldHaveBody {
		body = strings.NewReader(`
				{
					"ids":["react", "go"]
				}`,
		)
	} else {
		body = nil
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/something", body)
	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Set("username", "jane")
	return rec, c
}

type mockUsersInteractionService struct {
	err error
}

func (m *mockUsersInteractionService) GetUser(username string) (Account, error) {
	return Account{Username: username}, m.err
}

type mockTagsInteractionService struct {
	addTagsForUser, removeTagsFromUser bool
}

func (m *mockTagsInteractionService) GetUserTags(_ string) []tag.Tag {
	return []tag.Tag{
		{Name: "react"},
		{Name: "ruby"},
	}
}

func (m *mockTagsInteractionService) AddTagsForUser(_ string, _ []string) bool {
	return m.addTagsForUser
}

func (m *mockTagsInteractionService) RemoveTagsFromUser(_ string, _ []string) bool {
	return m.removeTagsFromUser
}
