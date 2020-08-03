package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"micron/commons"
	"micron/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("When Authorization header is missing/empty", func(t *testing.T) {
		// given
		userService := &mockService{}
		jwtService := &mockService{}
		recorder, c := recorderAndContext(false)

		// when
		Authorizer(userService, jwtService)(c)
		// then
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, commons.ToJSON(model.Response("Authorization header is missing or empty")), recorder.Body.String())
		checkUsernameKeyStatus(t, c, false)
	})

	t.Run("When parsing token fails", func(t *testing.T) {
		// given
		userService := &mockService{}
		jwtService := &mockService{
			parseErr: errors.New("parsing failed"),
		}

		recorder, c := recorderAndContext(true)

		// when
		Authorizer(userService, jwtService)(c)
		// then
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, commons.ToJSON(model.Response("parsing failed")), recorder.Body.String())
		checkUsernameKeyStatus(t, c, false)
	})

	t.Run("When user verification fails", func(t *testing.T) {
		// given
		userService := &mockService{}
		jwtService := &mockService{}
		recorder, c := recorderAndContext(true)

		// when
		Authorizer(userService, jwtService)(c)
		// then
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Equal(t, commons.ToJSON(model.Response("Unauthorized")), recorder.Body.String())
		checkUsernameKeyStatus(t, c, false)
	})

	t.Run("When everything is successful", func(t *testing.T) {
		// given
		userService := &mockService{verifyResult: true}
		jwtService := &mockService{}
		recorder, c := recorderAndContext(true)

		// when
		Authorizer(userService, jwtService)(c)
		// then
		assert.Equal(t, http.StatusOK, recorder.Code)
		checkUsernameKeyStatus(t, c, true)
	})

}

func recorderAndContext(header bool) (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("POST", "http://localhost/logout", nil)
	if header {
		request.Header.Add("Authorization", "Bearer random_token")
	}
	c.Request = request
	return recorder, c
}

func checkUsernameKeyStatus(t *testing.T, c *gin.Context, shouldBeSet bool) {
	usernameField := c.Keys["username"]
	if shouldBeSet {
		assert.NotNil(t, usernameField, "username field should be set after successful auth")
	} else {
		assert.Nil(t, usernameField, "username field should not be set")
	}
}

type mockService struct {
	verifyResult bool
	parseErr     error
}

func (m *mockService) Verify(_ string) bool {
	return m.verifyResult
}
func (m *mockService) ParseJwt(_ string) (*commons.ParsedJwtResult, error) {
	return &commons.ParsedJwtResult{
		Username: "john",
	}, m.parseErr
}
