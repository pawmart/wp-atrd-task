package listing

// Secret defines properties of secret to be listed
type Secret struct {
	Hash           string `json:"hash"`
	SecretText     string `json:"secretText"`
	CreatedAt      string `json:"createdAt"`
	ExpiresAt      string `json:"expiresAt"`
	RemainingViews int32  `json:"remainingViews"`
}
