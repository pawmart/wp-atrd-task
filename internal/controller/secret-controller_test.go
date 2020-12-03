package controller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alkmc/wp-atrd-task/internal/entity"
	"github.com/alkmc/wp-atrd-task/internal/repository"
	"github.com/alkmc/wp-atrd-task/internal/responder"
	"github.com/alkmc/wp-atrd-task/internal/service"
	"github.com/alkmc/wp-atrd-task/internal/validator"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	sRepo       = repository.NewRedis()
	sService    = service.NewService(sRepo)
	sValidator  = validator.NewValidator()
	sController = NewController(sService, sValidator)
)

const (
	exampleSecret       = "asdfasdfasdfasdf"
	startViews    int32 = 5
	remainViews   int32 = 4
)

func TestGetProductByID(t *testing.T) {
	const path = "/v1/secret/%v"

	// insert new product
	hash := uuid.New().String()
	setupWithHash(hash)

	// create a http GET request
	req := httptest.NewRequest("GET", fmt.Sprintf(path, hash), nil)

	// record http Response
	resp := httptest.NewRecorder()

	// assign http Handler function
	r := chi.NewRouter()
	r.Get("/v1/secret/{hash}", sController.GetSecretByHash)

	// dispatch the http request
	r.ServeHTTP(resp, req)

	// assert http status code
	checkResponseCode(t, http.StatusOK, resp.Code)

	// decode the http response
	var s entity.Secret
	if err := json.NewDecoder(io.Reader(resp.Body)).Decode(&s); err != nil {
		log.Fatal(err)
	}

	// assert http response
	assert.Equal(t, hash, s.Hash)
	assert.Equal(t, exampleSecret, s.SecretText)
	assert.Equal(t, remainViews, s.RemainingViews)
	assert.NotNil(t, s.CreatedAt)
	assert.NotNil(t, s.ExpiresAt)

	// clean up database
	tearDown(s.Hash)
}

func TestGetNotExistingSecret(t *testing.T) {
	const (
		errMsg = "Secret not found"
		path   = "/v1/secret/%v"
	)

	// create new hash
	hash := uuid.New().String()

	// create a http GET request
	req := httptest.NewRequest("GET", fmt.Sprintf(path, hash), nil)

	// record http Response
	resp := httptest.NewRecorder()

	// assign http Handler function
	r := chi.NewRouter()
	r.Get("/v1/secret/{hash}", sController.GetSecretByHash)

	// dispatch the http request
	r.ServeHTTP(resp, req)

	// assert http status code
	checkResponseCode(t, http.StatusNotFound, resp.Code)

	// decode the http response
	var se responder.SecretError
	if err := json.NewDecoder(io.Reader(resp.Body)).Decode(&se); err != nil {
		log.Fatal(err)
	}

	// assert http response
	assert.Equal(t, errMsg, se.Description)
}

func TestGetSecretWithIncorrectHash(t *testing.T) {
	const (
		errMsg        = "invalid hash"
		incorrectHash = "incorrect-secret-hash-correct-length"
		path          = "/v1/secret/%v"
	)

	// create a http GET request
	req := httptest.NewRequest("GET", fmt.Sprintf(path, incorrectHash), nil)

	// record http Response
	resp := httptest.NewRecorder()

	// assign http Handler function
	r := chi.NewRouter()
	r.Get("/v1/secret/{hash}", sController.GetSecretByHash)

	// dispatch the http request
	r.ServeHTTP(resp, req)

	// assert http status code
	checkResponseCode(t, http.StatusBadRequest, resp.Code)

	// decode the http response
	var se responder.SecretError
	if err := json.NewDecoder(io.Reader(resp.Body)).Decode(&se); err != nil {
		log.Fatal(err)
	}

	// assert http response
	assert.Equal(t, errMsg, se.Description)
}

func TestAddSecret(t *testing.T) {
	// insert new product
	body := strings.NewReader("secret=asdfasdfasdfasdf&expireAfterViews=5&expireAfter=60")

	// create a http GET request
	req := httptest.NewRequest("POST", "/v1/secret", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// record http Response
	resp := httptest.NewRecorder()

	// assign http Handler function
	r := chi.NewRouter()
	r.Post("/v1/secret", sController.AddSecret)

	// dispatch the http request
	r.ServeHTTP(resp, req)

	// assert http status code
	checkResponseCode(t, http.StatusOK, resp.Code)

	// decode the http response
	var s entity.Secret
	if err := json.NewDecoder(io.Reader(resp.Body)).Decode(&s); err != nil {
		log.Fatal(err)
	}

	// assert http response
	assert.Equal(t, exampleSecret, s.SecretText)
	assert.Equal(t, startViews, s.RemainingViews)
	assert.NotNil(t, s.Hash)
	assert.NotNil(t, s.CreatedAt)
	assert.NotNil(t, s.ExpiresAt)

	// clean up database
	tearDown(s.Hash)
}

func TestAddSecretWithoutExpiration(t *testing.T) {
	// insert new product
	body := strings.NewReader("secret=asdfasdfasdfasdf&expireAfterViews=5&expireAfter=0")

	// create a http GET request
	req := httptest.NewRequest("POST", "/v1/secret", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/xml")

	// record http Response
	resp := httptest.NewRecorder()

	// assign http Handler function
	r := chi.NewRouter()
	r.Post("/v1/secret", sController.AddSecret)

	// dispatch the http request
	r.ServeHTTP(resp, req)

	// assert http status code
	checkResponseCode(t, http.StatusOK, resp.Code)

	// decode the http response
	var s entity.Secret
	if err := xml.NewDecoder(io.Reader(resp.Body)).Decode(&s); err != nil {
		log.Fatal(err)
	}

	// assert http response
	assert.Equal(t, exampleSecret, s.SecretText)
	assert.Equal(t, startViews, s.RemainingViews)
	assert.NotNil(t, s.Hash)
	assert.NotNil(t, s.CreatedAt)
	assert.Nil(t, s.ExpiresAt)

	// clean up database
	tearDown(s.Hash)
}

func TestAddSecretIncorrectForm(t *testing.T) {
	const errMsg = "expireAfterViews must be greater than 0"

	// insert new product
	body := strings.NewReader("secret=asdfasdfasdfasdf&expireAfterViews=0&expireAfter=60")

	// create a http GET request
	req := httptest.NewRequest("POST", "/v1/secret", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// record http Response
	resp := httptest.NewRecorder()

	// assign http Handler function
	r := chi.NewRouter()
	r.Post("/v1/secret", sController.AddSecret)

	// dispatch the http request
	r.ServeHTTP(resp, req)

	// assert http status code
	checkResponseCode(t, http.StatusBadRequest, resp.Code)

	// decode the http response
	var se responder.SecretError
	if err := json.NewDecoder(io.Reader(resp.Body)).Decode(&se); err != nil {
		log.Fatal(err)
	}

	// assert http response
	assert.Equal(t, errMsg, se.Description)
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func setupWithHash(hash string) {
	s := entity.Secret{
		Hash:           hash,
		SecretText:     exampleSecret,
		CreatedAt:      time.Now(),
		ExpireAfter:    1,
		RemainingViews: startViews,
	}
	s.CalculateExpiration()
	addSecret(s)
}

func addSecret(s entity.Secret) {
	if err := sRepo.Set(s.Hash, &s, s.CastToDuration()); err != nil {
		log.Fatal(err)
	}
}

func tearDown(hash string) {
	if err := sRepo.Expire(hash); err != nil {
		log.Fatal(err)
	}
}
