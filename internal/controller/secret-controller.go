package controller

import (
	"net/http"

	"github.com/alkmc/wp-atrd-task/internal/entity"
	"github.com/alkmc/wp-atrd-task/internal/responder"
	"github.com/alkmc/wp-atrd-task/internal/service"
	"github.com/alkmc/wp-atrd-task/internal/validator"

	"github.com/go-chi/chi"
)

type secretController struct {
	secretService   service.Service
	secretValidator validator.Validator
}

//NewController returns Secret Controller
func NewController(s service.Service, v validator.Validator) Controller {
	return &secretController{
		secretService:   s,
		secretValidator: v,
	}
}

func (c *secretController) GetSecretByHash(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if err := c.secretValidator.Hash(hash); err != nil {
		responder.WithError(w, r, http.StatusBadRequest, "invalid hash")
		return
	}

	secret, err := c.secretService.FindAndUpdate(hash)
	if err != nil {
		responder.WithError(w, r, http.StatusNotFound, "Secret not found")
		return
	}
	secret.Respond(w, r)
}

func (c *secretController) AddSecret(w http.ResponseWriter, r *http.Request) {
	fd := getFormData(r)

	fi, err := c.secretValidator.FormData(fd)
	if err != nil {
		responder.WithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	s := &entity.Secret{
		SecretText:     fd.Secret,
		RemainingViews: fi.ExpireAfterViews,
		ExpireAfter:    fi.ExpireAfter,
	}

	if err := c.secretService.Create(s); err != nil {
		responder.WithError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	s.Respond(w, r)
}

func getFormData(r *http.Request) *entity.FormSecret {
	return &entity.FormSecret{
		Secret:           r.PostFormValue("secret"),
		ExpireAfter:      r.PostFormValue("expireAfter"),
		ExpireAfterViews: r.PostFormValue("expireAfterViews"),
	}
}
