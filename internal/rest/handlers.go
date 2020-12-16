package rest

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/systemz/wp-atrd-task/internal/config"
	"github.com/systemz/wp-atrd-task/internal/model"
	"github.com/systemz/wp-atrd-task/internal/service"
	"net/http"
	"time"
)

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
	// FIXME
	model.Redis.Set(redisKey, encryptedByte, time.Second*900)
	w.Write([]byte(newUuid))
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
	resByte, err := service.DecryptWithAesCfb([]byte(config.AES_KEY), []byte(rawResFromRedis))
	if err != nil {
		logrus.Errorf("can't decrypt string: %v", err)
		w.WriteHeader(500)
		return
	}
	logrus.Debugf("secret: %v", string(resByte[:]))
	w.Write(resByte)
}
