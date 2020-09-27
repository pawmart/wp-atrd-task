package docs

import "time"

type Secret struct {
	Hash           string    `json:"hash"`
	SecretText     string    `json:"secretText"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
	RemainingViews int32     `json:"remainingViews"`
}
