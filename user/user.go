package user

import (
	"errors"
	"log"
	"micron/tag"
)

var (
	TokenCouldNotBeCreated  = errors.New("could not create the token")
	IncorrectPassword       = errors.New("password is incorrect")
	NotFound                = errors.New("user does not exist")
	CouldNotEncryptPassword = errors.New("could not encrypt the password")
	EmptyToken              = ""
)

func (s *Service) Register(user User) error {
	password := user.Password
	encrypted, err := s.encrypt.Encrypt(password)
	if err != nil {
		log.Println(err)
		return CouldNotEncryptPassword
	}
	user.Password = encrypted
	s.store.SaveUser(user)
	return nil
}

func (s *Service) Login(incoming User) (string, error) {
	user := s.store.FindUser(incoming.Username)
	if user != DoesNotExist {
		if s.encrypt.Matches(incoming.Password, user.Password) {
			signingString, err := s.jwt.SignedToken(user.Username)
			if err == nil && s.jwt.SaveJwt(signingString) {
				return signingString, nil
			}
			log.Println(err)
			return EmptyToken, TokenCouldNotBeCreated
		} else {
			return EmptyToken, IncorrectPassword
		}
	}
	return EmptyToken, NotFound
}

func (s *Service) Verify(username string) bool {
	return s.store.FindUser(username) != DoesNotExist
}

var EmptyUser = Account{}

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

func (s *Service) GetUserTags(username string) []tag.Tag {
	log.Printf("User[%s] tags have been requested\n", username)
	var tagList []tag.Tag
	for _, id := range s.tags.GetUserTags(username) {
		aTag := s.tags.GetTagById(id)
		if aTag != tag.EmptyTag {
			tagList = append(tagList, aTag)
		}
	}
	return tagList
}

func (s *Service) AddTagsForUser(username string, newTagIds []string) bool {
	log.Printf("Add tags[%s] for user[%s]\n", newTagIds, username)
	return s.tags.AddTagsForUser(username, newTagIds)
}

func (s *Service) RemoveTagsFromUser(username string, tagIdsToRemove []string) bool {
	log.Printf("Remove tags[%s] for user[%s]\n", tagIdsToRemove, username)
	return s.tags.RemoveTagsFromUser(username, tagIdsToRemove)
}

func (s *Service) DeleteToken(token string) bool {
	return s.jwt.DeleteJwt(token)
}
