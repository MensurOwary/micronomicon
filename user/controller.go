package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"micron/commons"
	"micron/model"
	"strings"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

func HandleUserRegistration(c *gin.Context, service *Service) {
	createdUser := User{}
	_ = c.BindJSON(&createdUser)

	err := createdUser.validateRegister()

	if err != nil {
		c.JSON(400, model.DefaultResponse{
			Message: err.Error(),
		})
	}

	service.Register(createdUser)

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

func HandleUserAuthorization(c *gin.Context, service *Service) {
	createdUser := User{}
	_ = c.BindJSON(&createdUser)

	err := createdUser.validateLogin()

	if err != nil {
		c.JSON(400, model.DefaultResponse{
			Message: err.Error(),
		})
		c.Abort()
	}

	token, err := service.Login(createdUser)
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

var (
	InvalidUsername = errors.New("username is invalid")
	InvalidPassword = errors.New("password is invalid")
)

func (user *User) validateLogin() error {
	if strings.TrimSpace(user.Username) == "" {
		return InvalidUsername
	}

	if strings.TrimSpace(user.Password) == "" {
		return InvalidPassword
	}
	return nil
}

func HandleUserLogout(c *gin.Context, service *Service) {
	token := commons.ExtractToken(c.Request.Header)
	if !service.DeleteToken(token) {
		log.Println("Could not delete token successfully")
	}
	c.JSON(200, model.DefaultResponse{
		Message: "Successfully logged out",
	})
}

func HandleUserByTokenRetrieval(c *gin.Context, service *Service) {
	commons.WithUsername(c, func(username string) {
		retrievedUser, err := service.GetUser(username)
		if err != nil {
			c.JSON(404, err)
		} else {
			c.JSON(200, *retrievedUser)
		}
	})
}

func HandleUserTagsRetrieval(c *gin.Context, service *Service) {
	commons.WithUsername(c, func(username string) {
		tags := service.GetUserTags(username)
		c.JSON(200, tags)
	})
}

var (
	TagsAddedSuccessfully   = model.DefaultResponse{Message: "Successfully added the tags"}
	NewTagsCouldNotBeAdded  = model.DefaultResponse{Message: "Could not add the new tag(s)"}
	TagsRemovedSuccessfully = model.DefaultResponse{Message: "Successfully removed the tags"}
	TagsCouldNotBeRemoved   = model.DefaultResponse{Message: "Could not remove the tag(s)"}
)

func HandleUserTagsAddition(c *gin.Context, service *Service) {
	commons.WithUsername(c, func(username string) {
		body := new(map[string][]string)
		_ = c.BindJSON(body)
		tagIds, ok := (*body)["ids"]
		if ok && service.AddTagsForUser(username, tagIds) {
			c.JSON(200, TagsAddedSuccessfully)
		} else {
			c.JSON(400, NewTagsCouldNotBeAdded)
		}
	})
}

func HandleUserTagsRemoval(c *gin.Context, service *Service) {
	commons.WithUsername(c, func(username string) {
		body := new(map[string][]string)
		_ = c.BindJSON(body)
		tagIds, ok := (*body)["ids"]
		if ok && service.RemoveTagsFromUser(username, tagIds) {
			c.JSON(200, TagsRemovedSuccessfully)
		} else {
			c.JSON(400, TagsCouldNotBeRemoved)
		}
	})
}
