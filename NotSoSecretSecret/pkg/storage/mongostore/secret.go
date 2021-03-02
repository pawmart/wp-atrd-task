package mongostore

import (
	"notsosecretsercet/pkg/listing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Secret defines properties of secret to be stored in MongoDB
type Secret struct {
	Hash           primitive.ObjectID `bson:"_id,omitempty"`
	SecretText     string             `bson:"secretText,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty"`
	ExpiresAt      time.Time          `bson:"expiresAt,omitempty"`
	RemainingViews int32              `bson:"remainingViews,omitempty"`
}

// ToListingSecret returns a listing.Secret
func (s *Secret) ToListingSecret() listing.Secret {
	return listing.Secret{
		Hash:           s.Hash.Hex(),
		SecretText:     s.SecretText,
		CreatedAt:      s.CreatedAt,
		ExpiresAt:      s.ExpiresAt,
		RemainingViews: s.RemainingViews,
	}
}
