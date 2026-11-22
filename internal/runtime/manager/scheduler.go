package manager

import "time"

type Schedule struct {
	Interval time.Duration
}

func (s Schedule) Next(from time.Time) time.Time { return from.Add(s.Interval) }
