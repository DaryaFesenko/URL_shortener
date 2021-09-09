package auth

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type UserStorer interface {
	Insert(user *User) error
	Select(query string, params ...interface{}) ([]User, error)
	Update(query string, params ...interface{}) error
	Get(userID uuid.UUID, password string) (*User, error)
}

type AuthService struct {
	signingKey     []byte
	hashSalt       string
	expireDuration time.Duration

	userStore UserStorer
}

func NewAuthorizer(store UserStorer, hashSalt string, signingKey []byte, expireDuration time.Duration) *AuthService {
	return &AuthService{
		userStore:      store,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *AuthService) SignUp(user *User) (string, error) {
	if user.Password == "" || len(user.Password) < 4 {
		return "", fmt.Errorf("bad password")
	}

	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	err := a.userStore.Insert(user)
	if err != nil {
		return "", err
	}

	token := a.createToken(user.ID)

	return token.SignedString(a.signingKey)
}

func (a *AuthService) SignIn(user *User) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userStore.Get(user.ID, user.Password)
	if err != nil {
		return "", err
	}

	token := a.createToken(user.ID)

	return token.SignedString(a.signingKey)
}

func (a *AuthService) createToken(userID uuid.UUID) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserID: userID,
	})

	return token
}
