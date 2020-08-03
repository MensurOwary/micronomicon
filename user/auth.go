package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"micron/commons"
	"micron/model"
	"strings"
)

// Responses
var (
	ErrInvalidUsername = errors.New("username is invalid")
	ErrInvalidEmail    = errors.New("email is invalid")
	ErrInvalidName     = errors.New("name is invalid")
	ErrInvalidPassword = errors.New("password is invalid")
)

// AuthService represents a means that deals with authZ/N related actions
type AuthService interface {
	Register(user User) error
	Login(incoming User) (string, error)
	Logout(token string) bool
}

// AuthResponse represents authorization response
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// HandleUserRegistration deals with registering a user
func HandleUserRegistration(c *gin.Context, service AuthService) {
	createdUser := User{}
	_ = c.BindJSON(&createdUser)

	err := createdUser.validateRegister()

	if err != nil {
		c.JSON(400, model.Response(err.Error()))
		c.Abort()
	}

	if service.Register(createdUser) == ErrCouldNotEncryptPassword {
		c.JSON(422, model.Response("Could not process the request"))
	} else {
		c.JSON(201, model.Response("Created the user"))
	}
}

func (user *User) validateRegister() error {
	if strings.TrimSpace(user.Username) == "" {
		return ErrInvalidUsername
	}

	if strings.TrimSpace(user.Email) == "" {
		return ErrInvalidEmail
	}

	if strings.TrimSpace(user.Name) == "" {
		return ErrInvalidName
	}

	if strings.TrimSpace(user.Password) == "" {
		return ErrInvalidPassword
	}
	return nil
}

// HandleUserAuthorization deals with authN/Z of the user
func HandleUserAuthorization(c *gin.Context, service AuthService) {
	createdUser := User{}
	_ = c.BindJSON(&createdUser)

	err := createdUser.validateLogin()

	if err != nil {
		c.JSON(400, model.Response(err.Error()))
		c.Abort()
	}

	token, err := service.Login(createdUser)
	if err != nil {
		c.JSON(401, model.Response(err.Error()))
	} else {
		c.JSON(200, AuthResponse{
			AccessToken: token,
		})
	}
}

func (user *User) validateLogin() error {
	if strings.TrimSpace(user.Username) == "" {
		return ErrInvalidUsername
	}

	if strings.TrimSpace(user.Password) == "" {
		return ErrInvalidPassword
	}
	return nil
}

// HandleUserLogout deals with user logout
func HandleUserLogout(c *gin.Context, service AuthService) {
	token := commons.ExtractToken(c.Request.Header)
	if !service.Logout(token) {
		log.Println("Could not delete token successfully")
	}
	c.JSON(200, model.Response("Successfully logged out"))
}
