package commons

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type encryptService struct{}

// Deals with encryption concerns of the application
type EncryptService interface {
	Encrypt(input string) (string, error)
	Matches(incoming, existing string) bool
}

func NewEncryptService() EncryptService {
	return &encryptService{}
}

// Possible errors and default values
var (
	ErrCouldNotBeHashed = errors.New("could not hash the user password") // When the password could not be hashed
	EmptyInput          = ""                                             // No hashed password
)

func (e *encryptService) Encrypt(input string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input), Config.EncryptionCost)
	if err != nil {
		log.Printf("Error occurred while hashing the input (%s): %s", input, err)
		return EmptyInput, ErrCouldNotBeHashed
	}
	return string(hashedBytes), nil
}

func (e *encryptService) Matches(incoming, existing string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existing), []byte(incoming)) == nil
}
