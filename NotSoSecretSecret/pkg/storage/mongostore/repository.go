package mongostore

import (
	"context"
	"notsosecretsercet/pkg/adding"
	"notsosecretsercet/pkg/listing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CollectionSecret identifier of MongoDB colection
const CollectionSecret = "Secrets"

// NewStorage returns a new MongoDB storage
func NewStorage(db *mongo.Database) (*Storage, error) {
	return &Storage{
		SecretCollection: db.Collection(CollectionSecret),
	}, nil
}

// Storage stores Secrets data in mongodb
type Storage struct {
	SecretCollection *mongo.Collection
}

// AddSecret stores secret in database
func (s *Storage) AddSecret(secret adding.Secret) (*listing.Secret, error) {

	createdAt := time.Now()

	secretToStore := Secret{
		SecretText:     secret.SecretText,
		CreatedAt:      createdAt,
		ExpiresAt:      calculateExpireDate(createdAt, secret.ExpireAfter),
		RemainingViews: secret.ExpireAfterViews,
	}

	result, err := s.SecretCollection.InsertOne(context.TODO(), secretToStore)

	if err != nil {
		return nil, err
	}

	secretToStore.Hash = result.InsertedID.(primitive.ObjectID)
	secretToReturn := secretToStore.ToListingSecret()
	return &secretToReturn, nil
}

// GetSecretByHash returns secret with given hash, if not found returns ErrNotFound
func (s *Storage) GetSecretByHash(hash string) (*listing.Secret, error) {

	objectID, err := primitive.ObjectIDFromHex(hash)

	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":            objectID,
		"remainingViews": bson.M{"gt": 0},
		"expiresAt":      bson.M{"$gt": time.Now().Unix()},
	}

	result := s.SecretCollection.FindOne(context.TODO(), filter)

	secret := Secret{}

	if err := result.Decode(&secret); err != nil {
		return nil, err
	}

	listingSecret := secret.ToListingSecret()
	return &listingSecret, nil
}

// GetSecret trys to get a secret by hash
func (s *Storage) GetSecret(hash string) (*listing.Secret, error) {

	secret, err := s.GetSecretByHash(hash)

	if err != nil {
		return nil, err
	}

	return secret, nil
}

func calculateExpireDate(createdAt time.Time, expireTime int32) time.Time {
	if expireTime == 0 {
		return time.Date(7777, 7, 7, 7, 7, 7, 7, time.UTC)
	}

	return createdAt.Add(time.Duration(expireTime) * time.Minute)
}
