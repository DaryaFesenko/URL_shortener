package handler

import (
	"fmt"
	"html/template"
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

func NewRouter(auth *auth.AuthService, link *link.LinkService, signingKey []byte, serverAddress string) *Router {
	SIGNING_KEY = signingKey
	r := &Router{
		chi.NewRouter(),
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", r.handleMain)
	r.Get("/stat", r.handleStat)

	NewAuthRouter(r, auth).RegisterAPI()
	NewLinkRouter(r, link, serverAddress).RegisterAPI()

	return r
}

func (o *Router) handleMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../front/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (o *Router) handleStat(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../front/statistic.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (*Router) GetUserAuthFromCookie(r *http.Request) (*uuid.UUID, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, fmt.Errorf("can't header authorization")
		}
		return nil, fmt.Errorf("can't header authorization")
	}

	userID, err := ParseToken(cookie.Value, SIGNING_KEY)
	if err != nil {
		return nil, fmt.Errorf("bad header")
	}

	return userID, nil
}
