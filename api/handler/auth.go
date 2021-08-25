package handler

import (
	"encoding/json"
	"net/http"
	"urlshortener/app/services/auth"

	uuid "github.com/satori/go.uuid"
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
	a.r.HandleFunc("/signIn", http.HandlerFunc(a.signIn))
	a.r.HandleFunc("/signUp", http.HandlerFunc(a.signUp))
}

func (a *AuthRouter) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	request := signInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	newUser := auth.User{
		ID:       uuid.NewV4(),
		Login:    request.Login,
		Password: request.Password,
	}

	if err := a.auth.SignUp(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := a.auth.SignIn(&newUser)
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
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

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
