package validator

import "github.com/alkmc/wp-atrd-task/internal/entity"

//Validator is responsible for Secret entity validation
type Validator interface {
	FormData(s *entity.FormSecret) (*entity.FormInt, error)
	Hash(hash string) error
}
