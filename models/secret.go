package models

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"wp-atrd-task/connectors"
	"wp-atrd-task/docs"
)

type Secreter interface {
	CreateSecret(secret string, expireAfterViews string, expireAfter string) (*docs.Secret, error)
	GetSecret(hash string) (*docs.Secret, error)
}

type secretModel struct {
	connector connectors.RedisConnector
}

func NewSecretModel(connector connectors.RedisConnector) Secreter {
	return secretModel{connector: connector}
}

func (s secretModel) CreateSecret(secret string, expireAfterViews string, expireAfter string) (*docs.Secret, error) {
	sec := &docs.Secret{}

	createdAt := time.Now()
	remainingViews, err := strconv.Atoi(expireAfterViews)
	if err != nil {
		return nil, err
	}

	expireAt, err := strconv.Atoi(expireAfter)
	if err != nil {
		return nil, err
	}

	hash := strconv.Itoa(rand.Int()) + secret + expireAfterViews + expireAfter + createdAt.String()
	h := sha512.New().Sum([]byte(hash))

	sec.Hash = fmt.Sprintf("%x", h)
	sec.SecretText = secret
	sec.CreatedAt = createdAt
	sec.ExpiresAt = calculateExpired(createdAt, expireAt)
	sec.RemainingViews = int32(remainingViews)

	sec, err = s.setSecret(sec)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (s secretModel) GetSecret(hash string) (*docs.Secret, error) {
	sec := &docs.Secret{}

	b, err := s.connector.FetchSecret(hash)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&sec); err != nil {
		return nil, err
	}

	if sec.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("secret has expired")
	}

	if sec.RemainingViews == 0 {
		return nil, errors.New("secret is no longer available")
	}

	sec, err = s.updateSecret(sec)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (s secretModel) updateSecret(sec *docs.Secret) (*docs.Secret, error) {
	sec.RemainingViews -= 1
	sec, err := s.setSecret(sec)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (s secretModel) setSecret(sec *docs.Secret) (*docs.Secret, error) {
	a, err := s.connector.SetSecret(sec.Hash, sec)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(bytes.NewReader(a)).Decode(&sec); err != nil {
		return nil, err
	}

	return sec, nil
}

func calculateExpired(createdAt time.Time, expireAfter int) time.Time {
	if expireAfter == 0 {
		return time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return createdAt.Add(time.Duration(expireAfter) * time.Minute)
}
