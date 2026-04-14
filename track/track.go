package track

import (
	"episodes/file"
	"episodes/show"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func New(file file.File, values map[string]string) (*Track, error) {
	track := &Track{File: file}
	for k, v := range values {
		switch strings.ToLower(k) {
		case "number":
			number, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return nil, err
			}
			track.Number = uint8(number)
		case "season":
			number, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return nil, err
			}
			track.SeasonNumber = uint8(number)
		case "disk":
			number, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return nil, err
			}
			track.Disk = uint8(number)
		case "extension":
			track.Extension = v
		}
	}
	return track, nil
}

type Track struct {
	Number       uint8
	Disk         uint8
	SeasonNumber uint8
	Extension    string
	File         file.File
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
