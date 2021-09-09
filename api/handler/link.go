package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"urlshortener/app/services/link"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type linkRequest struct {
	LongLink string `json:"longLink"`
}

type linkResponse struct {
	ShortLink string `json:"shortLink"`
}

type longLinkResponse struct {
	LongLink string `json:"longLink"`
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
		r.Delete("/{linkId}", l.delLink)
		r.Get("/{id}", l.infoLink)
	})

	l.r.Get("/{shortLink}", l.getLongLink)
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

func (l *LinkRouter) delLink(w http.ResponseWriter, r *http.Request) {
	userID, err := l.r.GetUserAuth(r)
	if err != nil {
		http.Error(w, "can't get user from authorization", http.StatusUnauthorized)
		return
	}

	linkID := r.URL.Query().Get("linkId")
	if linkID == "" {
		http.Error(w, "parameters 'linkId' is empty", http.StatusBadRequest)
		return
	}

	uuidLinkID, err := uuid.Parse(linkID)
	if err != nil {
		http.Error(w, "parameters 'linkId' type is not uuid", http.StatusBadRequest)
		return
	}

	err = l.link.DeleteLink(*userID, uuidLinkID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (l *LinkRouter) infoLink(w http.ResponseWriter, r *http.Request) {
	userID, err := l.r.GetUserAuth(r)
	if err != nil {
		http.Error(w, "can't get user from authorization", http.StatusUnauthorized)
		return
	}

	linkID := r.URL.Query().Get("id")
	if linkID == "" {
		http.Error(w, "parameters 'id' is empty", http.StatusBadRequest)
		return
	}

	uuidLinkID, err := uuid.Parse(linkID)
	if err != nil {
		http.Error(w, "parameters 'id' type is not uuid", http.StatusBadRequest)
		return
	}

	statistic, err := l.link.GetLinkStatistic(userID, uuidLinkID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(statistic)
}

func (l *LinkRouter) getLongLink(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.Query().Get("shortLink")
	if shortLink == "" {
		http.Error(w, "parameters 'shortLink' is empty", http.StatusBadRequest)
	}

	forwarder := r.Header.Get("X-FORWARDED-FOR")
	if forwarder == "" {
		forwarder = r.RemoteAddr
	}
	usedUserID, _, _ := net.SplitHostPort(forwarder)

	longLink, err := l.link.GetLongLink(shortLink, usedUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(longLinkResponse{
		LongLink: longLink,
	})
}
