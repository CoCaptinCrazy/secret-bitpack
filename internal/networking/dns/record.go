package dns

import "net"

type Record struct {
	Name  string
	Type  string
	Value string
	TTL   int
}

func (r Record) IP() net.IP { return net.ParseIP(r.Value) }
