package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"micron/commons"
	"net/http"
)

type userService interface {
	Verify(username string) bool
}

type jwtService interface {
	ParseJwt(rawJwt string) (*commons.ParsedJwtResult, error)
}

var EmptyToken = ""

// Authorizer is a middleware that deals with bearer token authorizations
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
		log.Warnf("User[%s] does not exist", parsed.Username)
		return errors.New("unauthorized")
	}
	log.Infof("User[%s] successfully requested the resource : %s", parsed.Username, c.Request.RequestURI)
	c.Set("username", parsed.Username)
	return nil
}

func abort(c *gin.Context, status int, message string) {
	log.Errorf("Operation resulted in failure...aborting. Cause : %s", message)
	c.JSON(status, commons.Response(message))
	c.Abort()
}

func extractToken(c *gin.Context) string {
	return commons.ExtractToken(c.Request.Header)
}
