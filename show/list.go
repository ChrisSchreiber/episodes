package show

import (
	"episodes/directory"
	"episodes/matcher"
	"errors"
	"strconv"
)

func NewList(directory directory.Directory, matcher matcher.Matcher) (*List, error) {
	list := &List{}
	for _, file := range directory.Files {
		parts := matcher.Match(file)
		if parts != nil {
			var name string
			var season, episode uint64
			var err error
			var exists bool
			if name, exists = (*parts)["name"]; !exists {
				return nil, errors.New("Show name not found")
			}
			season, err = strconv.ParseUint((*parts)["season"], 10, 8)
			if err != nil {
				return nil, err
			}
			episode, err = strconv.ParseUint((*parts)["episode"], 10, 8)
			if err != nil {
				return nil, err
			}
			list.Add(name, uint8(season), uint8(episode))
		}
	}
	return list, nil
}

type List struct {
	Shows []Show
}

func (s *List) Add(name string, seasonNumber uint8, episodeNumber uint8) {
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

func (s *List) Get(name string) *Show {
	for _, show := range s.Shows {
		if show.Name == name {
			return &show
		}
	}
	return nil
}

func (s *List) GetDefault() (*Show, error) {
	if len(s.Shows) == 0 {
		return nil, errors.New("no shows found")
	}
	if len(s.Shows) > 1 {
		return nil, errors.New("multiple shows found")
	}
	return &s.Shows[0], nil
}
