package commons

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type encryptService struct{}

// EncryptService deals with encryption concerns of the application
type EncryptService interface {
	Encrypt(input string) (string, error)
	Matches(incoming, existing string) bool
}

// NewEncryptService initializes EncryptService
func NewEncryptService() EncryptService {
	return &encryptService{}
}

// Possible errors and default values
var (
	ErrCouldNotBeHashed = errors.New("could not hash the user password") // When the password could not be hashed
	EmptyInput          = ""                                             // No hashed password
)

// Encrypt encrypts the given input string
func (e *encryptService) Encrypt(input string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input), Config.EncryptionCost)
	if err != nil {
		log.Errorf("Error occurred while hashing the input (%s): %s", input, err)
		return EmptyInput, ErrCouldNotBeHashed
	}
	return string(hashedBytes), nil
}

// Matches checks if two given strings match where the former being the actual password,
// and the latter being the hashed one
func (e *encryptService) Matches(incoming, existing string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existing), []byte(incoming)) == nil
}
