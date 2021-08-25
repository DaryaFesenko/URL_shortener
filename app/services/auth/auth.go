package auth

import (
	"github.com/dgrijalva/jwt-go/v4"
	uuid "github.com/satori/go.uuid"
)

type Claims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userId"`
}
