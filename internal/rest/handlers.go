package rest

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/systemz/wp-atrd-task/internal/config"
	"github.com/systemz/wp-atrd-task/internal/model"
	"github.com/systemz/wp-atrd-task/internal/service"
	"github.com/thanhpk/randstr"
	"net/http"
	"time"
)

const redisKeyForSecret = "s."

func NewSecret(w http.ResponseWriter, r *http.Request) {
	newUuid := uuid.New()
	redisKey := config.REDIS_KEY_PREFIX + redisKeyForSecret + newUuid.String()
	// FIXME
	err, encryptedString := service.EncryptWithAes128(randstr.String(32))
	if err != nil {
		// FIXME
		logrus.Error(err)
		w.Write([]byte("error, check console"))
		return
	}
	model.Redis.Set(redisKey, encryptedString, time.Second*60)
	w.Write([]byte(newUuid.String()))
}
