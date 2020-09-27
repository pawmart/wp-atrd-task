package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"wp-atrd-task/config"
	"wp-atrd-task/connectors"
	"wp-atrd-task/endpoints"
)

type App struct {
	config   *config.Config
	router   *mux.Router
	redis    connectors.RedisConnector
	secretEP endpoints.Secreter
}

func NewApp() App {
	c, err := config.New("config/conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client := connectors.NewRedis(c)

	secretEP := endpoints.NewSecretEP(client)

	app := App{
		config:   c,
		router:   nil,
		redis:    client,
		secretEP: secretEP,
	}

	return app
}

func (a *App) AddHandlers() {
	r := mux.NewRouter()

	secret := r.PathPrefix("/v1").Subrouter().StrictSlash(true)
	secret.HandleFunc("/secret/{hash}", a.secretEP.GET).Methods("GET")
	secret.HandleFunc("/secret", a.secretEP.POST).Methods("POST")

	a.router = r
}

func (a *App) ListenAndServe() {
	s := &http.Server{
		Handler:      a.router,
		Addr:         fmt.Sprintf("%s:%s", a.config.Server.Address, a.config.Server.Port),
		ReadTimeout:  time.Duration(a.config.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(a.config.Server.Timeout) * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
