package show

type Season struct {
	Number   uint8
	Episodes []Episode
}

func (s *Season) HasEpisode(episodeNumber uint8) bool {
	for _, episode := range s.Episodes {
		if episode.EpisodeNumber == episodeNumber {
			return true
		}
	}
	return false
}

func (s *Season) AddEpisode(episodeNumber uint8) {
	if !s.HasEpisode(episodeNumber) {
		s.Episodes = append(s.Episodes, Episode{SeasonNumber: s.Number, EpisodeNumber: episodeNumber})
	}
}

// LastEpisode get last known episode. If there are no known episodes return episode 1.
func (s *Season) LastEpisode() (lastEpisode Episode) {
	if len(s.Episodes) == 0 {
		return Episode{SeasonNumber: s.Number, EpisodeNumber: 1}
	}
	for _, episode := range s.Episodes {
		if episode.EpisodeNumber > lastEpisode.EpisodeNumber {
			lastEpisode = episode
		}
	}
	return
}
