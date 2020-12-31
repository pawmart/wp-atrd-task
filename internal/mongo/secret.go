package mongo

import (
		"context"
		"github.com/mkalafior/wp-atrd-task/internal/app/secrets"
		"go.mongodb.org/mongo-driver/bson"
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/mongo/options"
		"time"
)

type SecretRepository struct {
		collection *mongo.Collection
}

type secretMongo struct {
		Hash           string
		CreatedAt      int
		ExpirationDate int
		Secret         string
		ViewsLeft      int
}

func (s *SecretRepository) Find(hash string) (secrets.Secret, error) {
		filter := bson.M{
				"hash": hash,
				"viewsleft": bson.M{
						"$gt": 0,
				},
				"$or": []bson.M{
						{"expirationdate": bson.M{"$gt": time.Now().Unix()}},
						{"expirationdate": 0},
				},
		}
		var sm secretMongo
		var secret secrets.Secret
		update := bson.D{{"$inc", bson.M{"viewsleft": -1}}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err := s.collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&sm)

		if err != nil {
				return secret, err
		}

		secret.Secret = sm.Secret
		secret.Expire.ViewsLeft = sm.ViewsLeft
		secret.Expire.Date = time.Unix(int64(sm.ExpirationDate), 0)
		secret.Created = time.Unix(int64(sm.CreatedAt), 0)
		secret.Hash = secrets.Hash(sm.Hash)

		return secret, nil
}

func (s *SecretRepository) Store(secret *secrets.Secret) error {
		opts := options.Update().SetUpsert(true)
		filter := bson.M{
				"hash": string(secret.Hash),
		}
		sm := secretMongo{
				Hash:           string(secret.Hash),
				CreatedAt:      int(secret.Created.Unix()),
				ExpirationDate: int(secret.Expire.Date.Unix()),
				Secret:         secret.Secret,
				ViewsLeft:      secret.Expire.ViewsLeft,
		}
		update := bson.D{
				{"$set", sm},
		}
		_, err := s.collection.UpdateOne(context.TODO(), filter, update, opts)

		return err
}

func NewSecretRepository(db *mongo.Database) *SecretRepository {
		col := db.Collection("secrets")
		sr := SecretRepository{
				collection: col,
		}

		return &sr
}
