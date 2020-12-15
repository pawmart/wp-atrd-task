package rest

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func StartWebServer() {
	// define all REST API routes
	r := mux.NewRouter()
	r.Handle("/test", http.HandlerFunc(NewSecret)).Methods("GET")

	// add logging to our router
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// start HTTP server
	logrus.Fatal(http.ListenAndServe(":3000", loggedRouter))
}
