package storage

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

// InMemoryStorage implementation struct
type InMemoryStorage struct {
	mutex  *sync.Mutex
	values map[uuid.UUID]Secret
}

// NewInMemoryStorage returns new InMemoryStorage
func NewInMemoryStorage() Storage {
	return &InMemoryStorage{
		mutex:  &sync.Mutex{},
		values: make(map[uuid.UUID]Secret),
	}
}

func NewInMemoryStorageWithData(data map[uuid.UUID]Secret) Storage {
	return &InMemoryStorage{
		mutex:  &sync.Mutex{},
		values: data,
	}
}

// AddSecret adds new Secret to InMemoryStorage
func (s *InMemoryStorage) AddSecret(
	secretValue string,
	viewsLeft uint32,
	createdAt time.Time,
	expiresAt *time.Time,
) Secret {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := uuid.New()
	s.values[id] = Secret{
		Id:             id,
		Value:          secretValue,
		RemainingViews: viewsLeft,
		CreatedAt:      createdAt,
		ExpiresAfter:   expiresAt,
	}

	return s.values[id]
}

// GetSecret fetch, decrement remaining views and returns Secret with bool for existence of secret
func (s *InMemoryStorage) GetSecret(id uuid.UUID) (*Secret, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	secret, exist := s.values[id]

	if !exist {
		return nil, false
	}

	if !secret.isFetchable() {
		delete(s.values, id)
		return nil, false
	}

	secret.RemainingViews--
	s.values[id] = secret

	return &secret, true
}
