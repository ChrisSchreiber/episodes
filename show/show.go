package show

type Show struct {
	Name    string
	Seasons []Season
}

func (s *Show) AddEpisode(seasonNumber uint8, episodeNumber uint8) {
	for i, season := range s.Seasons {
		if season.Number == seasonNumber {
			s.Seasons[i].AddEpisode(episodeNumber)
			return
		}
	}
	season := Season{Number: seasonNumber}
	season.AddEpisode(episodeNumber)
	s.Seasons = append(s.Seasons, season)
}

func (s *Show) GetSeason(number uint8) (season Season) {
	for _, season = range s.Seasons {
		if season.Number == number {
			return
		}
	}
	return Season{Number: number}
}

// NextEpisode get the next known episode.
// If there are no known seasons, season 1 episode 1 is returned.
func (s *Show) NextEpisode() (next Episode) {
	if len(s.Seasons) == 0 {
		return Episode{SeasonNumber: 1, EpisodeNumber: 1}
	}
	var lastSeason Season
	for _, season := range s.Seasons {
		if season.Number > lastSeason.Number {
			lastSeason = season
		}
	}
	return lastSeason.NextEpisode()
}
