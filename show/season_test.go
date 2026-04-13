package show_test

import (
	"episodes/show"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestSeasonHasEpisode(t *testing.T) {
	seasonNumber := gofakeit.Uint8()
	season := show.Season{
		Number: seasonNumber,
		Episodes: []show.Episode{
			{SeasonNumber: seasonNumber, EpisodeNumber: 1},
			{SeasonNumber: seasonNumber, EpisodeNumber: 3},
		},
	}

	assert.True(t, season.HasEpisode(1))
	assert.True(t, season.HasEpisode(3))
	assert.False(t, season.HasEpisode(2))
	assert.False(t, season.HasEpisode(0))
}

func TestSeasonAddEpisode(t *testing.T) {
	seasonNumber := gofakeit.Uint8()
	season := show.Season{Number: seasonNumber}

	t.Run("adds new episode", func(t *testing.T) {
		season.AddEpisode(1)
		assert.Len(t, season.Episodes, 1)
		assert.Equal(t, show.Episode{SeasonNumber: seasonNumber, EpisodeNumber: 1}, season.Episodes[0])
	})

	t.Run("does not add duplicate episode", func(t *testing.T) {
		season.AddEpisode(1)
		assert.Len(t, season.Episodes, 1)
	})

	t.Run("adds another episode", func(t *testing.T) {
		season.AddEpisode(2)
		assert.Len(t, season.Episodes, 2)
		assert.Equal(t, show.Episode{SeasonNumber: seasonNumber, EpisodeNumber: 2}, season.Episodes[1])
	})
}

func TestSeasonLastEpisode(t *testing.T) {
	seasonNumber := gofakeit.Uint8()

	t.Run("returns episode 1 when no episodes", func(t *testing.T) {
		season := show.Season{Number: seasonNumber}
		last := season.LastEpisode()
		assert.Equal(t, show.Episode{SeasonNumber: seasonNumber, EpisodeNumber: 1}, last)
	})

	t.Run("returns only episode when one episode", func(t *testing.T) {
		season := show.Season{
			Number:   seasonNumber,
			Episodes: []show.Episode{{SeasonNumber: seasonNumber, EpisodeNumber: 5}},
		}
		last := season.LastEpisode()
		assert.Equal(t, show.Episode{SeasonNumber: seasonNumber, EpisodeNumber: 5}, last)
	})

	t.Run("returns highest episode number", func(t *testing.T) {
		season := show.Season{
			Number: seasonNumber,
			Episodes: []show.Episode{
				{SeasonNumber: seasonNumber, EpisodeNumber: 3},
				{SeasonNumber: seasonNumber, EpisodeNumber: 7},
				{SeasonNumber: seasonNumber, EpisodeNumber: 2},
			},
		}
		last := season.LastEpisode()
		assert.Equal(t, show.Episode{SeasonNumber: seasonNumber, EpisodeNumber: 7}, last)
	})
}
