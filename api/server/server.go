package server

import (
	"net/http"
	"time"
)

const (
	timeout = 30
)

type Server struct {
	srv http.Server
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       timeout * time.Second,
		WriteTimeout:      timeout * time.Second,
		ReadHeaderTimeout: timeout * time.Second,
	}
	return s
}

func (s *Server) Stop() {
}

func (s *Server) Start() {
}
