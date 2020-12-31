package mongo

import (
		"context"
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/mongo/options"
		"log"
		"sync"
)

var session *mongo.Client
var db *mongo.Database

func NewDb(uri string) *mongo.Database {
		var connectOnce sync.Once
		connectOnce.Do(func() {
				var err error
				opts := options.Client().ApplyURI(uri)
				session, err = mongo.NewClient(opts)
				if err != nil {
						log.Fatal(err)
				}
				err = session.Connect(context.TODO())
				if err != nil {
						log.Fatal(err)
				}

				db = session.Database("secrets")
				log.Println("DB connected")
		})

		return db
}

func Close() {
		if session != nil {
				ctx := context.TODO()
				func() {
						if err := session.Disconnect(ctx); err != nil {
								log.Println(err)
						} else {
								log.Println("DB disconnected")
						}
				}()
		}
}
