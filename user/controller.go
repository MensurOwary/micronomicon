package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micron/commons"
	"micron/model"
	"strings"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

func HandleUserRegistration(c *gin.Context) {
	createdUser := User{}
	_ = c.BindJSON(&createdUser)

	err := createdUser.validateRegister()

	if err != nil {
		c.JSON(400, model.DefaultResponse{
			Message: err.Error(),
		})
	}

	Register(createdUser)

	c.JSON(201, model.DefaultResponse{
		Message: "Created the user",
	})
}

func (user *User) validateRegister() error {
	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is invalid")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is invalid")
	}

	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is invalid")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is invalid")
	}
	return nil
}

func HandleUserAuthorization(c *gin.Context) {
	createdUser := User{}
	_ = c.BindJSON(&createdUser)

	err := createdUser.validateLogin()

	if err != nil {
		c.JSON(400, model.DefaultResponse{
			Message: err.Error(),
		})
		c.Abort()
	}

	token, err := Login(createdUser)
	if err != nil {
		c.JSON(401, model.DefaultResponse{
			Message: err.Error(),
		})
	} else {
		c.JSON(200, AuthResponse{
			AccessToken: *token,
		})
	}
}

func (user *User) validateLogin() error {
	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is invalid")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is invalid")
	}
	return nil
}

func HandleUserByTokenRetrieval(c *gin.Context) {
	commons.WithUsername(c, func(username string) {
		retrievedUser, err := GetUser(username)
		if err != nil {
			c.JSON(404, err)
		} else {
			c.JSON(200, *retrievedUser)
		}
	})
}

func HandleUserTagsRetrieval(c *gin.Context) {
	commons.WithUsername(c, func(username string) {
		tags := GetUserTags(username)
		c.JSON(200, tags)
	})
}

func HandleUserTagsAddition(c *gin.Context) {
	commons.WithUsername(c, func(username string) {
		body := new(map[string][]string)
		_ = c.BindJSON(body)
		tagIds, ok := (*body)["ids"]
		if ok && AddTagsForUser(username, tagIds) {
			c.JSON(200, model.DefaultResponse{
				Message: "Successfully added the tags",
			})
		} else {
			c.JSON(400, model.DefaultResponse{
				Message: "Could not add the new tag(s)",
			})
		}
	})
}

func HandleUserTagsRemoval(c *gin.Context) {
	commons.WithUsername(c, func(username string) {
		body := new(map[string][]string)
		_ = c.BindJSON(body)
		tagIds, ok := (*body)["ids"]
		if ok && RemoveTagsFromUser(username, tagIds) {
			c.JSON(200, model.DefaultResponse{Message: "Successfully removed the tags"})
		} else {
			c.JSON(400, model.DefaultResponse{Message: "Could not remove the tag(s)"})
		}
	})
}
