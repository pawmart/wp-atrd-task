package server

import (
	"github.com/gorilla/mux"

	"github.com/seblw/wp-atrd-task/store"
)

// Server represents HTTP server.
type Server struct {
	Router *mux.Router
	Store  store.Storer
}

// NewServer initializes new HTTP Server.
func NewServer(r *mux.Router, st store.Storer) *Server {
	s := &Server{
		Router: r,
		Store:  st,
	}

	s.Router.HandleFunc("/v1/secret", s.HandleCreateSecret()).Methods("POST")
	s.Router.HandleFunc("/v1/secret/{id}", s.HandleGetSecretByID()).Methods("GET")

	return s
}
