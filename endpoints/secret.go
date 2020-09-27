package endpoints

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"wp-atrd-task/connectors"
	"wp-atrd-task/models"
	"wp-atrd-task/validators"
)

type Secreter interface {
	GET(w http.ResponseWriter, r *http.Request)
	POST(w http.ResponseWriter, r *http.Request)
}

type secret struct {
	secretModel models.Secreter
}

func NewSecretEP(connector connectors.RedisConnector) Secreter {
	return &secret{secretModel: models.NewSecretModel(connector)}
}

func (s secret) GET(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	b, err := s.secretModel.GetSecret(hash)
	if err != nil {
		Response(w, r, err.Error(), http.StatusNotFound)
		return
	}

	Response(w, r, b, http.StatusOK)
}

func (s secret) POST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		Response(w, r, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	v := validators.FormValidator(r)
	if len(v) != 0 {
		Response(w, r, fmt.Sprintf("required formData are missing: %v", v), http.StatusBadRequest)
		return
	}

	if !validators.ExpireViewsValidator(r) {
		Response(w, r, fmt.Sprintf("expireAfterViews must be greater than 0"), http.StatusBadRequest)
		return
	}

	cs, err := s.secretModel.CreateSecret(
		r.FormValue("secret"),
		r.FormValue("expireAfterViews"),
		r.FormValue("expireAfter"))

	if err != nil {
		Response(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	Response(w, r, cs, http.StatusOK)
}
