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
	apiVersion                 = "/v1"
)

// Handler returns secret service handler
func Handler(as adding.Service, ls listing.Service) http.Handler {
	router := chi.NewRouter()

	router.Route(apiVersion, func(r chi.Router) {
		r.Post("/secret", addSecret(as))
		r.Get("/secret/{secretHash}", getSecret(ls))
	})

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

		if newSecret.ExpireAfterViews < 0 {
			http.Error(w, "Invalid input", http.StatusMethodNotAllowed)
			return
		}

		createdSecret, err := as.AddSecret(newSecret)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		handleResponse(w, r, createdSecret)

	}
}

func getSecret(ls listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secretHash := chi.URLParam(r, "secretHash")

		secret, err := ls.GetSecret(secretHash)

		if err == listing.ErrNotFound {
			http.Error(w, "Secret not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		handleResponse(w, r, secret)
	}
}

func handleResponse(w http.ResponseWriter, r *http.Request, secret *listing.Secret) {
	if r.Header.Get("Accept") == ContentTypeApplicationXML {
		w.Header().Set("Content-Type", ContentTypeApplicationXML)
		xml.NewEncoder(w).Encode(secret)
		return
	}

	w.Header().Set("Content-Type", ContentTypeApplicationJSON)
	json.NewEncoder(w).Encode(secret)
}
