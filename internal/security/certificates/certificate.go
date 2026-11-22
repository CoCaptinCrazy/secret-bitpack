package certificates

import (
	"crypto/x509"
	"time"
)

type Certificate struct {
	Parsed *x509.Certificate
}

func (c Certificate) ValidAt(at time.Time) bool {
	return c.Parsed != nil && at.After(c.Parsed.NotBefore) && at.Before(c.Parsed.NotAfter)
}
