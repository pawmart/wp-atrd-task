package listing

import "time"

// Secret defines properties of secret to be listed
type Secret struct {
	Hash           string    `xml:"hash,attr" json:"hash"`
	SecretText     string    `xml:"secretText,atr" json:"secretText"`
	CreatedAt      time.Time `xml:"createdAt,atr" json:"createdAt"`
	ExpiresAt      time.Time `xml:"expiresAt,atr" json:"expiresAt"`
	RemainingViews int32     `xml:"remainingViews,atr" json:"remainingViews"`
}
