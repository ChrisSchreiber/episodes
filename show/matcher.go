package show

import (
	"episodes/directory"
	"episodes/file"
	"errors"
	"regexp"
	"strconv"
)

func NewMatcher(pattern string, directory *directory.Directory) *Matcher {
	matcher := &Matcher{pattern: regexp.MustCompile(pattern)}
	for _, file := range directory.Files {
		matcher.Check(file)
	}
	return matcher
}

type Matcher struct {
	pattern *regexp.Regexp
	Shows   []Show
}

func (s *Matcher) Add(name string, seasonNumber uint8, episodeNumber uint8) {
	for i, show := range s.Shows {
		if show.Name == name {
			s.Shows[i].AddEpisode(seasonNumber, episodeNumber)
			return
		}
	}
	show := Show{Name: name}
	show.AddEpisode(seasonNumber, episodeNumber)
	s.Shows = append(s.Shows, show)
}

func (s *Matcher) Get(name string) ( *Show) {
	for _, show := range s.Shows {
		if show.Name == name {
			return &show
		}
	}
	return nil
}

func (s *Matcher) GetDefault() (*Show, error) {
	if len(s.Shows) == 0 {
		return nil, errors.New("no shows found")
	}
	if len(s.Shows) > 1 {
		return nil, errors.New("multiple shows found")
	}
	return &s.Shows[0], nil
}

func (s *Matcher) Check(file file.File) {
	if parts := file.Match(s.pattern); parts != nil {
		seasonNumber, _ := strconv.Atoi(parts[2])
		episodeNumber, _ := strconv.Atoi(parts[3])
		s.Add(parts[1], uint8(seasonNumber), uint8(episodeNumber))
	}
}