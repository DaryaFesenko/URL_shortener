package handler

import (
	"fmt"
	"net/http"
	"strings"
	"urlshortener/app/services/auth"
	"urlshortener/app/services/link"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Router struct {
	*chi.Mux
}

var (
	SIGNING_KEY []byte
)

const (
	CountParts = 2
	JWT        = "Bearer"
)

func NewRouter(auth *auth.AuthService, link *link.LinkService, signingKey []byte) *Router {
	SIGNING_KEY = signingKey
	r := &Router{
		chi.NewRouter(),
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	NewAuthRouter(r, auth).RegisterAPI()
	NewLinkRouter(r, link).RegisterAPI()

	return r
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "can't header authorization", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != CountParts {
				http.Error(w, "can't parse header", http.StatusUnauthorized)
				return
			}

			if headerParts[0] != JWT {
				http.Error(w, "bad header", http.StatusUnauthorized)
				return
			}

			_, err := ParseToken(headerParts[1], SIGNING_KEY)
			if err != nil {
				http.Error(w, "bad header", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func (*Router) GetUserAuth(r *http.Request) (*uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("can't header authorization")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != CountParts {
		return nil, fmt.Errorf("can't parse header")
	}

	if headerParts[0] != JWT {
		return nil, fmt.Errorf("bad header")
	}

	userID, err := ParseToken(headerParts[1], SIGNING_KEY)
	if err != nil {
		return nil, fmt.Errorf("bad header")
	}

	return userID, nil
}
