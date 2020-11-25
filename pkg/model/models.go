package model

import (
	"fmt"
	"time"

	guuid "github.com/google/uuid"
)

type Secret struct {
	Hash           string     `json:"hash"`
	SecretText     string     `json:"secretText"`
	CreatedAt      time.Time  `json:"createdAt"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	RemainingViews int        `json:"remainingViews"`
}

type SecretResponse struct {
	Hash           string     `json:"hash"`
	SecretText     string     `json:"secretText"`
	CreatedAt      time.Time  `json:"createdAt"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
	RemainingViews int        `json:"remainingViews"`
}

type FormData struct {
	Secret           string
	ExpireAfterViews int
	ExpireAfter      int
}

func NewSecret(formData FormData) Secret {
	createdAt := time.Now()
	var expiresAtPtr *time.Time
	if formData.ExpireAfter > 0 {
		expiresAt := createdAt.Add(time.Duration(formData.ExpireAfter) * time.Minute)
		expiresAtPtr = &expiresAt
	}
	return Secret{
		Hash:           genUUID(),
		RemainingViews: formData.ExpireAfterViews,
		SecretText:     formData.Secret,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAtPtr,
	}
}

func NewSecretResponse(secret Secret) SecretResponse {
	return SecretResponse{
		Hash:           secret.Hash,
		CreatedAt:      secret.CreatedAt,
		ExpiresAt:      secret.ExpiresAt,
		RemainingViews: secret.RemainingViews,
		SecretText:     secret.SecretText,
	}
}

func genUUID() string {
	id := guuid.New()
	return fmt.Sprintf("%s", id.String())
}
