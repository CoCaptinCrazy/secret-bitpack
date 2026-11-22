package core

import (
	"io"
	"os"

	"secret-bitpack/internal/core/config"
	"secret-bitpack/internal/core/events"
	"secret-bitpack/internal/core/logging"
	"secret-bitpack/internal/core/state"
	"secret-bitpack/internal/runtime/manager"
)

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
		Version: "0.1.17",
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
