package commons

import (
	"github.com/gin-gonic/gin"
	"micron/model"
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
