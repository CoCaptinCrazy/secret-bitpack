package events

import "context"

type Dispatcher struct{ Bus *Bus }

func (d Dispatcher) Publish(ctx context.Context, event Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	d.Bus.Publish(event)
	return nil
}
