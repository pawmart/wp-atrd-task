package server

import "net/http"

// CreateSecret handles CreateSecret request.
func (s *Server) CreateSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("CreateSecret called\n"))
	}
}

// GetSecret handles GetSecret request.
func (s *Server) GetSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("GetSecret called\n"))
	}
}
