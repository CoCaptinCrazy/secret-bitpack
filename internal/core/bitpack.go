package core

import (
	"io"
	"os"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/core/config"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/core/events"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/core/logging"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/core/state"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/runtime/manager"
)

var Version = "0.1.17"

type Bitpack struct {
	Config  *config.Config
	Logger  *logging.Logger
	State   *state.Manager
	Events  *events.Bus
	Runtime *manager.Manager
	Output  io.Writer
	Version string
}

func New(output io.Writer) *Bitpack {
	return &Bitpack{
		Config:  config.Default(),
		Logger:  logging.New(output),
		State:   state.New(),
		Events:  events.New(),
		Runtime: manager.New(),
		Output:  output,
		Version: Version,
	}
}

func (b *Bitpack) Initialize() error {
	if err := b.State.Load(); err != nil {
		return err
	}
	configured, err := config.Load(".bitpack/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	b.Config = configured
	return nil
}
