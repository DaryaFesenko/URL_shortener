package handler

import (
	"encoding/json"
	"net/http"
	"urlshortener/app/services/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type signInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type AuthRouter struct {
	r    *Router
	auth *auth.AuthService
}

func NewAuthRouter(r *Router, auth *auth.AuthService) *AuthRouter {
	return &AuthRouter{
		r:    r,
		auth: auth,
	}
}

func (a *AuthRouter) RegisterAPI() {
	a.r.Route("/auth", func(r chi.Router) {
		r.Post("/signIn", a.signIn)

		r.Post("/signUp", a.signUp)
	})
}

func (a *AuthRouter) signUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request := signInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	newUser := auth.User{
		ID:       uuid.New(),
		Login:    request.Login,
		Password: request.Password,
	}

	token, err := a.auth.SignUp(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(
		signInResponse{
			Token: token,
		},
	)
}

func (a *AuthRouter) signIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request := signInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := auth.User{
		Login:    request.Login,
		Password: request.Password,
	}

	token, err := a.auth.SignIn(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(
		signInResponse{
			Token: token,
		},
	)
}
