package search

import "sort"

type Tags map[string][]string

func (t Tags) List(image string) []string {
	result := append([]string(nil), t[image]...)
	sort.Strings(result)
	return result
}
