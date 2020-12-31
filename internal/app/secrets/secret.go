package secrets

import (
		"github.com/google/uuid"
		"time"
)

type Hash string

type Secret struct {
		Hash    Hash
		Secret  string
		Expire  Expiration
		Created time.Time
}

func NewSecret(secret string, views int, minutes int) Secret {
		expire := NewExpiration(views)
		if minutes > 0 {
				expire.inMinutes(time.Duration(minutes))
		}else {
				expire.never()
		}

		return Secret{
				Hash:    NewHash(),
				Secret:  secret,
				Expire:  expire,
				Created: time.Now(),
		}
}

func NewHash() Hash {
		uid := uuid.New()
		return Hash(uid.String())
}

type Expiration struct {
		ViewsLeft int
		Date      time.Time
}

func NewExpiration(views int) Expiration {
		return Expiration{
				ViewsLeft: views,
				Date:      time.Now(),
		}
}

func (e *Expiration) inMinutes(t time.Duration) {
		e.Date = time.Now().Add(t * time.Minute)
}

func (e *Expiration) never() {
		e.Date = time.Unix(int64(0), 0)
}

type SecretRepository interface {
		Find(hash string) (Secret, error)
		Store(secret *Secret) error
}
