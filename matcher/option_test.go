package matcher_test

import (
	"episodes/file"
	"episodes/matcher"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OptionSuite struct {
	suite.Suite
}

func TestOptionSuite(t *testing.T) {
	suite.Run(t, new(OptionSuite))
}

func (s *OptionSuite) TestWithPattern() {
	m := matcher.New(matcher.WithPattern(`^(?P<slug>[a-z-]+)\.md$`))

	s.Run("matches filename", func() {
		entry := makeEntry("hello-world.md")
		result := m.Match(file.File{Entry: entry})
		s.NotNil(result)
		s.Equal("hello-world", (*result)["slug"])
		entry.AssertExpectations(s.T())
	})

	s.Run("no match for unrelated filename", func() {
		entry := makeEntry("hello-world.txt")
		result := m.Match(file.File{Entry: entry})
		s.Nil(result)
		entry.AssertExpectations(s.T())
	})
}

func (s *OptionSuite) TestWithShowPatterns() {
	m := matcher.New(matcher.WithShowPatterns())

	cases := []struct {
		filename string
		name     string
		season   string
		episode  string
		ext      string
	}{
		{"Breaking-Bad-S01E01.mkv", "Breaking-Bad", "01", "01", ".mkv"},
		{"My-Show-S12E34.mp4", "My-Show", "12", "34", ".mp4"},
		{"A-S02E10", "A", "02", "10", ""},
	}

	for _, tc := range cases {
		s.Run(tc.filename, func() {
			entry := makeEntry(tc.filename)
			result := m.Match(file.File{Entry: entry})
			s.NotNil(result)
			s.Equal(tc.name, (*result)["name"])
			s.Equal(tc.season, (*result)["season"])
			s.Equal(tc.episode, (*result)["episode"])
			s.Equal(tc.ext, (*result)["extenstion"])
			entry.AssertExpectations(s.T())
		})
	}

	s.Run("does not match track filename", func() {
		entry := makeEntry("artist_t03.mp3")
		result := m.Match(file.File{Entry: entry})
		s.Nil(result)
		entry.AssertExpectations(s.T())
	})
}

func (s *OptionSuite) TestWithTrackPatterns() {
	m := matcher.New(matcher.WithTrackPatterns())

	cases := []struct {
		filename string
		track    string
		ext      string
	}{
		{"artist_t01.mp3", "01", ".mp3"},
		{"some_long_name_t12.flac", "12", ".flac"},
		{"prefix_t99", "99", ""},
	}

	for _, tc := range cases {
		s.Run(tc.filename, func() {
			entry := makeEntry(tc.filename)
			result := m.Match(file.File{Entry: entry})
			s.NotNil(result)
			s.Equal(tc.track, (*result)["track"])
			s.Equal(tc.ext, (*result)["extenstion"])
			entry.AssertExpectations(s.T())
		})
	}

	s.Run("does not match show filename", func() {
		entry := makeEntry("Show-S01E01.mkv")
		result := m.Match(file.File{Entry: entry})
		s.Nil(result)
		entry.AssertExpectations(s.T())
	})
}
