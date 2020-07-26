package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"micron/commons"
	"micron/tag"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Account struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Tags     tag.Tags `json:"tags"`
}

func Register(user User) {
	password := user.Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal("Could not hash the user password")
	}
	user.Password = string(hashedBytes)
	DoSaveUser(&user)
}

func Login(incoming User) (*string, error) {
	user := DoFindUser(incoming.Username)
	if user != nil {
		hashedPassword := user.Password
		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(incoming.Password)) == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": user.Username,
			})
			signingString, err := token.SignedString([]byte(commons.Config.JwtSecret))
			if err == nil {
				return &signingString, nil
			}
			log.Println(err.Error())
			return nil, errors.New("could not create the token")
		} else {
			return nil, errors.New("password is incorrect")
		}
	}
	return nil, errors.New("user does not exist")
}

func Verify(username string) bool {
	return DoFindUser(username) != nil
}

func GetUser(username string) (*Account, error) {
	user := DoFindUser(username)
	if user != nil {
		if user.Username == username {
			tags := GetUserTags(username)
			return &Account{
				Username: user.Username,
				Email:    user.Email,
				Name:     user.Name,
				Tags: tag.Tags{
					Tags: tags,
					Size: len(tags),
				},
			}, nil
		}
	}
	return nil, errors.New("user not found for username " + username)
}

func GetUserTags(username string) []tag.Tag {
	log.Printf("User[%s] tags have been requested\n", username)
	var tagList []tag.Tag
	for _, id := range DoGetUserTags(username) {
		aTag := tag.GetTagById(id)
		if aTag != nil {
			tagList = append(tagList, *aTag)
		}
	}
	return tagList
}

func AddTagsForUser(username string, newTagIds []string) bool {
	log.Printf("Add tags[%s] for user[%s]\n", newTagIds, username)
	return DoAddTagsForUser(username, newTagIds)
}

func RemoveTagsFromUser(username string, tagIdsToRemove []string) bool {
	log.Printf("Remove tags[%s] for user[%s]\n", tagIdsToRemove, username)
	return DoRemoveTagsFromUser(username, tagIdsToRemove)
}
