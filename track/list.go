package track

import (
	"cmp"
	"episodes/directory"
	"episodes/matcher"
	"slices"
)

type List struct {
	Tracks []Track
}

func (m *List) FromDirectory(dir *directory.Directory, matcher *matcher.Matcher) error {
	for _, file := range dir.Files {
		if parts := matcher.Match(file); parts != nil {
			track, err := New(file, *parts)
			if err != nil {
				return err
			}
			m.Tracks = append(m.Tracks, *track)
		}
	}
	return nil
}

func (m *List) SortAscending() {
	slices.SortFunc(m.Tracks, func(a, b Track) int {
		if disk := cmp.Compare(a.Disk, b.Disk); disk != 0 {
			return disk
		} else {
			return cmp.Compare(a.Number, b.Number)
		}
	})
}

func (m *List) SortDescending() {
	slices.SortFunc(m.Tracks, func(a, b Track) int {
		if disk := cmp.Compare(a.Disk, b.Disk); disk != 0 {
			return disk
		} else {
			return cmp.Compare(b.Number, a.Number)
		}
	})
}
