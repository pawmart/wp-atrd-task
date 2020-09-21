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

// AddSecret adds new Secret to Storage
func (s *InMemoryStorage) AddSecret(secretValue string, viewsLeft uint32, expiresAfterMinutes uint32) Secret {
	var expireDate *time.Time
	now := time.Now()

	if expiresAfterMinutes != 0 {
		t := now.Add(time.Minute * time.Duration(expiresAfterMinutes))
		expireDate = &t
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := uuid.New()
	s.values[id] = Secret{
		Id:             id,
		Value:          secretValue,
		RemainingViews: viewsLeft,
		CreatedAt:      now,
		ExpiresAfter:   expireDate,
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
