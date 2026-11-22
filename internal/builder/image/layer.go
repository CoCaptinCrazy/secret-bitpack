package image

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type projectLayer struct {
	compressed   []byte
	uncompressed []byte
	digest       v1.Hash
	diffID       v1.Hash
}

func newProjectLayer(uncompressed []byte) (v1.Layer, error) {
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(uncompressed); err != nil {
		return nil, fmt.Errorf("compress OCI layer: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("finish compressed OCI layer: %w", err)
	}
	return &projectLayer{
		compressed:   compressed.Bytes(),
		uncompressed: append([]byte(nil), uncompressed...),
		digest:       sha256Hash(compressed.Bytes()),
		diffID:       sha256Hash(uncompressed),
	}, nil
}

func sha256Hash(data []byte) v1.Hash {
	sum := sha256.Sum256(data)
	return v1.Hash{Algorithm: "sha256", Hex: hex.EncodeToString(sum[:])}
}

func (l *projectLayer) Digest() (v1.Hash, error) { return l.digest, nil }
func (l *projectLayer) DiffID() (v1.Hash, error) { return l.diffID, nil }
func (l *projectLayer) Size() (int64, error)     { return int64(len(l.compressed)), nil }
func (l *projectLayer) MediaType() (types.MediaType, error) {
	return types.OCILayer, nil
}
func (l *projectLayer) Compressed() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(l.compressed)), nil
}
func (l *projectLayer) Uncompressed() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(l.uncompressed)), nil
}
