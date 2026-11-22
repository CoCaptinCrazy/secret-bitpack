package routing

import "net/netip"

type Route struct {
	Prefix  netip.Prefix
	Gateway netip.Addr
}
