package autoscale

type Policy struct {
	Minimum int
	Maximum int
	Target  float64
}

func (p Policy) Desired(current int, utilization float64) int {
	if p.Target <= 0 {
		return current
	}
	desired := int(float64(current)*utilization/p.Target + 0.5)
	if desired < p.Minimum {
		return p.Minimum
	}
	if desired > p.Maximum {
		return p.Maximum
	}
	return desired
}
