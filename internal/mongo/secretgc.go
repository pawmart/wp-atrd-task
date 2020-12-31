package mongo

import (
		"context"
		"go.mongodb.org/mongo-driver/bson"
		"go.mongodb.org/mongo-driver/mongo"
		"log"
		"time"
)

type SecretGCRepository struct {
		collection *mongo.Collection
}

func (s *SecretGCRepository) GC() error {
		filter := bson.M{
				"$or": []bson.M{
						{
								"$and": []bson.M{
										{"expirationdate": bson.M{"$lt": time.Now().Unix()}},
										{"expirationdate": bson.M{"$gt": 0}},
								},
						},
						{"viewsleft": bson.M{"$lt": 1}},
				},
		}

		res, err := s.collection.DeleteMany(context.TODO(), filter)
		if res != nil {
				log.Println("GC:", res.DeletedCount, "removed")
				return nil
		}

		return err
}

func NewSecretGCRepository(db *mongo.Database) *SecretGCRepository {
		col := db.Collection("secrets")
		sr := SecretGCRepository{
				collection: col,
		}

		return &sr
}
