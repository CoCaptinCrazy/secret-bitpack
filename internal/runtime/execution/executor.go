package execution

import (
	"context"
	"os/exec"
)

type Executor struct{}

func (Executor) Run(ctx context.Context, command string, args ...string) error {
	return exec.CommandContext(ctx, command, args...).Run()
}
