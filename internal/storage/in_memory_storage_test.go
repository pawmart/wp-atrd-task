package storage

import (
	"github.com/google/uuid"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestInMemoryStorage_GetSecret(t *testing.T) {
	firstId := uuid.New()
	value := "test"
	now := time.Now()
	expiresAt := now.Add(time.Minute * 10)
	remainingViews := uint32(5)

	secondId := uuid.New()
	secondExpiresAt := now.Add(- time.Minute * 6)

	thirdId := uuid.New()

	m := map[uuid.UUID]Secret{
		firstId: {
			Value:          value,
			CreatedAt:      now,
			ExpiresAfter:   &expiresAt,
			RemainingViews: remainingViews,
		},
		secondId: {
			Value:          value,
			CreatedAt:      now,
			ExpiresAfter:   &secondExpiresAt,
			RemainingViews: remainingViews,
		},
		thirdId: {
			Value:          value,
			CreatedAt:      now.Add(- time.Minute * 5),
			ExpiresAfter:   nil,
			RemainingViews: remainingViews,
		},
	}

	s := InMemoryStorage{
		mutex:  &sync.Mutex{},
		values: m,
	}

	if secret, exist := s.GetSecret(firstId); secret == nil || exist == false {
		t.Errorf("Secret should exist")
	}

	if secret, exist := s.GetSecret(secondId); secret != nil || exist == true {
		t.Errorf("Secret shouldnt exist")
	}

	if secret, exist := s.GetSecret(thirdId); secret == nil || exist != true {
		t.Errorf("Secret should exist")
	}

	if secret, exist := s.GetSecret(uuid.New()); secret != nil || exist == true {
		t.Errorf("Secret shouldnt exist")
	}

	if secret, exist := s.GetSecret(firstId); exist != true || secret.RemainingViews != remainingViews-2 {
		t.Errorf(
			"Secret remaining views after two fetches should be equal: %v, got %v",
			secret.RemainingViews,
			remainingViews-2,
		)
	}

}

func TestInMemoryStorage_AddSecret(t *testing.T) {
	s := InMemoryStorage{
		mutex:  &sync.Mutex{},
		values: make(map[uuid.UUID]Secret),
	}

	value := "test"
	createdAt := time.Date(2020, time.January, 1, 1, 1, 1, 1, time.UTC)
	expiresAt := time.Date(2020, time.January, 1, 1, 5, 1, 1, time.UTC)
	remainingViews := uint32(5)

	secret := s.AddSecret(value, remainingViews, createdAt, &expiresAt)

	testSecret := Secret{
		Id:             secret.Id,
		Value:          value,
		CreatedAt:      createdAt,
		ExpiresAfter:   &expiresAt,
		RemainingViews: remainingViews,
	}

	if !reflect.DeepEqual(secret, testSecret) {
		t.Errorf("AddSecret() = %v, want %v", secret, testSecret)
	}

	value = "test2"
	createdAt = time.Date(2020, time.January, 1, 1, 1, 1, 1, time.UTC)
	remainingViews = uint32(1)

	secret = s.AddSecret(value, remainingViews, createdAt, nil)

	testSecret = Secret{
		Id:             secret.Id,
		Value:          value,
		CreatedAt:      createdAt,
		ExpiresAfter:   nil,
		RemainingViews: remainingViews,
	}

	if !reflect.DeepEqual(secret, testSecret) {
		t.Errorf("AddSecret() = %v, want %v", secret, testSecret)
	}
}
