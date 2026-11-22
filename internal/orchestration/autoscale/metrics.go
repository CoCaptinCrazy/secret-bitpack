package autoscale

type Sample struct {
	CPU    float64
	Memory float64
}

func (s Sample) Average() float64 { return (s.CPU + s.Memory) / 2 }
