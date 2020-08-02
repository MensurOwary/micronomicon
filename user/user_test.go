package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"micron/tag"
	"testing"
)

func TestService_Register(t *testing.T) {
	t.Run("Registers user successfully", func(t *testing.T) {
		store, service := getService(mockStore{}, mockEncrypt{}, mockJwt{})
		user := User{
			Username: "username",
			Email:    "email@email.com",
			Name:     "Jane Doe",
			Password: "very strong password",
		}

		assert.Nil(
			t,
			service.Register(user),
			"The result should be nil because of the success",
		)

		type value struct {
			username, name, email, password string
		}

		assert.Equal(t,
			value{
				username: "username",
				email:    "email@email.com",
				name:     "Jane Doe",
				password: "encrypted-password",
			},
			value{
				username: store.user.Username,
				email:    store.user.Email,
				name:     store.user.Name,
				password: store.user.Password,
			})
	})

	t.Run("Password is nil", func(t *testing.T) {
		_, service := getService(mockStore{}, mockEncrypt{
			shouldError: true,
		}, mockJwt{})
		user := User{
			Username: "username",
			Email:    "email@email.com",
			Name:     "Jane Doe",
			Password: "",
		}
		err := service.Register(user)

		assert.NotNil(t, err, "The result should error")
		assert.Equal(t, CouldNotEncryptPassword, err, "Should not be able to encrypt the password")
	})
}

func TestService_Login(t *testing.T) {

	tt := []struct {
		testName      string
		store         mockStore
		encrypt       mockEncrypt
		jwt           mockJwt
		tokenExpected string
		errorExpected error
	}{
		{
			testName: "User is not found",
			store: mockStore{
				user: DoesNotExist,
			},
			encrypt:       mockEncrypt{},
			jwt:           mockJwt{},
			tokenExpected: EmptyToken,
			errorExpected: NotFound,
		},
		{
			testName: "Password does not match",
			store: mockStore{
				user: User{Username: "demo"},
			},
			encrypt: mockEncrypt{
				usernameMatches: false,
			},
			jwt:           mockJwt{},
			tokenExpected: EmptyToken,
			errorExpected: IncorrectPassword,
		},
		{
			testName: "Token could not be created",
			store: mockStore{
				user: User{Username: "demo"},
			},
			encrypt: mockEncrypt{
				usernameMatches: true,
			},
			jwt: mockJwt{
				token: "",
				err:   errors.New("some error"),
			},
			tokenExpected: EmptyToken,
			errorExpected: TokenCouldNotBeCreated,
		},
		{
			testName: "Successfully gets the token",
			store: mockStore{
				user: User{Username: "demo"},
			},
			encrypt: mockEncrypt{
				usernameMatches: true,
			},
			jwt: mockJwt{
				token: "jwt token",
				err:   nil,
			},
			tokenExpected: "jwt token",
			errorExpected: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			_, service := getService(tc.store, tc.encrypt, tc.jwt)

			token, err := service.Login(User{})

			assert.Equal(t, tc.tokenExpected, token)
			assert.Equal(t, tc.errorExpected, err)
		})
	}

}

func TestService_Verify(t *testing.T) {
	t.Run("A user exists", func(t *testing.T) {
		_, service := getService(
			mockStore{user: User{Username: "demo"}},
			mockEncrypt{},
			mockJwt{})

		assert.True(t, service.Verify("demo-username"))
	})

	t.Run("No user exists", func(t *testing.T) {
		_, service := getService(
			mockStore{},
			mockEncrypt{},
			mockJwt{})

		assert.False(t, service.Verify("demo-username"))
	})
}

