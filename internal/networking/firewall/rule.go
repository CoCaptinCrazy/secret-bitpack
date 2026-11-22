package firewall

type Rule struct {
	Protocol string
	Port     int
	Action   string
}

func (r Rule) Matches(protocol string, port int) bool {
	return r.Protocol == protocol && (r.Port == 0 || r.Port == port)
}
