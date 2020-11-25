package db

import (
	"database/sql"
	"fmt"

	"github.com/maciejem/secret/pkg/apperrors"
	"github.com/maciejem/secret/pkg/model"
)

func (db Database) CreateSecret(secret model.Secret) error {
	var err error
	if secret.ExpiresAt != nil {
		query := `INSERT INTO secrets (id, secret_text, remaining_views, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)`
		err = db.Conn.QueryRow(query, secret.Hash, secret.SecretText, secret.RemainingViews, secret.CreatedAt, *secret.ExpiresAt).Scan()
	} else {
		query := `INSERT INTO secrets (id, secret_text, remaining_views, created_at) VALUES ($1, $2, $3, $4)`
		err = db.Conn.QueryRow(query, secret.Hash, secret.SecretText, secret.RemainingViews, secret.CreatedAt).Scan()
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("creating secret: %v", err)
	}
	return nil
}

func (db Database) GetSecretById(secretId string) (model.Secret, error) {
	secret := model.Secret{}
	query := `SELECT * FROM secrets WHERE id = $1;`
	row := db.Conn.QueryRow(query, secretId)
	err := row.Scan(&secret.Hash, &secret.SecretText, &secret.RemainingViews, &secret.CreatedAt, &secret.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return secret, apperrors.ErrNoMatch
		}
		return secret, fmt.Errorf("getting secret by id: %v", err)
	}
	return secret, nil
}

func (db Database) DecreaseSecretRemainingViewsById(secretId string) error {
	query := `UPDATE secrets SET remaining_views=remaining_views-1 WHERE id=$1;`
	err := db.Conn.QueryRow(query, secretId).Scan()
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("decreasing secret remaining views by id: %v", err)
	}
	return nil
}
