package track

import (
	"cmp"
	"episodes/file"
	"regexp"
	"slices"
	"strconv"
)

func NewMatcher(pattern string) *Matcher {
	return &Matcher{pattern: regexp.MustCompile(pattern)}
}

type Matcher struct {
	pattern *regexp.Regexp
	Tracks []Track
}

func (m *Matcher) Check(file file.File) {
	if parts := file.Match(m.pattern); parts != nil {
		number, _ := strconv.ParseUint(parts[2], 10, 8)
		m.Tracks = append(m.Tracks, Track{Number: uint8(number), Extension: parts[3], File: file})
	}
}

func (m *Matcher) SortAscending() {
	slices.SortFunc(m.Tracks, func(a, b Track) int {
		return cmp.Compare(a.Number, b.Number)
	})
}

func (m *Matcher) SortDescending() {
	slices.SortFunc(m.Tracks, func(a, b Track) int {
		return cmp.Compare(b.Number, a.Number)
	})
}