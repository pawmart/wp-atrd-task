package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
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

		sc := Secret{
			Content:        msg.Secret,
			RemainingViews: int32(expireAfterViews),
			CreatedAt:      now,
		}

		if msg.ExpireAfter != "" {
			expireAfter, err := strconv.Atoi(msg.ExpireAfter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			t := now.Add(time.Duration(expireAfter) * time.Minute)
			sc.ExpiresAt = &t
		}

		o, err := s.Store.Insert(context.Background(), convToStore(sc))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		PrepareResponse(req, w, convToAPI(o))
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

		if sc.RemainingViews == 0 || (sc.ExpiresAt != nil && sc.ExpiresAt.Before(time.Now())) {
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

		PrepareResponse(req, w, convToAPI(scu))
	}
}

// PrepareResponse prepares HTTP response with content negotiation.
func PrepareResponse(r *http.Request, w http.ResponseWriter, s Secret) {
	const (
		ContentTypeApplicationXML  = "application/xml"
		ContentTypeApplicationJSON = "application/json"
	)

	if r.Header.Get("Accept") == ContentTypeApplicationXML {
		sj, err := xml.Marshal(s)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal secret: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", ContentTypeApplicationXML)
		w.Write(sj)
		return
	}

	sj, err := json.Marshal(s)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal secret: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentTypeApplicationJSON)
	w.Write(sj)
}

// Secret representation for API layer.
type Secret struct {
	ID             string     `json:"id" xml:"id,attr"`
	Content        string     `json:"content" xml:"content,attr" `
	RemainingViews int32      `json:"remaining_views" xml:"remaining_views,attr"`
	CreatedAt      time.Time  `json:"created_at" xml:"created_at,attr"`
	ExpiresAt      *time.Time `json:"expires_at" xml:"expires_at,attr"`
}
