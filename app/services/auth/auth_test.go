package auth_test

import (
	"fmt"
	"testing"
	"time"
	"urlshortener/app/services/auth"

	"github.com/stretchr/testify/assert"
)

var _ auth.UserStorer = &Store{}

type Store struct {
	user *auth.User
}

func InitStore() Store {
	s := Store{}

	s.user = &auth.User{
		Login:    "tets",
		Password: "12345",
	}

	return s
}

func (s *Store) Insert(user *auth.User) error {
	s.user = user
	return nil
}

func (s *Store) Get(login, password string) (*auth.User, error) {
	if login == s.user.Login {
		return s.user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *Store) ExistUserByLogin(login string) (bool, error) {
	if login == s.user.Login {
		return true, nil
	}

	return false, nil
}

func TestSignUpOK(t *testing.T) {
	store := InitStore()
	service := auth.NewAuthorizer(&store, "hashSalt", []byte("signingKey"), 1*time.Second)

	user := &auth.User{
		Login:    "test",
		Password: "12345",
	}

	_, err := service.SignUp(user)

	assert.Equal(t, err, nil)
}

func TestSignInOK(t *testing.T) {
	store := InitStore()
	service := auth.NewAuthorizer(&store, "hashSalt", []byte("signingKey"), 1*time.Second)

	store.user = &auth.User{
		Login:    "test",
		Password: "12345",
	}

	_, err := service.SignIn(store.user)
	assert.Equal(t, err, nil)
}

func TestSignUpFAIL(t *testing.T) {
	store := InitStore()
	service := auth.NewAuthorizer(&store, "hashSalt", []byte("signingKey"), 1*time.Second)

	store.user = &auth.User{
		Login:    "test",
		Password: "125",
	}
	_, err := service.SignUp(store.user)

	assert.Equal(t, err, fmt.Errorf("bad password"))
}

func TestSignInFAIL(t *testing.T) {
	store := InitStore()
	service := auth.NewAuthorizer(&store, "hashSalt", []byte("signingKey"), 1*time.Second)

	user := &auth.User{
		Login:    "test",
		Password: "12345",
	}

	_, err := service.SignUp(user)
	assert.Equal(t, err, nil)

	u := &auth.User{
		Login:    "test",
		Password: "1235",
	}
	_, err = service.SignIn(u)
	assert.Equal(t, nil, err)
}
