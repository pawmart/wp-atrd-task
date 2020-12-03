package repository

import (
	"time"

	"github.com/alkmc/wp-atrd-task/internal/entity"
)

//Repository is responsible for DB operation on Secret entity
type Repository interface {
	Set(key string, value *entity.Secret, exp time.Duration) error
	Get(key string) (*entity.Secret, error)
	Expire(key string) error
	CloseDB()
}
