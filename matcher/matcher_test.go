package matcher_test

import (
	"episodes/file"
	"episodes/matcher"
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

func makeEntry(name string) *MockDirEntry {
	e := &MockDirEntry{}
	e.On("Name").Return(name)
	return e
}

type MatcherSuite struct {
	suite.Suite
}

func TestMatcherSuite(t *testing.T) {
	suite.Run(t, new(MatcherSuite))
}

func (s *MatcherSuite) TestMatch_NoPatternMatches() {
	m := matcher.New()
	m.AddPattern(`^(?P<id>\d+)$`)
	entry := makeEntry("no-match-here")
	result := m.Match(file.File{Entry: entry})
	s.Nil(result)
	entry.AssertExpectations(s.T())
}

func (s *MatcherSuite) TestMatch_StopsAtFirstMatchingPattern() {
	m := matcher.New(
		matcher.WithPattern(`^(?P<first>first)$`),
		matcher.WithPattern(`^(?P<second>first)$`),
	)
	entry := makeEntry("first")
	result := m.Match(file.File{Entry: entry})
	s.NotNil(result)
	s.Equal("first", (*result)["first"])
	s.Equal("", (*result)["second"])
	entry.AssertExpectations(s.T())
}

func (s *MatcherSuite) TestAddPattern() {
	m := matcher.New()
	m.AddPattern(`^(?P<id>\d+)\.txt$`)

	s.Run("matches file", func() {
		entry := makeEntry("42.txt")
		result := m.Match(file.File{Entry: entry})
		s.NotNil(result)
		s.Equal("42", (*result)["id"])
		entry.AssertExpectations(s.T())
	})

	s.Run("no match for unrelated filename", func() {
		entry := makeEntry(gofakeit.Animal() + ".mp4")
		result := m.Match(file.File{Entry: entry})
		s.Nil(result)
		entry.AssertExpectations(s.T())
	})
}
