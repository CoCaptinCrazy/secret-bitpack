package routing

import "net/netip"

type Table struct{ Routes []Route }

func (t *Table) Add(route Route) { t.Routes = append(t.Routes, route) }

func (t Table) Lookup(address netip.Addr) (Route, bool) {
	var best Route
	found := false
	for _, route := range t.Routes {
		if route.Prefix.Contains(address) && (!found || route.Prefix.Bits() > best.Prefix.Bits()) {
			best, found = route, true
		}
	}
	return best, found
}
