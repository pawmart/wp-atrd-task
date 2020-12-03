package validator

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/alkmc/wp-atrd-task/internal/entity"

	"github.com/google/uuid"
)

type secretValidator struct {
}

//NewValidator returns new Secret Validator
func NewValidator() Validator {
	return &secretValidator{}
}

func (v *secretValidator) Hash(secretHash string) error {
	if _, err := uuid.Parse(secretHash); err != nil {
		return err
	}
	return nil
}

func (v *secretValidator) FormData(s *entity.FormSecret) (*entity.FormInt, error) {
	trimmedSecret := strings.TrimSpace(s.Secret)
	if len(trimmedSecret) == 0 {
		return nil, errors.New("secret shall be provided")
	}

	if s.ExpireAfter == "" {
		return nil, errors.New("expireAfter must by provided")
	}
	ea, err := strconv.ParseInt(s.ExpireAfter, 10, 32)
	if err != nil {
		log.Println(err)
		return nil, errors.New("expireAfter shall be integer")
	}
	if ea < 0 {
		return nil, errors.New("expireAfter must be greater or equal 0")
	}

	if s.ExpireAfterViews == "" {
		return nil, errors.New("expireAfterViews shall be provided")
	}
	eav, err := strconv.ParseInt(s.ExpireAfterViews, 10, 32)
	if err != nil {
		return nil, errors.New("expireAfter shall be integer")
	}
	if eav < 1 {
		return nil, errors.New("expireAfterViews must be greater than 0")
	}

	fi := &entity.FormInt{
		ExpireAfter:      int32(ea),
		ExpireAfterViews: int32(eav),
	}

	return fi, nil
}
