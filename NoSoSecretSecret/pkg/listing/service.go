package listing

// Repository provides access to secrets
type Repository interface {
	GetSecret(string) (*Secret, error)
}

// Service provides operation(s) for listing secrets
type Service interface {
	GetSecret(string) (*Secret, error)
}

// NewService creates listing services with aproperate dependences
func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

func (s *service) GetSecret(secretHash string) (*Secret, error) {
	return s.repo.GetSecret(secretHash)
}
