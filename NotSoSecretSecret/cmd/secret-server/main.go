package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notsosecretsercet/pkg/adding"
	"notsosecretsercet/pkg/config"
	"notsosecretsercet/pkg/http/rest"
	"notsosecretsercet/pkg/listing"
	"notsosecretsercet/pkg/storage/mongostore"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DbURI))

	if err != nil {
		log.Fatalf("Unable to connect to database. Secret server shuts down :(")
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("There was a problem durning mongo clien disconecting :(")
		}
	}()

	pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCtxCancel()

	err = client.Ping(pingCtx, readpref.Primary())

	if err != nil {
		log.Fatalf("Unable to connect to database. Secret server shuts down :(")
	}

	db := client.Database(config.DbName)

	repo := mongostore.NewStorage(db)
	addingService := adding.NewService(repo)
	listingService := listing.NewService(repo)

	router := rest.Handler(addingService, listingService)

	port := ":" + config.Port
	fmt.Println("")
	fmt.Println("The secret server is on tap now: http://localhost" + port)

	http.ListenAndServe(port, router)
}
