package rest

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/systemz/wp-atrd-task/internal/config"
	"github.com/systemz/wp-atrd-task/internal/model"
	"github.com/systemz/wp-atrd-task/internal/service"
	"net/http"
	"time"
)

type SecretResponse struct {
	Hash           string    `json:"hash"`
	SecretText     string    `json:"secretText"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
	RemainingViews int64     `json:"remainingViews"`
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
	encryptedByte, err := service.EncryptWithAesCfb([]byte(config.AES_KEY), []byte(secretVal))
	if err != nil {
		// FIXME
		logrus.Error(err)
		w.Write([]byte("error, check console"))
		return
	}
	// FIXME save this in DB
	createdAt := time.Now()
	// FIXME get expire time from request
	ttlSeconds := 900
	expiresAt := createdAt.Add(time.Second * time.Duration(ttlSeconds))
	model.Redis.Set(redisKey, encryptedByte, time.Second*time.Duration(ttlSeconds))
	rawResponse := SecretResponse{
		Hash:       newUuid,
		SecretText: secretVal,
		CreatedAt:  createdAt,
		ExpiresAt:  expiresAt,
		// FIXME implement max views
		RemainingViews: 0,
	}
	//response, err := json.Marshal(&rawResponse)
	response, err := json.MarshalIndent(rawResponse, "", "    ")
	if err != nil {
		logrus.Error(err)
		w.Write([]byte("error, check console"))
		return
	}
	w.Write(response)
}

func GetSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	logrus.Debugf("hash: %v", params["hash"])
	rawResFromRedis, err := model.Redis.Get(redisKey(params["hash"])).Result()
	if err != nil {
		logrus.Errorf("can't find hash: %v", err)
		w.WriteHeader(404)
		return
	}
	secretByte, err := service.DecryptWithAesCfb([]byte(config.AES_KEY), []byte(rawResFromRedis))
	if err != nil {
		logrus.Errorf("can't decrypt string: %v", err)
		w.WriteHeader(500)
		return
	}
	logrus.Debugf("secret: %v", string(secretByte[:]))

	rawResponse := SecretResponse{
		Hash:       params["hash"],
		SecretText: string(secretByte[:]),
		CreatedAt:  time.Now(), //FIXME
		ExpiresAt:  time.Now(), //FIXME
		// FIXME implement max views
		RemainingViews: 0,
	}
	//response, err := json.Marshal(&rawResponse)
	response, err := json.MarshalIndent(rawResponse, "", "    ")
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
