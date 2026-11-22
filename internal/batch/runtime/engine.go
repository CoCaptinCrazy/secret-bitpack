package runtime

import "context"

type Step func(context.Context) error

type Engine struct{ Steps []Step }

func (e Engine) Run(ctx context.Context) error {
	for _, step := range e.Steps {
		if err := step(ctx); err != nil {
			return err
		}
	}
	return nil
}
