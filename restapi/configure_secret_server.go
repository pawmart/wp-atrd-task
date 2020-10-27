// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	"wp-atrd-task/models"
	"wp-atrd-task/restapi/operations"
	"wp-atrd-task/restapi/operations/secret"
)

//go:generate swagger generate server --target ../../wp-atrd-task --name SecretServer --spec ../api/swagger/swagger.yml --principal interface{}

var secrets = make(map[string]models.Secret)
var secretsLock = &sync.Mutex{}

func configureFlags(api *operations.SecretServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SecretServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.UrlformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()
	api.XMLProducer = runtime.XMLProducer()

	api.SecretAddSecretHandler = secret.AddSecretHandlerFunc(func(params secret.AddSecretParams) middleware.Responder {
		now := time.Now()
		expires := now.Add(time.Minute * time.Duration(params.ExpireAfter))
		if params.ExpireAfter == 0 {
			expires = expires.AddDate(100, 0, 0) // good enough, this won't run for 100 years...
		}
		hash := uuid.New().String()

		newSecret := models.Secret{
			CreatedAt:      strfmt.DateTime(now),
			ExpiresAt:      strfmt.DateTime(expires),
			Hash:           hash,
			RemainingViews: params.ExpireAfterViews,
			SecretText:     params.Secret,
		}

		AddSecret(&newSecret)
		return secret.NewAddSecretOK().WithPayload(&newSecret)
	})

	api.SecretGetSecretByHashHandler = secret.GetSecretByHashHandlerFunc(func(params secret.GetSecretByHashParams) middleware.Responder {
		foundSecret, prs := secrets[params.Hash]
		if !prs {
			return secret.NewGetSecretByHashNotFound()
		}

		expired := time.Time(foundSecret.ExpiresAt).Sub(time.Now())

		if expired > 0 && foundSecret.RemainingViews > 0 {
			foundSecret.RemainingViews--
			UpdateSecret(&foundSecret)

			return secret.NewGetSecretByHashOK().WithPayload(&foundSecret)
		}

		DeleteSecret(&foundSecret)
		return secret.NewGetSecretByHashNotFound()
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

// AddSecret add a secret to the storage
func AddSecret(secret *models.Secret) {
	secretsLock.Lock()
	defer secretsLock.Unlock()

	secrets[secret.Hash] = *secret
}

// UpdateSecret updates state of the secret
func UpdateSecret(secret *models.Secret) {
	secretsLock.Lock()
	defer secretsLock.Unlock()

	secrets[secret.Hash] = *secret

}

// DeleteSecret deletes the secret from the storage
func DeleteSecret(secret *models.Secret) {
	secretsLock.Lock()
	defer secretsLock.Unlock()

	delete(secrets, secret.Hash)
}
