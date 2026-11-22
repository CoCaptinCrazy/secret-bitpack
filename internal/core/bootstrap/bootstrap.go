package bootstrap

import (
	"io"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/core"
)

type Bootstrap struct{}

func Initialize(output io.Writer) (*core.Bitpack, error) {
	bitpack := core.New(output)
	if err := bitpack.Initialize(); err != nil {
		return nil, err
	}
	return bitpack, nil
}
