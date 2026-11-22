package alerts

import "time"

type Alert struct {
	Name     string
	Message  string
	FiredAt  time.Time
	Resolved bool
}
