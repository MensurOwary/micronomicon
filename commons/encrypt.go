package commons

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type encryptService struct{}

type EncryptService interface {
	Encrypt(input string) (string, error)
	Matches(incoming, existing string) bool
}

func NewEncryptService() EncryptService {
	return &encryptService{}
}

var (
	CouldNotBeHashed = errors.New("could not hash the user password")
	EmptyInput       = ""
)

func (e *encryptService) Encrypt(input string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input), 14)
	if err != nil {
		log.Printf("Error occurred while hashing the input (%s): %s", input, err)
		return EmptyInput, CouldNotBeHashed
	}
	return string(hashedBytes), nil
}

func (e *encryptService) Matches(incoming, existing string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existing), []byte(incoming)) == nil
}
