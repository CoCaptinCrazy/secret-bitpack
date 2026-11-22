package deployment

type Rollout struct {
	DeploymentID string
	Desired      int
	Updated      int
}

func (r Rollout) Complete() bool { return r.Desired == r.Updated }
