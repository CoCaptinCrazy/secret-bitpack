package layers

import "crypto/sha256"

type Layer struct {
	Digest [32]byte
	Size   int64
}

func New(content []byte) Layer {
	return Layer{Digest: sha256.Sum256(content), Size: int64(len(content))}
}
