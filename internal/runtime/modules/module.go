package modules

type Module struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	Version         string `json:"version,omitempty"`
	Blueprint       string `json:"blueprint,omitempty"`
	Runtime         string `json:"runtime,omitempty"`
	ProcessID       string `json:"process_id,omitempty"`
	PID             int    `json:"pid,omitempty"`
	Started         int64  `json:"started,omitempty"`
	DesiredReplicas int    `json:"desired_replicas,omitempty"`
	AutoscaleMin    int    `json:"autoscale_min,omitempty"`
	AutoscaleMax    int    `json:"autoscale_max,omitempty"`
	AutoscaleCPU    int    `json:"autoscale_cpu,omitempty"`
}
