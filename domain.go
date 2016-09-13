package lshensemble

import (
	"sort"
)

type Domain struct {
	Key        string
	Size int
	Signature  Signature
}

type BySize []*Domain

func (rs BySize) Len() int           { return len(rs) }
func (rs BySize) Less(i, j int) bool { return rs[i].Size < rs[j].Size }
func (rs BySize) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }

func (rs BySize) Subset(lower, upper int) []*Domain {
	if !sort.IsSorted(rs) {
		panic("Must be sorted by domain size first")
	}
	start, end := -1, -1
	for i := range rs {
		if start == -1 && rs[i].Size >= lower {
			start = i
		}
		if end == -1 && (rs[i].Size > upper || i == len(rs)-1) {
			end = i
			break
		}
	}
	if start == -1 || end == -1 {
		panic("Cannot find such domain size range")
	}
	if end == len(rs)-1 {
		end++
	}
	return []*Domain(rs[start:end])
}