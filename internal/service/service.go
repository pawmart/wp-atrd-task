package service

import "github.com/alkmc/wp-atrd-task/internal/entity"

//Service is responsible for interaction with Repository interface
type Service interface {
	Create(p *entity.Secret) error
	FindAndUpdate(id string) (*entity.Secret, error)
}
