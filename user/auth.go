package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"micron/commons"
	"net/http"
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
		log.Errorf("User payload validation failed : %s", err.Error())
		c.JSON(http.StatusBadRequest, commons.Response(err.Error()))
		c.Abort()
	}

	if service.Register(createdUser) == ErrCouldNotEncryptPassword {
		log.Error("Could not register the user : encrypting the password failed")
		c.JSON(http.StatusUnprocessableEntity, commons.Response("Could not process the request"))
	} else {
		log.Info("Successfully registered the user")
		c.JSON(http.StatusCreated, commons.Response("Created the user"))
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
		log.Errorf("User payload validation failed : %s", err.Error())
		c.JSON(http.StatusBadRequest, commons.Response(err.Error()))
		c.Abort()
	}

	token, err := service.Login(createdUser)
	if err != nil {
		log.Errorf("User[%s] authorization failed : %s", createdUser.Username, err.Error())
		c.JSON(http.StatusUnauthorized, commons.Response(err.Error()))
	} else {
		log.Infof("User[%s] successfully authorized", createdUser.Username)
		c.JSON(http.StatusOK, AuthResponse{
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
		log.Errorf("Could not delete token successfully")
	}
	c.JSON(http.StatusOK, commons.Response("Successfully logged out"))
}
