package handler

import (
	"encoding/json"
	"net/http"
	"urlshortener/app/services/link"

	"github.com/go-chi/chi/v5"
)

type linkRequest struct {
	LongLink string `json:"longLink"`
}

type linkResponse struct {
	ShortLink string `json:"shortLink"`
}

type LinkRouter struct {
	r    *Router
	link *link.LinkService
}

func NewLinkRouter(r *Router, link *link.LinkService) *LinkRouter {
	return &LinkRouter{
		r:    r,
		link: link,
	}
}

func (l *LinkRouter) RegisterAPI() {
	l.r.Route("/link", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/", l.addLink)
	})
}

func (l *LinkRouter) addLink(w http.ResponseWriter, r *http.Request) {
	userID, err := l.r.GetUserAuth(r)
	if err != nil {
		http.Error(w, "can't get user from authorization", http.StatusUnauthorized)
		return
	}

	request := linkRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	shortLink, err := l.link.CreateLink(userID, request.LongLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(
		linkResponse{
			ShortLink: shortLink,
		},
	)
}
