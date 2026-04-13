package show_test

import (
	"episodes/directory"
	"episodes/file"
	"episodes/show"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockDirEntry struct {
	mock.Mock
}

func (m *MockDirEntry) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockDirEntry) Info() (os.FileInfo, error) {
	args := m.Called()
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockDirEntry) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDirEntry) Type() fs.FileMode {
	args := m.Called()
	return args.Get(0).(fs.FileMode)
}

const testPattern = `^([\w\s]+)-(\d+)\+(\d+)\.(\w+)$`

func (s *MatcherSuite) makeEntry(name string) *MockDirEntry {
	e := &MockDirEntry{}
	e.On("Name").Return(name)
	return e
}

type MatcherSuite struct {
	suite.Suite
	emptyDir *directory.Directory
}

func (s *MatcherSuite) SetupTest() {
	s.emptyDir = &directory.Directory{}
}

func TestMatcherSuite(t *testing.T) {
	suite.Run(t, new(MatcherSuite))
}

func (s *MatcherSuite) TestNewMatcher() {
	name := gofakeit.Animal()
	season := gofakeit.Uint8()
	episode := gofakeit.Uint8()
	extension := gofakeit.FileExtension()

	matching := s.makeEntry(fmt.Sprintf("%s-%d+%d.%s", name, season, episode, extension))
	nonMatching := s.makeEntry(fmt.Sprintf("%s%d%d%s", name, season, episode, extension))

	dir := &directory.Directory{
		Files: []file.File{
			{Entry: matching},
			{Entry: nonMatching},
		},
	}

	matcher := show.NewMatcher(testPattern, dir)

	s.Len(matcher.Shows, 1, "non-matching file should be ignored")
	s.Equal(name, matcher.Shows[0].Name)
	s.Len(matcher.Shows[0].Seasons, 1)
	s.Equal(season, matcher.Shows[0].Seasons[0].Number)
	s.Len(matcher.Shows[0].Seasons[0].Episodes, 1)
	s.Equal(season, matcher.Shows[0].Seasons[0].Episodes[0].SeasonNumber)
	s.Equal(episode, matcher.Shows[0].Seasons[0].Episodes[0].EpisodeNumber)

	matching.AssertExpectations(s.T())
	nonMatching.AssertExpectations(s.T())
}

func (s *MatcherSuite) TestCheck() {
	s.Run("matching file is added", func() {
		entry := s.makeEntry("MyShow-2+5.mkv")
		matcher := show.NewMatcher(testPattern, s.emptyDir)
		matcher.Check(file.File{Entry: entry})
		s.Len(matcher.Shows, 1)
		s.Equal("MyShow", matcher.Shows[0].Name)
		season := matcher.Shows[0].GetSeason(2)
		s.True(season.HasEpisode(5))
		entry.AssertExpectations(s.T())
	})

	s.Run("non-matching file is ignored", func() {
		entry := s.makeEntry("no-match-here")
		matcher := show.NewMatcher(testPattern, s.emptyDir)
		matcher.Check(file.File{Entry: entry})
		s.Empty(matcher.Shows)
		entry.AssertExpectations(s.T())
	})
}

func (s *MatcherSuite) TestGetDefault_NoShows() {
	matcher := show.Matcher{}
	_, err := matcher.GetDefault()
	s.Error(err)
}

func (s *MatcherSuite) TestOneShow() {
	seasonNumber := gofakeit.Uint8()
	expectedShow := show.Show{
		Name: gofakeit.AnimalType(),
		Seasons: []show.Season{
			{Number: seasonNumber, Episodes: []show.Episode{
				{SeasonNumber: seasonNumber, EpisodeNumber: 1},
				{SeasonNumber: seasonNumber, EpisodeNumber: 2},
			}},
		},
	}

	matcher := show.Matcher{}
	matcher.Add(expectedShow.Name, seasonNumber, 1)
	matcher.Add(expectedShow.Name, seasonNumber, 2)

	s.Run("Get existing show", func() {
		s.Equal(expectedShow, *matcher.Get(expectedShow.Name))
	})

	s.Run("Get unknown show returns nil", func() {
		s.Nil(matcher.Get(gofakeit.City()))
	})

	s.Run("GetDefault returns the only show", func() {
		result, err := matcher.GetDefault()
		s.NoError(err)
		s.Equal(expectedShow, *result)
	})
}

func (s *MatcherSuite) TestMultipleShows() {
	show1 := show.Show{
		Name: gofakeit.AnimalType(),
		Seasons: []show.Season{
			{Number: 1, Episodes: []show.Episode{
				{SeasonNumber: 1, EpisodeNumber: 1},
				{SeasonNumber: 1, EpisodeNumber: 2},
			}},
		},
	}
	show2 := show.Show{
		Name: gofakeit.CarType(),
		Seasons: []show.Season{
			{Number: 1, Episodes: []show.Episode{
				{SeasonNumber: 1, EpisodeNumber: 1},
			}},
			{Number: 2, Episodes: []show.Episode{
				{SeasonNumber: 2, EpisodeNumber: 42},
			}},
		},
	}

	matcher := show.Matcher{}
	matcher.Add(show1.Name, 1, 1)
	matcher.Add(show1.Name, 1, 2)
	matcher.Add(show2.Name, 1, 1)
	matcher.Add(show2.Name, 2, 42)

	s.Run("both shows are present", func() {
		s.Len(matcher.Shows, 2)
		s.Contains(matcher.Shows, show1)
		s.Contains(matcher.Shows, show2)
	})

	s.Run("Get returns correct show", func() {
		s.Equal(show1, *matcher.Get(show1.Name))
		s.Equal(show2, *matcher.Get(show2.Name))
	})

	s.Run("Get unknown show returns nil", func() {
		s.Nil(matcher.Get(gofakeit.City()))
	})

	s.Run("GetDefault errors with multiple shows", func() {
		_, err := matcher.GetDefault()
		s.Error(err)
	})
}
