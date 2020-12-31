package secrets

import (
		"testing"
		"time"
)

func TestNewExpiration(t *testing.T) {
		views := 5
		minutes := time.Duration(5)
		e := NewExpiration(views)

		if e.ViewsLeft != views {
				t.Errorf("number of views doesnt match, got %d expected %d", views, e.ViewsLeft)
		}

		gt := time.Now().Add(minutes * time.Minute)
		time.Sleep(1 * time.Second)
		e.inMinutes(minutes)
		time.Sleep(1 * time.Second)
		lt := time.Now().Add(minutes * time.Minute)

		if e.Date.Unix() < gt.Unix() || e.Date.Unix() > lt.Unix() {
				t.Errorf("expiration date should be between %s and %s, got %s", gt, lt, e.Date)
		}

		e.never()

		if e.Date.Unix() != 0 {
				t.Errorf("expiration date should not be set")
		}
}

func TestNewSecret(t *testing.T) {
		secretText := "psst"
		views := 5
		minutes := 5
		s := NewSecret(secretText, views, 0)
		s2:= NewSecret(secretText, views, minutes)

		if s.Hash == "" {
				t.Errorf("hash should be set")
		}
		if s.Created.Unix() > 0 && s.Created.Unix() < time.Now().Unix() {
				t.Errorf("wrong creation date")
		}
		if s.Secret != secretText {
				t.Errorf("secret text doesnt match, expected %s has %s", secretText, s.Secret)
		}
		if s.Expire.ViewsLeft != views {
				t.Errorf("wrong number of allowed views, expected %d has %d", views, s.Expire.ViewsLeft)
		}
		if s.Expire.Date.Unix() != 0 {
				t.Errorf("expiration date should not be set, expected %d has %d", 0, s.Expire.Date.Unix())
		}
		if s2.Expire.Date.Unix() < time.Now().Unix() {
				t.Errorf("wrong expiration date, expected current date + ~ %d minutes, has %s", 0, s.Expire.Date)
		}
}
