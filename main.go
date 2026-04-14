package main

import (
	"episodes/config"
	"episodes/directory"
	"episodes/matcher"
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
	log.Debug().EmbedObject(config).Msg("Config")

	destination, err := directory.New(config.Destination)
	if err != nil {
		log.Fatal().Err(err).Msg("Read destination directory failed")
	}

	var source *directory.Directory
	if config.Source == "" {
		source = destination
	} else {
		source, err = directory.New(config.Source)
		if err != nil {
			log.Fatal().Err(err).Msg("Read source directory failed")
		}
	}

	showList, err := show.NewList(*destination, *matcher.New(matcher.WithShowPatterns()))
	if err != nil {
		log.Fatal().Err(err).Msg("Read show files from destination failed")
	}

	var destinationShow *show.Show
	if config.Show == "" {
		destinationShow, err = showList.GetDefault()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not determin show name from existing files")
		}
	} else {
		destinationShow = showList.Get(config.Show)
		if destinationShow == nil {
			destinationShow = &show.Show{Name: config.Show}
		}
	}

	var next show.Episode
	if config.Season == 0 {
		// get last known episode of last known season; defaults to season 1 episode 1
		next = destinationShow.NextEpisode()
		if config.Episode != 0 {
			// override season last episode
			next.EpisodeNumber = uint8(config.Episode) - 1
		}
	} else {
		// get known episodes for season
		season := destinationShow.GetSeason(uint8(config.Season))
		if config.Episode == 0 {
			// get last known episode for season. defaults to episode 1
			next = season.NextEpisode()
		} else {
			// override season last episode
			next = show.Episode{SeasonNumber: next.SeasonNumber, EpisodeNumber: uint8(config.Episode) - 1}
		}
	}

	log.Info().Str("Show", destinationShow.Name).Uint8("Season", next.SeasonNumber).Uint8("Episode", next.EpisodeNumber).Msg("Next episode")

	tracks := new(track.List)
	err = tracks.FromDirectory(source, matcher.New(matcher.WithTrackPatterns()))
	if err != nil {
		log.Fatal().Err(err).Msg("Read tracks from source directory failed")
	}

	if len(tracks.Tracks) == 0 {
		log.Fatal().Msg("Tracks not found")
	}

	if config.Reverse {
		tracks.SortDescending()
	} else {
		tracks.SortAscending()
	}
	for _, track := range tracks.Tracks {
		err := track.Move(config.Destination, destinationShow.Name, next, config.DryRun)
		if err != nil {
			log.Fatal().Err(err).Msg("Move file failed")
		}
		next.EpisodeNumber++
	}
}
