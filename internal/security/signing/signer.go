package signing

import (
	"crypto"
	"crypto/rand"
)

func Sign(key crypto.Signer, content []byte, opts crypto.SignerOpts) ([]byte, error) {
	digest := crypto.Hash(opts.HashFunc()).New()
	if _, err := digest.Write(content); err != nil {
		return nil, err
	}
	return key.Sign(rand.Reader, digest.Sum(nil), opts)
}
