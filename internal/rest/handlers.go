package rest

import (
	"github.com/google/uuid"
	"github.com/systemz/wp-atrd-task/internal/config"
	"github.com/systemz/wp-atrd-task/internal/model"
	"net/http"
	"time"
)

const redisKeyForSecret = "s."

func NewSecret(w http.ResponseWriter, r *http.Request) {
	newUuid := uuid.New()
	redisKey := config.REDIS_KEY_PREFIX + redisKeyForSecret + newUuid.String()
	// FIXME
	model.Redis.Set(redisKey, "CHANGEME", time.Second*60)
	w.Write([]byte(newUuid.String()))
}
