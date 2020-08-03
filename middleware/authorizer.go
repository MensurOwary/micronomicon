package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"micron/commons"
	"micron/model"
	"net/http"
)

type userService interface {
	Verify(username string) bool
}

type jwtService interface {
	ParseJwt(rawJwt string) (*commons.ParsedJwtResult, error)
}

var EmptyToken = ""

// middleware that deals with bearer token authorizations
func Authorizer(userService userService, jwtService jwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == EmptyToken {
			abort(c, http.StatusBadRequest, "Authorization header is missing or empty")
		} else {
			parseJwt, err := jwtService.ParseJwt(accessToken)

			if err != nil {
				abort(c, http.StatusBadRequest, err.Error())
			} else if err := userVerification(c, parseJwt, userService); err != nil {
				abort(c, http.StatusUnauthorized, "Unauthorized")
			}
		}
	}
}

func userVerification(c *gin.Context, parsed *commons.ParsedJwtResult, userService userService) error {
	if userService.Verify(parsed.Username) == false {
		return errors.New("unauthorized")
	}
	log.Printf("User [%s] successfully requested the resource\n", parsed.Username)
	c.Set("username", parsed.Username)
	return nil
}

func abort(c *gin.Context, status int, message string) {
	c.JSON(status, model.Response(message))
	c.Abort()
}

func extractToken(c *gin.Context) string {
	return commons.ExtractToken(c.Request.Header)
}
