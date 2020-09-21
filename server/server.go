package server

import (
	"github.com/gorilla/mux"
)

// Server represents HTTP server.
type Server struct {
	Router *mux.Router
}

// NewServer initializes new HTTP Server.
func NewServer(r *mux.Router) *Server {
	s := &Server{
		Router: r,
	}

	s.Router.HandleFunc("/v1/secret/", s.CreateSecret()).Methods("POST")
	s.Router.HandleFunc("/v1/secret/{id}", s.GetSecret()).Methods("GET")

	return s
}
