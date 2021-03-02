package rest

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"notsosecretsercet/pkg/adding"
	"notsosecretsercet/pkg/listing"

	"github.com/go-chi/chi/v5"
)

const (
	ContentTypeApplicationXML  = "application/xml"
	ContentTypeApplicationJSON = "application/json"
)

// Handler returns secret service handler
func Handler(as adding.Service, ls listing.Service) http.Handler {
	router := chi.NewRouter()

	router.Post("/secret", addSecret(as))
	router.Get("/secrets/{secretHash}", getSecret(ls))

	return router
}

func addSecret(as adding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		newSecret := adding.Secret{}

		if err := parseForm(r, &newSecret); err != nil {
			http.Error(w, "Invalid input", http.StatusMethodNotAllowed)
			return
		}

		if newSecret.ExpireAfterViews < 1 {
			http.Error(w, "Invalid input", http.StatusMethodNotAllowed)
			return
		}

		createdSecret, err := as.AddSecret(newSecret)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Accept") == ContentTypeApplicationXML {
			w.Header().Set("Content-Type", ContentTypeApplicationXML)
			xml.NewEncoder(w).Encode(createdSecret)
			return
		}

		w.Header().Set("Content-Type", ContentTypeApplicationJSON)
		json.NewEncoder(w).Encode(createdSecret)

	}
}

func getSecret(ls listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
