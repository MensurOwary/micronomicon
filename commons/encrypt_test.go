package commons

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncryptService_Encrypt(t *testing.T) {
	t.Run("Hashing and matching work", func(t *testing.T) {
		Config.EncryptionCost = 14
		service := NewEncryptService()

		result, err := service.Encrypt("hello")

		assert.Nil(t, err)
		assert.True(t, service.Matches("hello", result))
	})

	t.Run("When error", func(t *testing.T) {
		Config.EncryptionCost = bcrypt.MaxCost + 1 // causes error
		service := NewEncryptService()

		result, err := service.Encrypt("hello")

		assert.Equal(t, ErrCouldNotBeHashed, err)
		assert.Equal(t, EmptyInput, result)
	})
}
