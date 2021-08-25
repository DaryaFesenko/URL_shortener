package handler

import (
	"net/http"
	"urlshortener/app/services/link"
)

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
	l.r.HandleFunc("/add_link", l.r.AuthMiddleware(http.HandlerFunc(l.addLink)).ServeHTTP)
}

func (l *LinkRouter) addLink(w http.ResponseWriter, r *http.Request) {
	_, err := l.r.GetUserAuth(w, r)
	if err != nil {
		http.Error(w, "can't get user by id", http.StatusUnauthorized)
		return
	}
}
