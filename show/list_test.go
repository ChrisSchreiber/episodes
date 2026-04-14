package show_test

import (
	"episodes/directory"
	"episodes/show"
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

func (s *ListTestSuite) makeEntry(name string) *MockDirEntry {
	e := &MockDirEntry{}
	e.On("Name").Return(name)
	return e
}

type ListTestSuite struct {
	suite.Suite
	emptyDir *directory.Directory
}

func (s *ListTestSuite) SetupTest() {
	s.emptyDir = &directory.Directory{}
}

func TestMatcherSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}

func (s *ListTestSuite) TestGetDefault_NoShows() {
	matcher := show.List{}
	_, err := matcher.GetDefault()
	s.Error(err)
}

func (s *ListTestSuite) TestOneShow() {
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

	matcher := show.List{}
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

func (s *ListTestSuite) TestMultipleShows() {
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

	matcher := show.List{}
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
