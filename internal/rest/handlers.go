package rest

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/systemz/wp-atrd-task/internal/config"
	"github.com/systemz/wp-atrd-task/internal/model"
	"github.com/systemz/wp-atrd-task/internal/service"
	"net/http"
	"strconv"
	"time"
)

type SecretResponse struct {
	Hash           string    `json:"hash"`
	SecretText     string    `json:"secretText"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
	RemainingViews int       `json:"remainingViews"`
}

func redisKey(str string) string {
	return config.REDIS_KEY_PREFIX + "s." + str
}

func NewSecret(w http.ResponseWriter, r *http.Request) {
	newUuid := uuid.New().String()
	redisKey := redisKey(newUuid)

	secretVal := r.FormValue("secret")
	// deny creating too short secrets
	if len([]rune(secretVal)) < 1 {
		// FIXME
		w.Write([]byte("too short secret"))
		return
	}

	// detect how much time secret should be valid
	expireAfterMinutesStr := r.FormValue("expireAfter")
	expireAfterMinutes, err := strconv.Atoi(expireAfterMinutesStr)
	if err != nil {
		// fallback to lack of time expiration
		expireAfterMinutes = 0
	}

	// detect how many views are allowed
	expireAfterViewsStr := r.FormValue("expireAfterViews")
	expireAfterViews, err := strconv.Atoi(expireAfterViewsStr)
	if err != nil {
		// fallback to unlimited fetches from DB
		expireAfterViews = 0
	}
	if expireAfterViews < 1 {
		// FIXME
		w.Write([]byte("expireAfterViews must be greater than 0"))
		return
	}

	encryptedByte, err := service.EncryptWithAesCfb([]byte(config.AES_KEY), []byte(secretVal))
	if err != nil {
		// FIXME
		logrus.Error(err)
		w.Write([]byte("error, check console"))
		return
	}
	createdAt := time.Now()
	ttl := time.Minute * time.Duration(expireAfterMinutes)
	expiresAt := createdAt.Add(ttl)
	b64Secret := base64.StdEncoding.EncodeToString(encryptedByte)

	// prepare shared json for user and DB
	rawResult := SecretResponse{
		Hash:           newUuid,
		SecretText:     b64Secret,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
		RemainingViews: expireAfterViews,
	}

	// create JSON for DB with encrypted secret
	result, err := json.Marshal(&rawResult)

	// throw JSON to DB
	model.Redis.Set(redisKey, result, ttl)

	// create JSON for API result with plaintext secret
	rawResult.SecretText = secretVal
	result, err = json.MarshalIndent(rawResult, "", "    ")
	if err != nil {
		logrus.Error(err)
		w.Write([]byte("error, check console"))
		return
	}

	// show result to user
	w.Write(result)
}

func GetSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	logrus.Debugf("hash: %v", params["hash"])
	key := redisKey(params["hash"])
	rawResFromRedis, err := model.Redis.Get(key).Result()
	if err != nil {
		logrus.Errorf("can't find hash: %v", err)
		w.WriteHeader(404)
		return
	}
	var jsonFromDb SecretResponse
	err = json.Unmarshal([]byte(rawResFromRedis), &jsonFromDb)
	secretB64Byte, err := base64.StdEncoding.DecodeString(string(jsonFromDb.SecretText[:]))
	// decrypt secret for user
	secretByte, err := service.DecryptWithAesCfb([]byte(config.AES_KEY), []byte(secretB64Byte))
	if err != nil {
		logrus.Errorf("can't decrypt string: %v", err)
		w.WriteHeader(500)
		return
	}
	secretPlaintext := string(secretByte[:])
	//logrus.Debugf("secret: %v", string(secretByte[:]))
	logrus.Debugf("secret: %v", secretPlaintext)

	rawResponse := jsonFromDb
	rawResponse.SecretText = secretPlaintext

	// view counter depleted actions
	if rawResponse.RemainingViews == 1 {
		// this was last view, remove record from DB and show result instantly
		model.Redis.Del(key)
		// show user that no further views are allowed
		rawResponse.RemainingViews--
		// pretty formatted JSON result
		response, _ := json.MarshalIndent(rawResponse, "", "    ")
		w.Write(response)
		return
	}

	// update counters for DB record and user facing JSON
	rawResponse.RemainingViews--
	jsonFromDb.RemainingViews--

	// view counter still allowing more views than 1
	//if rawResponse.RemainingViews > 1 {
	// update JSON for DB with encrypted secret
	result, err := json.Marshal(&jsonFromDb)
	if err != nil {
		logrus.Errorf("can't marshal: %v", err)
		w.WriteHeader(500)
	}

	// update JSON in DB
	if time.Now().After(jsonFromDb.ExpiresAt) {
		// time expired for this record already, redis probably didn't run GC yet
		// delete it instantly and don't show secret to user, pretend that secret doesn't exist in DB already
		model.Redis.Del(key)
		w.WriteHeader(404)
		return
	}
	// this redis lib probably doesn't have option to leave TTL as is, calculate it and set again
	ttl := jsonFromDb.ExpiresAt.Sub(time.Now())
	model.Redis.Set(key, result, ttl)

	// respond with updated view counter
	response, _ := json.MarshalIndent(rawResponse, "", "    ")
	w.Write(response)
}

func GetDocsRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/v1/docs/", 302)
}

func GetSwaggerUi(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.HTTP_DOCS_DIR+"swagger.html")
}

func GetSwaggerYml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.HTTP_DOCS_DIR+"swagger.yml")
}
