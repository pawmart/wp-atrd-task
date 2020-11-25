package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/maciejem/secret/db"
	"github.com/maciejem/secret/pkg/apperrors"
	"github.com/maciejem/secret/pkg/model"
)

var secretKey = "secret"
var dbInstance db.Database

func NewHandler(db db.Database) http.Handler {
	dbInstance = db
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeForm))

	r.Route("/v1/secret", func(r chi.Router) {
		r.Post("/", createSecret)
		r.Route("/{secretID}", func(r chi.Router) {
			r.Use(secretCtx)
			r.Get("/", getSecret)
		})
	})

	return r
}

func secretCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretID := chi.URLParam(r, "secretID")
		if secretID == "" {
			render.Render(w, r, ErrNotFound)
			return
		}
		secret, err := dbInstance.GetSecretById(secretID)
		if err != nil {
			if errors.Is(err, apperrors.ErrNoMatch) {
				render.Render(w, r, ErrNotFound)
				return
			}
			render.Render(w, r, ErrInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), secretKey, secret)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	secretText := r.FormValue("secret")
	expireAfterViews := r.FormValue("expireAfterViews")
	expireAfter := r.FormValue("expireAfter")
	if secretText == "" || expireAfterViews == "" || expireAfter == "" {
		render.Render(w, r, ErrInvalidInput)
		return
	}
	expireAfterViewsInt, err := strconv.Atoi(expireAfterViews)
	if err != nil {
		render.Render(w, r, ErrInvalidInput)
		return
	}
	if expireAfterViewsInt <= 0 {
		render.Render(w, r, ErrInvalidInput)
		return
	}
	expireAfterInt, err := strconv.Atoi(expireAfter)
	if err != nil {
		render.Render(w, r, ErrInvalidInput)
		return
	}
	formData := model.FormData{
		Secret:           secretText,
		ExpireAfter:      expireAfterInt,
		ExpireAfterViews: expireAfterViewsInt,
	}
	secret := model.NewSecret(formData)
	err = dbInstance.CreateSecret(secret)
	if err != nil {
		render.Render(w, r, ErrInvalidInput)
		return
	}
	renderSecret(w, r, secret)
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	secret := r.Context().Value(secretKey).(model.Secret)
	if secret.ExpiresAt != nil && secret.ExpiresAt.Before(time.Now()) {
		render.Render(w, r, ErrNotFound)
		return
	}
	if secret.RemainingViews < 1 {
		render.Render(w, r, ErrNotFound)
		return
	}
	err := dbInstance.DecreaseSecretRemainingViewsById(secret.Hash)
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}
	secret.RemainingViews--
	renderSecret(w, r, secret)
}

func renderSecret(w http.ResponseWriter, r *http.Request, secret model.Secret) {
	secretResponse := model.NewSecretResponse(secret)
	urlFormat, _ := r.Context().Value(middleware.URLFormatCtxKey).(string)
	switch urlFormat {
	case "json":
		render.JSON(w, r, secretResponse)
	case "xml":
		render.XML(w, r, secretResponse)
	default:
		render.JSON(w, r, secretResponse)
	}
	render.Status(r, http.StatusOK)
}
