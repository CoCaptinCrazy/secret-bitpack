package analytics

type Trend struct {
	Values []float64
}

func (t Trend) Delta() float64 {
	if len(t.Values) < 2 {
		return 0
	}
	return t.Values[len(t.Values)-1] - t.Values[0]
}
