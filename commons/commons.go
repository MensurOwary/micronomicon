package commons

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"micron/model"
	"net/http"
	"os"
	"strings"
)

type callable func(username string)

// Extracts the username from the context and
// passes it to the provided callable function
func WithUsername(c *gin.Context, callable callable) {
	username, ok := c.Keys["username"].(string)
	if ok {
		callable(username)
	} else {
		c.JSON(http.StatusBadRequest, model.Response("username was not found"))
	}
}

// Extracts the Bearer token from the Authorization header
func ExtractToken(header http.Header) string {
	headerToken := header.Get("Authorization")
	if strings.Index(headerToken, "Bearer ") == 0 && len(headerToken) > 7 {
		headerToken = strings.ReplaceAll(headerToken, "Bearer ", "")
		headerToken = strings.TrimSpace(headerToken)
		return headerToken
	}
	return ""
}

// Union of two string arrays
func Union(a []string, b []string) []string {
	memo := make(map[string]int)
	c := append(a, b...)
	for _, e := range c {
		memo[e] = 1
	}
	z := make([]string, 0)
	for k := range memo {
		z = append(z, k)
	}
	return z
}

// Difference of two string arrays
func Difference(a []string, b []string) []string {
	memo := make(map[string]int)
	for _, e := range a {
		memo[e] = 1
	}
	for _, e := range b {
		memo[e] = 0
	}
	z := make([]string, 0)
	for k, v := range memo {
		if v == 1 {
			z = append(z, k)
		}
	}
	return z
}

// Gets the environment variable
// panics when the value is not present or is empty
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Missing the value for the environment variable [%s]", key))
	}
	return value
}
