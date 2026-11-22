package policies

type Policy struct {
	Name   string
	Action string
	Effect string
}

func (p Policy) Allows(action string) bool {
	return p.Effect == "allow" && (p.Action == "*" || p.Action == action)
}
