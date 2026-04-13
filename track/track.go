package track

import (
	"episodes/file"
	"episodes/show"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Track struct {
	Number    uint8
	Extension string
	File      file.File
}

func (t *Track) Move(directory, name string, episode show.Episode, dryRun bool) (err error) {
	source := filepath.Join(t.File.Path, t.File.Entry.Name())
	destination := filepath.Join(directory, fmt.Sprintf("%s-S%02dE%02d%s", name, episode.SeasonNumber, episode.EpisodeNumber, t.Extension))
	log.Info().Str("Source", source).Str("Destination", destination).Msg("Move")
	if dryRun {
		return
	}
	return os.Rename(source, destination)
}
