package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alkmc/wp-atrd-task/internal/entity"

	"github.com/go-redis/redis"
)

type redisConn struct {
	client *redis.Client
}

//NewRedis creates new redis Repository
func NewRedis() Repository {
	host, db := getEnvVars()
	opt := &redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	}
	client := redis.NewClient(opt)
	if err := client.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to DB")

	return &redisConn{client}
}

func (r *redisConn) Set(key string, value *entity.Secret, exp time.Duration) error {
	secret, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := r.client.Set(key, secret, exp).Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *redisConn) Get(key string) (*entity.Secret, error) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s := &entity.Secret{}
	if err := json.Unmarshal([]byte(val), s); err != nil {
		log.Println(err)
		return nil, err
	}

	return s, nil
}

func (r *redisConn) Expire(key string) error {
	if err := r.client.Del(key).Err(); err != nil {
		log.Println("err")
		return errors.New("failed to expire secret")
	}
	return nil
}

func (r *redisConn) CloseDB() {
	if err := r.client.Close(); err != nil {
		log.Println("failed to close database connection")
	}
	log.Println("connection to database closed")
}

func getEnvVars() (string, int) {
	host := os.Getenv("redisHOST")
	if host == "" {
		log.Fatal("environment variable redisHOST is required")
	}

	db := os.Getenv("redisDB")
	if db == "" {
		log.Fatal("environment variable redisDB is required")
	}
	dbNum, err := strconv.Atoi(db)
	if err != nil {
		log.Fatal("incorrect database number")
	}

	return host, dbNum
}
