package profiling

import "time"

type Report struct {
	Spans []Span
}

func (r Report) TotalTime() time.Duration {
	var total time.Duration
	for _, span := range r.Spans {
		total += span.Duration
	}
	return total
}