func TestService_GetUser(t *testing.T) {
	t.Run("User does not exist", func(t *testing.T) {
		service := getServiceWithTags(mockStore{})
		user, err := service.GetUser("xusername")

		assert.Equal(t, EmptyUser, user)
		assert.Equal(t, errors.New("user not found for username xusername"), err)
	})

	t.Run("User exists", func(t *testing.T) {
		service := getServiceWithTags(mockStore{
			user: User{
				Username: "xusername",
				Email:    "email@email.com",
				Name:     "Jane Doe",
				Password: "garbage",
			},
		})

		user, err := service.GetUser("xusername")

		assert.Equal(t, "xusername", user.Username)
		assert.Equal(t, "email@email.com", user.Email)
		assert.Equal(t, 2, user.Tags.Size)
		assert.ElementsMatch(t, []tag.Tag{
			{Name: "react"},
			{Name: "python"},
		}, user.Tags.Tags)
		assert.Equal(t, nil, err)
	})
}

func TestService_the_rest(t *testing.T) {
	t.Run("AddTagsForUser", func(t *testing.T) {
		service := getServiceWithTags(mockStore{})

		resultTrue := service.AddTagsForUser("will-return-true", []string{
			"react", "python",
		})

		resultFalse := service.AddTagsForUser("will-return-false", []string{
			"react", "python",
		})

		assert.True(t, resultTrue)
		assert.False(t, resultFalse)
	})

	t.Run("RemoveTagsFromUser", func(t *testing.T) {
		service := getServiceWithTags(mockStore{})

		resultTrue := service.RemoveTagsFromUser("will-return-true", []string{
			"react", "python",
		})

		resultFalse := service.RemoveTagsFromUser("will-return-false", []string{
			"react", "python",
		})

		assert.True(t, resultTrue)
		assert.False(t, resultFalse)
	})

	t.Run("DeleteToken", func(t *testing.T) {
		service := getServiceWithTags(mockStore{})

		resultTrue := service.DeleteToken("will-return-true")
		resultFalse := service.DeleteToken("will-return-false")

		assert.True(t, resultTrue)
		assert.False(t, resultFalse)
	})
}

func getService(store mockStore, encrypt mockEncrypt, jwt mockJwt) (*mockStore, *Service) {
	service := NewService(ServiceConfig{
		Store:   &store,
		Encrypt: &encrypt,
		Jwt:     &jwt,
	})
	return &store, service
}

func getServiceWithTags(store mockStore) *Service {
	service := NewService(ServiceConfig{
		Store:   &store,
		Encrypt: &mockEncrypt{},
		Jwt:     &mockJwt{},
		Tags:    &mockTags{},
	})
	return service
}

// mock zone

type mockStore struct {
	user User
}

func (s *mockStore) SaveUser(user User) bool {
	s.user = user
	return true
}

func (s *mockStore) FindUser(_ string) User {
	return s.user
}

type mockEncrypt struct {
	shouldError     bool
	usernameMatches bool
}

func (e *mockEncrypt) Encrypt(_ string) (string, error) {
	if e.shouldError {
		return "", errors.New("error occurred")
	}
	return "encrypted-password", nil
}

func (e *mockEncrypt) Matches(_, _ string) bool {
	return e.usernameMatches
}

type mockJwt struct {
	token string
	err   error
}

func (j *mockJwt) SignedToken(_ string) (string, error) {
	return j.token, j.err
}

func (j *mockJwt) SaveJwt(_ string) bool {
	return true
}
func (j *mockJwt) DoesJwtExist(_ string) bool {
	return false
}
func (j *mockJwt) DeleteJwt(token string) bool {
	return token == "will-return-true"
}

type mockTags struct {
}

func (t mockTags) AddTagsForUser(username string, newTagIds []string) bool {
	return username == "will-return-true"
}
func (t mockTags) RemoveTagsFromUser(username string, removable []string) bool {
	return username == "will-return-true"
}
func (t mockTags) GetUserTags(username string) []string {
	return []string{
		"react", "python",
	}
}
func (t mockTags) GetTagById(name string) tag.Tag {
	if name == "react" {
		return tag.Tag{
			Name: "react",
		}
	} else if name == "python" {
		return tag.Tag{
			Name: "python",
		}
	}
	return tag.EmptyTag
}
