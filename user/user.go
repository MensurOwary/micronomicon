package user

import (
	"errors"
	"log"
	"micron/tag"
)

// Responses to actions and some constants
var (
	ErrTokenCouldNotBeCreated  = errors.New("could not create the token")     // creation failure
	ErrIncorrectPassword       = errors.New("password is incorrect")          // password match failure
	ErrNotFound                = errors.New("user does not exist")            // user existence failure
	ErrCouldNotEncryptPassword = errors.New("could not encrypt the password") // encryption failure
	EmptyToken                 = ""
	EmptyUser                  = Account{}
)

// Registers the user
func (s *Service) Register(user User) error {
	password := user.Password
	encrypted, err := s.encrypt.Encrypt(password)
	if err != nil {
		log.Println(err)
		return ErrCouldNotEncryptPassword
	}
	user.Password = encrypted
	s.store.SaveUser(user)
	return nil
}

// Logs the user in
func (s *Service) Login(incoming User) (string, error) {
	user := s.store.FindUser(incoming.Username)
	if user != DoesNotExist {
		if s.encrypt.Matches(incoming.Password, user.Password) {
			signingString, err := s.jwt.SignedToken(user.Username)
			if err == nil && s.jwt.SaveJwt(signingString) {
				return signingString, nil
			}
			log.Println(err)
			return EmptyToken, ErrTokenCouldNotBeCreated
		} else {
			return EmptyToken, ErrIncorrectPassword
		}
	}
	return EmptyToken, ErrNotFound
}

// Verifies the existence of the user by username
func (s *Service) Verify(username string) bool {
	return s.store.FindUser(username) != DoesNotExist
}

// Fetches the user data
func (s *Service) GetUser(username string) (Account, error) {
	user := s.store.FindUser(username)
	if user != DoesNotExist {
		tags := s.GetUserTags(username)
		return Account{
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Name,
			Tags: tag.Tags{
				Tags: tags,
				Size: len(tags),
			},
		}, nil
	}
	return EmptyUser, errors.New("user not found for username " + username)
}

// Fetches the tags of the user
func (s *Service) GetUserTags(username string) []tag.Tag {
	log.Printf("User[%s] tags have been requested\n", username)
	var tagList []tag.Tag
	for _, id := range s.tags.GetUserTags(username) {
		aTag := s.tags.GetTagByID(id)
		if aTag != tag.EmptyTag {
			tagList = append(tagList, aTag)
		}
	}
	return tagList
}

// Adds new tags for user
func (s *Service) AddTagsForUser(username string, newTagIds []string) bool {
	log.Printf("Add tags[%s] for user[%s]\n", newTagIds, username)
	return s.tags.AddTagsForUser(username, newTagIds)
}

// Removes tags from user
func (s *Service) RemoveTagsFromUser(username string, tagIdsToRemove []string) bool {
	log.Printf("Remove tags[%s] for user[%s]\n", tagIdsToRemove, username)
	return s.tags.RemoveTagsFromUser(username, tagIdsToRemove)
}

// Logs the user out
func (s *Service) Logout(token string) bool {
	return s.jwt.DeleteJwt(token)
}
