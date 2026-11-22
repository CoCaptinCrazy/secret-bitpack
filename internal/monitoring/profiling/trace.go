package profiling

import "time"

type Span struct {
	Name      string
	StartedAt time.Time
	Duration  time.Duration
}
