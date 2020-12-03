package service

import (
	"time"

	"github.com/alkmc/wp-atrd-task/internal/entity"
	"github.com/alkmc/wp-atrd-task/internal/repository"

	"github.com/google/uuid"
)

type secretService struct {
	repo repository.Repository
}

//NewService returns new Secret Service
func NewService(r repository.Repository) Service {
	return &secretService{repo: r}
}

func (s *secretService) Create(p *entity.Secret) error {
	p.Hash = uuid.New().String()
	p.CreatedAt = time.Now()
	p.CalculateExpiration()

	return s.repo.Set(p.Hash, p, p.CastToDuration())
}

func (s *secretService) FindAndUpdate(hash string) (*entity.Secret, error) {
	sec, err := s.repo.Get(hash)
	if err != nil {
		return nil, err
	}

	decreaseViews(sec)
	if sec.RemainingViews == 0 {
		if err := s.repo.Expire(hash); err != nil {
			return nil, err
		}
		return sec, nil
	}

	if err := s.repo.Set(sec.Hash, sec, sec.NewExpirationAt()); err != nil {
		return nil, err
	}

	return sec, nil
}

func decreaseViews(s *entity.Secret) {
	s.RemainingViews--
}
