package adding

import "nososecretsercet/pkg/listing"

// Service provides a Secret adding operations
type Service interface {
	AddSecret(Secret) (*listing.Secret, error)
}

// Repository provides access to secret repository
type Repository interface {
	AddSecret(Secret) (string, error)
	GetSecret(string) (*listing.Secret, error) 
}

// NewService creates and adding service with neeede dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

// AddSecret persists given secret to storage
func (s *service) AddSecret(secret Secret) (*listing.Secret, error) {

	newHash, err := s.repo.AddSecret(secret)

	if err != nil {
		return nil, err
	}

	newCreatedSecret, err := s.repo.GetSecret(newHash)

	if err != nil {
		return nil, err
	}

	return newCreatedSecret, nil
}
