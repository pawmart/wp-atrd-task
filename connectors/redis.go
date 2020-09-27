package connectors

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"wp-atrd-task/config"
)

type RedisConnector interface {
	FetchSecret(hash string) ([]byte, error)
	SetSecret(hash string, v interface{}) ([]byte, error)
}

type redisConnector struct {
	*redis.Client
}

func NewRedis(c *config.Config) RedisConnector {
	return redisConnector{redis.NewClient(&redis.Options{
		Addr:     c.Redis.Address,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})}
}

func (r redisConnector) FetchSecret(hash string) ([]byte, error) {
	return r.Get(hash).Bytes()
}

func (r redisConnector) SetSecret(hash string, v interface{}) ([]byte, error) {
	s, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	r.Set(hash, s, 0)
	return r.Get(hash).Bytes()
}
