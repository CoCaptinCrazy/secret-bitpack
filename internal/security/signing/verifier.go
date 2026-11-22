package signing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

func Verify(key *rsa.PublicKey, content, signature []byte) error {
	digest := sha256.Sum256(content)
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature)
}
