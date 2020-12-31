package secrets

import (
		"context"
		transport "github.com/go-kit/kit/transport/http"
		"github.com/gorilla/mux"
		"go.mongodb.org/mongo-driver/mongo"
		"log"
		"net/http"
)

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
		switch err {
		case ErrInvalidArgument:
				w.WriteHeader(http.StatusMethodNotAllowed)
		case mongo.ErrNoDocuments:
				w.WriteHeader(http.StatusNotFound)
		default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
		}
}

func MakeHandler(endpoints Endpoints) *mux.Router {
		r := mux.NewRouter()
		s := r.PathPrefix("/v1").Subrouter()

		opts := []transport.ServerOption{
				transport.ServerErrorEncoder(encodeError),
		}

		s.Handle("/secret", transport.NewServer(
				endpoints.NewSecretEndpoint,
				decodeNewSecretRequest,
				encodeResponse,
				opts...,
		)).Methods("POST")

		s.Handle("/secret/{hash}", transport.NewServer(
				endpoints.FetchEndpoint,
				decodeFetchSecretRequest,
				encodeResponse,
				opts...,
		)).Methods("GET")

		return r
}
