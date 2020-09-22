package server

import "github.com/seblw/wp-atrd-task/store"

func convToStore(sc Secret) store.Secret {
	return store.Secret{
		// ID is created in Store layer.
		Content:        sc.SecretText,
		RemainingViews: sc.RemainingViews,
		CreatedAt:      sc.CreatedAt,
		ExpiresAt:      sc.ExpiresAt,
	}
}

func convToAPI(sc store.Secret) Secret {
	return Secret{
		Hash:           sc.ID.String(),
		SecretText:     sc.Content,
		RemainingViews: sc.RemainingViews,
		CreatedAt:      sc.CreatedAt,
		ExpiresAt:      sc.ExpiresAt,
	}
}
