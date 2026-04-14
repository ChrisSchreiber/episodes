package show_test

import (
	"episodes/show"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestShowAddEpisode(t *testing.T) {
	name := gofakeit.Animal()
	s := show.Show{Name: name}

	t.Run("creates new season when adding to empty show", func(t *testing.T) {
		s.AddEpisode(1, 1)
		assert.Len(t, s.Seasons, 1)
		assert.Equal(t, uint8(1), s.Seasons[0].Number)
		assert.Len(t, s.Seasons[0].Episodes, 1)
		assert.Equal(t, uint8(1), s.Seasons[0].Episodes[0].EpisodeNumber)
	})

	t.Run("adds episode to existing season", func(t *testing.T) {
		s.AddEpisode(1, 2)
		assert.Len(t, s.Seasons, 1)
		assert.Len(t, s.Seasons[0].Episodes, 2)
	})

	t.Run("does not add duplicate episode", func(t *testing.T) {
		s.AddEpisode(1, 1)
		assert.Len(t, s.Seasons[0].Episodes, 2)
	})

	t.Run("creates new season for different season number", func(t *testing.T) {
		s.AddEpisode(2, 1)
		assert.Len(t, s.Seasons, 2)
		assert.Equal(t, uint8(2), s.Seasons[1].Number)
		assert.Len(t, s.Seasons[1].Episodes, 1)
	})
}

func TestShowGetSeason(t *testing.T) {
	s := show.Show{
		Name: gofakeit.Animal(),
		Seasons: []show.Season{
			{Number: 1, Episodes: []show.Episode{{SeasonNumber: 1, EpisodeNumber: 3}}},
			{Number: 2, Episodes: []show.Episode{{SeasonNumber: 2, EpisodeNumber: 5}}},
		},
	}

	t.Run("returns existing season", func(t *testing.T) {
		season := s.GetSeason(1)
		assert.Equal(t, uint8(1), season.Number)
		assert.Len(t, season.Episodes, 1)
		assert.Equal(t, uint8(3), season.Episodes[0].EpisodeNumber)
	})

	t.Run("returns empty season for unknown number", func(t *testing.T) {
		season := s.GetSeason(9)
		assert.Equal(t, uint8(9), season.Number)
		assert.Empty(t, season.Episodes)
	})
}

func TestShowNextEpisode(t *testing.T) {
	t.Run("returns season 1 episode 1 when no seasons", func(t *testing.T) {
		s := show.Show{Name: gofakeit.Animal()}
		next := s.NextEpisode()
		assert.Equal(t, show.Episode{SeasonNumber: 1, EpisodeNumber: 1}, next)
	})

	t.Run("returns next episode after last of only season", func(t *testing.T) {
		s := show.Show{
			Name: gofakeit.Animal(),
			Seasons: []show.Season{
				{Number: 3, Episodes: []show.Episode{
					{SeasonNumber: 3, EpisodeNumber: 2},
					{SeasonNumber: 3, EpisodeNumber: 5},
				}},
			},
		}
		next := s.NextEpisode()
		assert.Equal(t, show.Episode{SeasonNumber: 3, EpisodeNumber: 6}, next)
	})

	t.Run("returns next episode after last of highest season", func(t *testing.T) {
		s := show.Show{
			Name: gofakeit.Animal(),
			Seasons: []show.Season{
				{Number: 2, Episodes: []show.Episode{
					{SeasonNumber: 2, EpisodeNumber: 10},
				}},
				{Number: 1, Episodes: []show.Episode{
					{SeasonNumber: 1, EpisodeNumber: 6},
				}},
			},
		}
		next := s.NextEpisode()
		assert.Equal(t, show.Episode{SeasonNumber: 2, EpisodeNumber: 11}, next)
	})
}
