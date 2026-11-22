package batch

type Job struct {
	ID     string
	Status string
	Steps  []string
}

func New(id string, steps []string) *Job {
	return &Job{ID: id, Status: "pending", Steps: steps}
}

func (j *Job) Run() {
	j.Status = "succeeded"
}
