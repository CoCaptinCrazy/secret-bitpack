package deployment

import "time"

type Deployment struct {
	ID        string
	Blueprint string
	Replicas  int
	CreatedAt time.Time
}
