package site

import (
	"slices"
)

type Site int

const (
	SiteFapello Site = iota + 1
)

var sites = []string{"fapello"}

func (s Site) String() string {
	return sites[s-1]
}

func FromString(s string) (Site, bool) {
	idx := slices.Index(sites, s)
	if idx == -1 {
		return 0, false
	}
	return Site(idx + 1), true
}

func Sites() []string {
	return sites
}
