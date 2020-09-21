package storage

import (
	"github.com/google/uuid"
	"time"
)

// Storage interface for interacting with secret storage
type Storage interface {
	AddSecret(secretValue string, viewsLeft uint32, createdAt time.Time, expiresAt *time.Time) Secret
	GetSecret(id uuid.UUID) (*Secret, bool)
}

// Secret struct
type Secret struct {
	Id             uuid.UUID  `json:"hash" xml:"hash"`
	Value          string     `json:"secretText" xml:"secretText"`
	CreatedAt      time.Time  `json:"createdAt" xml:"createdAt"`
	ExpiresAfter   *time.Time `json:"expiresAt" xml:"expiresAt"`
	RemainingViews uint32     `json:"remainingViews" xml:"remainingViews"`
}

// isFetchable checks if secret can be fetch depending on remaining views and expire date
func (s Secret) isFetchable() bool {
	if s.RemainingViews <= 0 {
		return false
	}

	if s.ExpiresAfter != nil && !time.Now().Before(*s.ExpiresAfter) {
		return false
	}

	return true
}
