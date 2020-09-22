package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/seblw/wp-atrd-task/server"
	"github.com/seblw/wp-atrd-task/store"
)

var addr = flag.String("addr", ":8080", "Address for listening")

func main() {
	flag.Parse()

	db, err := initDB(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = backoff.Retry(func() error {
		err := db.Ping()
		if err != nil {
			log.Printf("Failed to ping DB (%s)\n", err)
			return err
		}
		return nil
	}, backoff.NewExponentialBackOff())
	if err != nil {
		panic(err)
	}
	log.Println("Connection to DB initialized!")

	ctx := context.Background()
	s, err := store.New(ctx, db)
	if err != nil {
		panic(err)
	}

	api := server.NewServer(mux.NewRouter(), s)

	srv := &http.Server{
		Addr:         *addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      api.Router,
	}

	// TODO: Add graceful shutdown.
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initDB(host, port, username, password, database string) (*sqlx.DB, error) {
	s := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	return sqlx.Open("postgres", s)
}
