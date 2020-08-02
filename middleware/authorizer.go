package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"micron/commons"
	"micron/model"
	"net/http"
)

type UserService interface {
	Verify(username string) bool
}

type JwtService interface {
	DoesJwtExist(jwt string) bool
}

func Authorizer(userService UserService, jwtService JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == "" {
			abort(c, http.StatusBadRequest, "Authorization Header is Missing")
		} else {
			parsedToken, err := parseToken(accessToken)
			if err != nil {
				abort(c, http.StatusBadRequest, err.Error())
			} else {
				if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
					if verifyToken(c, jwtService, parsedToken.Raw) {
						userVerification(c, claims, userService)
					}
				} else {
					abort(c, http.StatusBadRequest, "Token is invalid")
				}
			}
		}
	}
}

func userVerification(c *gin.Context, claims jwt.MapClaims, userService UserService) {
	if username, ok := claims["username"]; ok && userService.Verify(username.(string)) {
		log.Printf("User [%s] successfully requested the resource\n", username.(string))
		c.Set("username", username)
	} else {
		abort(c, http.StatusUnauthorized, "Unauthorized")
	}
}

func verifyToken(c *gin.Context, jwtService JwtService, token string) bool {
	if !jwtService.DoesJwtExist(token) {
		abort(c, http.StatusUnauthorized, "Token expired")
		return false
	}
	return true
}

func abort(c *gin.Context, status int, message string) {
	c.JSON(status, model.DefaultResponse{
		Message: message,
	})
	c.Abort()
}

func parseToken(header string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt error: signing method is wrong")
		}
		return []byte(commons.Config.JwtSecret), nil
	})
	return parsedToken, err
}

func extractToken(c *gin.Context) string {
	return commons.ExtractToken(c.Request.Header)
}
