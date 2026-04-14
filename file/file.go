package file

import (
	"os"
)

func New(path string, entry os.DirEntry) *File {
	return &File{Path: path, Entry: entry}
}

type File struct {
	Path  string
	Entry os.DirEntry
}
