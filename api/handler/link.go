package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
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

type LinkRouter struct {
	r             *Router
	link          *link.LinkService
	serverAddress string
}

func NewLinkRouter(r *Router, link *link.LinkService, serverAddress string) *LinkRouter {
	return &LinkRouter{
		r:             r,
		link:          link,
		serverAddress: serverAddress,
	}
}

func (l *LinkRouter) RegisterAPI() {
	l.r.Route("/link", func(r chi.Router) {
		r.Use(AuthMiddleware)

		r.Post("/", l.addLink)
		r.Delete("/", l.delLink)
		r.Get("/myLinks", l.getLinks)
		r.Get("/", l.infoLink)
	})

	l.r.Get("/shortlink/{shortLink}", l.getLongLink)
}

func (l *LinkRouter) getLinks(w http.ResponseWriter, r *http.Request) {
	userID, err := l.r.GetUserAuth(r)
	if err != nil {
		http.Error(w, "can't get user from authorization", http.StatusUnauthorized)
		return
	}

	links, err := l.link.GetLinks(*userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't get links for user with id %s", userID), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(links)
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
			ShortLink: "https://url-shortener212.herokuapp.com/shortlink/" + shortLink,
		},
	)
}

func (l *LinkRouter) delLink(w http.ResponseWriter, r *http.Request) {
	userID, err := l.r.GetUserAuth(r)
	if err != nil {
		http.Error(w, "can't get user from authorization", http.StatusUnauthorized)
		return
	}

	linkID := r.URL.Query().Get("linkID")
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
	userID, err := l.r.GetUserAuthFromCookie(r)
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

	pwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "can't get path to current directory", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(pwd + "/front/statistic.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statistic.ShortLink = "https://url-shortener212.herokuapp.com/shortlink/" + statistic.ShortLink

	err = tmpl.Execute(w, statistic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (l *LinkRouter) getLongLink(w http.ResponseWriter, r *http.Request) {
	shortLink := chi.URLParam(r, "shortLink")
	if shortLink == "" {
		http.Error(w, "parameters 'shortLink' is empty", http.StatusBadRequest)
	}

	forwarder := r.Header.Get("X-FORWARDED-FOR")

	longLink, err := l.link.GetLongLink(shortLink, forwarder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longLink, http.StatusTemporaryRedirect)
}
