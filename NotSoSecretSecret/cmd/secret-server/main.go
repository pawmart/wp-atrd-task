package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notsosecretsercet/pkg/adding"
	"notsosecretsercet/pkg/http/rest"
	"notsosecretsercet/pkg/listing"
	"notsosecretsercet/pkg/storage/mongostore"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/secret-server/config/")
	viper.AddConfigPath("/cmd/secret-server/config/")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("%v", err)
	}
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("DB_URI")))

	if err != nil {
		log.Fatalf("Unable to connect to database. Secret server shuts down :(")
	}

	pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCtxCancel()

	err = client.Ping(pingCtx, readpref.Primary())

	if err != nil {
		log.Fatalf("Unable to connect to database. Secret server shuts down :(")
	}

	db := client.Database(viper.GetString("DB_NAME"))

	repo := mongostore.NewStorage(db)
	addingService := adding.NewService(repo)
	listingService := listing.NewService(repo)

	router := rest.Handler(addingService, listingService)

	fmt.Println("The secret server is on tap now: http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
