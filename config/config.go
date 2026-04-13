package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func New() *Config {
	config := &Config{}

	flag.Usage = usage

	flag.StringVar(&config.Destination, "destination", "", "Destination directory path. (default: currect directory)")
	flag.StringVar(&config.Source, "source", "", "Source directory path. (default: current directory)")
	flag.StringVar(&config.Show, "show", "", "Name of show. (default: Show name is determined from files in destination directory)")
	flag.UintVar(&config.Season, "season", 0, "Season to add new episodes to. (default: Season is determined from files in destination directory)")
	flag.UintVar(&config.Episode, "episode", 0, "Episode number to start at. (default: Start episode is determined from files in destination directory)")
	flag.BoolVar(&config.Reverse, "reverse", false, "Reverse track order. Highest number track is next episode.")
	flag.Var(&config.LogLevel, "log-level", "Log level. (debug|info|warn|error|fatal|disabled; default: info)")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Display source and destination files but do not rename.")

	flag.Parse()

	if config.Destination == "" {
		currentDirectory, err := os.Getwd()
		if err != nil {
			log.Panic().Err(err).Msg("Current directory unknown; destination directory required")
		}
		config.Destination = currentDirectory
	}

	return config
}

type Config struct {
	Destination string
	Source      string
	Show        string
	Season      uint
	Episode     uint
	Reverse     bool
	LogLevel    LogLevel
	DryRun      bool
}

func usage() {
	fmt.Fprintln(os.Stderr, `Rename video files wtih series name, season and episode.

episodes --show Andor
	Rename files in current directory. Source file names must end with "_t##.[EXTENSION]", for example "title_01.mkv". Source files are renamed with the format "Andor-S[SEASON]E[EPISODE].[EXTENSTION]". If no existing files with the same format exist SEASON and EPISODE start at "01" i.e. "Andor-S01E01.mkv". Episodes are ordered by the digits in the source file names.

	`)
	flag.PrintDefaults()
}
