package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/seblw/wp-atrd-task/server"
)

func main() {
	api := server.NewServer(mux.NewRouter())

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      api.Router,
	}
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
