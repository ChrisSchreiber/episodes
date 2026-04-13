package directory

import (
	"episodes/file"
	"os"
)

func New(path string) (dir *Directory, err error) {
	entry, err := os.ReadDir(path)
	if err != nil {
		return
	}
	dir = &Directory{Path: path, Files: []file.File{}}
	for _, entry := range entry {
		if !entry.IsDir() {
			dir.Files = append(dir.Files, *file.New(path, entry))
		}
	}
	return
}

type Directory struct {
	Path  string
	Files []file.File
}
