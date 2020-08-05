package user

import (
	"errors"
	log "github.com/sirupsen/logrus"
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

// Register registers the user
func (s *Service) Register(user User) error {
	password := user.Password
	encrypted, err := s.encrypt.Encrypt(password)
	if err != nil {
		log.Errorf("Password encryption failed: %s", err)
		return ErrCouldNotEncryptPassword
	}
	user.Password = encrypted
	s.store.SaveUser(user)
	return nil
}

// Login logs the user in
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
		}
		return EmptyToken, ErrIncorrectPassword
	}
	return EmptyToken, ErrNotFound
}

// Verify verifies the existence of the user by username
func (s *Service) Verify(username string) bool {
	return s.store.FindUser(username) != DoesNotExist
}

// GetUser fetches the user data
func (s *Service) GetUser(username string) (Account, error) {
	user := s.store.FindUser(username)
	if user != DoesNotExist {
		tags := s.GetUserTags(username)
		return Account{
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Name,
			Tags:     tags,
		}, nil
	}
	return EmptyUser, errors.New("user not found for username " + username)
}

// GetUserTags fetches the tags of the user
func (s *Service) GetUserTags(username string) tag.Tags {
	log.Infof("User[%s] tags have been requested", username)
	tagList := []tag.Tag{}
	for _, id := range s.tags.GetUserTags(username) {
		aTag := s.tags.GetTagByID(id)
		if aTag != tag.EmptyTag {
			tagList = append(tagList, aTag)
		}
	}
	return tag.Tags{
		Tags: tagList,
		Size: len(tagList),
	}
}

// AddTagsForUser adds new tags for user
func (s *Service) AddTagsForUser(username string, newTagIds []string) bool {
	log.Infof("Add tags[%s] for user[%s]", newTagIds, username)
	return s.tags.AddTagsForUser(username, newTagIds)
}

// RemoveTagsFromUser removes tags from user
func (s *Service) RemoveTagsFromUser(username string, tagIdsToRemove []string) bool {
	log.Infof("Remove tags[%s] for user[%s]", tagIdsToRemove, username)
	return s.tags.RemoveTagsFromUser(username, tagIdsToRemove)
}

// Logout logs the user out
func (s *Service) Logout(token string) bool {
	return s.jwt.DeleteJwt(token)
}
