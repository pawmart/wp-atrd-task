package entity

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alkmc/wp-atrd-task/internal/responder"
)

//Secret entity is the core business object
type Secret struct {
	Hash           string     `json:"hash" xml:"hash"`
	SecretText     string     `json:"secretText" xml:"secretText"`
	CreatedAt      time.Time  `json:"createdAt" xml:"createdAt"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty" xml:"expiresAt,omitempty"`
	RemainingViews int32      `json:"remainingViews" xml:"remainingViews"`
	ExpireAfter    int32      `json:"-" xml:"-"`
}

//Respond serializes the Secret entity into the response body
func (s *Secret) Respond(w http.ResponseWriter, r *http.Request) {
	responder.Response(w, r, http.StatusOK, s)
}

//CalculateExpiration calculates expiration based on CreatedAt and ExpireAfter fields
func (s *Secret) CalculateExpiration() {
	if s.ExpireAfter != 0 {
		x := s.CreatedAt.Add(time.Duration(s.ExpireAfter) * time.Minute)
		s.ExpiresAt = &x
	}
}

//CastToDuration casts ExpireAfter field to Duration type
func (s *Secret) CastToDuration() time.Duration {
	return time.Duration(s.ExpireAfter) * time.Minute
}

//NewExpirationAt calculates new expiration time based on ExpiresAt field
func (s *Secret) NewExpirationAt() time.Duration {
	if s.ExpiresAt == nil {
		return 0
	}
	fmt.Println(s.ExpiresAt.Sub(time.Now()))
	return s.ExpiresAt.Sub(time.Now())
}
