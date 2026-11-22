package firewall

type Firewall struct{ Rules []Rule }

func (f Firewall) Allows(protocol string, port int) bool {
	allowed := false
	for _, rule := range f.Rules {
		if rule.Matches(protocol, port) {
			allowed = rule.Action == "allow"
		}
	}
	return allowed
}
