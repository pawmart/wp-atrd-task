package controller

import "net/http"

//Controller is responsible for controlling the application logic
type Controller interface {
	GetSecretByHash(w http.ResponseWriter, r *http.Request)
	AddSecret(w http.ResponseWriter, r *http.Request)
}
