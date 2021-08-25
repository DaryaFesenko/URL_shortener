package handler

import (
	"fmt"

	"urlshortener/app/services/auth"

	"github.com/dgrijalva/jwt-go/v4"
	uuid "github.com/satori/go.uuid"
)

func ParseToken(accessToken string, signingKey []byte) (*uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return &claims.UserID, nil
	}

	return nil, fmt.Errorf("invalid access token")
}
