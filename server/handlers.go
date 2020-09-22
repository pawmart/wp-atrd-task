package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/seblw/wp-atrd-task/store"
)

const uuidLen = 36

// HandleCreateSecret handles CreateSecret request.
func (s *Server) HandleCreateSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		msg := &Message{
			Secret:           req.PostFormValue("secret"),
			ExpireAfterViews: req.PostFormValue("expireAfterViews"),
			ExpireAfter:      req.PostFormValue("expireAfter"),
		}

		if !msg.Validate() {
			// Return HTTP 405 on invalid input according to swagger.
			http.Error(w, msg.PrintErrors(), http.StatusMethodNotAllowed)
			return
		}

		now := time.Now()
		expireAfterViews, err := strconv.Atoi(msg.ExpireAfterViews)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expireAfter, err := strconv.Atoi(msg.ExpireAfter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Make ExpiresAt optional when expireAfter is zero.
		sc := Secret{
			Content:        msg.Secret,
			RemainingViews: int32(expireAfterViews),
			CreatedAt:      now,
			ExpiresAt:      now.Add(time.Duration(expireAfter) * time.Minute),
		}

		o, err := s.Store.Insert(context.Background(), convToStore(sc))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		oj, err := json.Marshal(convToAPI(o))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal secret"), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(oj)
	}
}

// HandleGetSecretByID handles GetSecret request.
func (s *Server) HandleGetSecretByID() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]

		if id == "" || len(id) != uuidLen {
			http.Error(w, fmt.Sprintf("invalid id: '%s'", id), http.StatusBadRequest)
			return
		}

		sc, err := s.Store.GetByID(context.Background(), id)
		if err != nil {
			http.Error(w, "failed to get secret by ID", http.StatusNotFound)
			return
		}

		if sc.RemainingViews == 0 || sc.ExpiresAt.Before(time.Now()) {
			if err := s.Store.Delete(context.Background(), id); err != nil {
				http.Error(w, fmt.Sprintf("failed to delete secret by ID: %v", err), http.StatusInternalServerError)
				return
			}
			http.Error(w, "failed to get secret by ID", http.StatusNotFound)
			return
		}

		sc.RemainingViews--
		scu, err := s.Store.Update(context.Background(), id, store.Secret{
			Content:        sc.Content,
			RemainingViews: sc.RemainingViews,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update secret by ID: %v", err), http.StatusInternalServerError)
			return
		}

		cj, err := json.Marshal(convToAPI(scu))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal secret: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(cj)
	}
}

// Secret ..
type Secret struct {
	ID             string    `json:"id"`
	Content        string    `json:"content"`
	RemainingViews int32     `json:"remaining_views"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}
