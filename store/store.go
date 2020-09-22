package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store provides persistence functionality for Secret records.
type Store struct {
	DB *sqlx.DB
}

// New returns initialized Store object.
func New(ctx context.Context, db *sqlx.DB) (*Store, error) {
	s := &Store{
		DB: db,
	}
	if err := s.InitSecretsTable(ctx); err != nil {
		return nil, fmt.Errorf("failed to init secrets table: %w", err)
	}

	s.Seed(ctx)

	return s, nil
}

// Storer interface for persisting secrets.
type Storer interface {
	GetByID(ctx context.Context, id string) (Secret, error)
	Insert(ctx context.Context, s Secret) (Secret, error)
	Update(ctx context.Context, id string, sc Secret) (Secret, error)
	Delete(ctx context.Context, id string) error
}

// Secret represents secret entry in DB.
type Secret struct {
	ID             uuid.UUID  `db:"id"`
	Content        string     `db:"content"`
	RemainingViews int32      `db:"remaining_views"`
	CreatedAt      time.Time  `db:"created_at"`
	ExpiresAt      *time.Time `db:"expires_at"`
}

// InitSecretsTable creates secrets table if not exists.
func (s *Store) InitSecretsTable(ctx context.Context) error {
	if _, err := s.DB.ExecContext(ctx, queryInitSecretsTable); err != nil {
		return fmt.Errorf("cannot create 'secrets' table: %w", err)
	}
	return nil
}

// Seed populates DB with test data.
func (s *Store) Seed(ctx context.Context) {
	s.DB.MustExec(querySeedTestData)
}

// GetByID gets secret record by ID.
func (s *Store) GetByID(ctx context.Context, id string) (Secret, error) {
	sc := Secret{}
	if err := s.DB.GetContext(ctx, &sc, queryGetSecretByID, id); err != nil {
		return Secret{}, fmt.Errorf("cannot get secret from DB: %w", err)
	}
	return sc, nil
}

// Insert inserts new secret record to DB.
func (s *Store) Insert(ctx context.Context, sc Secret) (Secret, error) {
	sc.ID = uuid.New()
	_, err := s.DB.NamedQueryContext(ctx, queryInsertSecret, sc)
	if err != nil {
		return Secret{}, fmt.Errorf("cannot insert new secret into DB: %w", err)
	}
	return sc, nil
}

// Update updates secret record (Content and RemainingViews fields) in DB.
func (s *Store) Update(ctx context.Context, id string, new Secret) (Secret, error) {
	old := Secret{}
	if err := s.DB.GetContext(ctx, &old, queryGetSecretByID, id); err != nil {
		return Secret{}, fmt.Errorf("cannot get secret from DB: %w", err)
	}

	old.ID = uuid.MustParse(id)
	old.Content = new.Content
	old.RemainingViews = new.RemainingViews
	_, err := s.DB.NamedQueryContext(ctx, queryUpdateSecret, old)
	if err != nil {
		return Secret{}, fmt.Errorf("cannot update secret into DB: %w", err)
	}
	return old, nil
}

// Delete deletes an secret record from DB.
func (s *Store) Delete(ctx context.Context, id string) error {
	if _, err := s.DB.ExecContext(ctx, queryDeleteSecretByID, id); err != nil {
		return fmt.Errorf("cannot delete secret from DB: %w", err)
	}
	return nil
}
