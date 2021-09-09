package auth

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userId"`
}
