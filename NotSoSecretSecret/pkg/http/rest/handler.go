package rest

import (
	"net/http"
	"notsosecretsercet/pkg/adding"
	"notsosecretsercet/pkg/listing"

	"github.com/go-chi/chi/v5"
)

// Handler returns secret service handler
func Handler(as adding.Service, ls listing.Service) http.Handler {
	router := chi.NewRouter()

	return router
}

func addSecret(as adding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
