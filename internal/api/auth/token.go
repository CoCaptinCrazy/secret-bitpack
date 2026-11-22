package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Token struct {
	Value     string
	Subject   string
	ExpiresAt time.Time
}

func Issue(subject string, ttl time.Duration) (Token, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return Token{}, err
	}
	return Token{Value: hex.EncodeToString(raw), Subject: subject, ExpiresAt: time.Now().Add(ttl)}, nil
}

func (t Token) Valid(now time.Time) bool {
	return t.Value != "" && now.Before(t.ExpiresAt)
}
