package main

import (
	"episodes/config"
	"episodes/directory"
	"episodes/show"
	"episodes/track"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	config := config.New()

	destination, err := directory.New(config.Destination)
	if err != nil {
		log.Panic().Err(err).Msg("read destination directory failed")
	}

	var source *directory.Directory
	if config.Source == "" {
		source = destination
	} else {
		source, err = directory.New(config.Source)
		if err != nil {
			log.Panic().Err(err).Msg("read source directory failed")
		}
	}

	showMatcher := show.NewMatcher(`^(.*)-S(\d{2})E(\d{2})(\.[a-zA-Z0-9]+)?$`, destination)

	var destinationShow *show.Show
	if config.Show == "" {
		destinationShow, err = showMatcher.GetDefault()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not determin show name from existing files")
		}
	} else {
		destinationShow = showMatcher.Get(config.Show)
		if destinationShow == nil {
			destinationShow = &show.Show{Name: config.Show}
		}
	}

	var episode show.Episode
	if config.Season == 0 {
		// get last known episode of last known season; defaults to season 1 episode 1
		episode = destinationShow.LastEpisode()
		if config.Episode != 0 {
			// override season last episode
			episode.EpisodeNumber = uint8(config.Episode) - 1
		}
	} else {
		// get known episodes for season
		season := destinationShow.GetSeason(uint8(config.Season))
		if config.Episode == 0 {
			// get last known episode for season. defaults to episode 1
			episode = season.LastEpisode()
		} else {
			// override season last episode
			episode = show.Episode{SeasonNumber: episode.SeasonNumber, EpisodeNumber: uint8(config.Episode) - 1}
		}
	}
	last := destinationShow.LastEpisode()

	log.Info().Str("Show", destinationShow.Name).Uint8("Season", last.SeasonNumber).Uint8("Episode", last.EpisodeNumber).Msg("Show identified")

	trackMatcher := track.NewMatcher(`^(.+_t)(\d{2})(\.[a-zA-Z0-9]+)?$`)
	for _, file := range source.Files {
		trackMatcher.Check(file)
	}

	if len(trackMatcher.Tracks) == 0 {
		log.Fatal().Msg("Tracks not found")
	}

	if config.Reverse {
		trackMatcher.SortDescending()
	} else {
		trackMatcher.SortAscending()
	}
	for _, track := range trackMatcher.Tracks {
		last.EpisodeNumber++
		err := track.Move(config.Destination, destinationShow.Name, last, config.DryRun)
		if err != nil {
			log.Fatal().Err(err).Msg("Move file failed")
		}
	}
}
