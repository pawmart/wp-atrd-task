package mongostore

import (
	"context"
	"fmt"
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
func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		SecretCollection: db.Collection(CollectionSecret),
	}
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

// GetSecret returns secret with given hash, if not found returns ErrNotFound
func (s *Storage) GetSecret(hash string) (*listing.Secret, error) {

	objectID, err := primitive.ObjectIDFromHex(hash)

	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"$and": bson.A{
			bson.M{"_id": objectID},
			bson.M{"remainingViews": bson.M{"$gt": 0}},
			bson.M{"expiresAt": bson.M{"$gt": time.Now()}},
		},
	}

	result := s.SecretCollection.FindOne(context.TODO(), filter)

	if result.Err() == mongo.ErrNoDocuments {
		return nil, listing.ErrNotFound
	}

	secret := Secret{}

	if err := result.Decode(&secret); err != nil {
		return nil, err
	}

	update := s.SecretCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		bson.M{"$set": bson.M{"remainingViews": secret.RemainingViews - 1}},
	)

	if update.Err() != nil {
		fmt.Println(update.Err())
		return nil, update.Err()
	}

	secret.RemainingViews = secret.RemainingViews - 1
	listingSecret := secret.ToListingSecret()
	return &listingSecret, nil
}

func calculateExpireDate(createdAt time.Time, expireTime int32) time.Time {
	if expireTime == 0 {
		return time.Date(7777, 7, 7, 7, 7, 7, 7, time.UTC)
	}

	return createdAt.Add(time.Duration(expireTime) * time.Minute)
}
