package commons

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"micron/model"
	"net/http"
	"os"
	"strings"
)

type Callable func(username string)

func WithUsername(c *gin.Context, callable Callable) {
	username, ok := c.Keys["username"].(string)
	if ok {
		callable(username)
	} else {
		c.JSON(400, model.DefaultResponse{
			Message: "username was not found",
		})
	}
}

func ExtractToken(header http.Header) string {
	headerToken := header.Get("Authorization")
	if strings.Index(headerToken, "Bearer ") == 0 {
		headerToken = strings.ReplaceAll(headerToken, "Bearer ", "")
		headerToken = strings.TrimSpace(headerToken)
		return headerToken
	}
	return ""
}

func Union(a []string, b []string) []string {
	memo := make(map[string]int)
	c := append(a, b...)
	for _, e := range c {
		memo[e] = 1
	}
	z := make([]string, 0)
	for k, _ := range memo {
		z = append(z, k)
	}
	return z
}

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

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Missing the value for the environment variable [%s]", key))
	}
	return value
}
