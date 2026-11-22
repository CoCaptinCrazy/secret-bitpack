package scanning

type Finding struct {
	Rule     string
	Severity string
	Message  string
}

type Report struct {
	Findings []Finding
}

func (r Report) HasCritical() bool {
	for _, finding := range r.Findings {
		if finding.Severity == "critical" {
			return true
		}
	}
	return false
}
