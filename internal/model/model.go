package model

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/systemz/wp-atrd-task/internal/config"
)

var (
	Redis *redis.Client
)

func RedisInit() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_HOST,
		Password: config.REDIS_PASSWORD,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Panic("ping to redis failed")
	}
	logrus.Info("redis connected")
}
