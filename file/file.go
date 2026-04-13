package file

import (
	"os"
	"regexp"
)

func New(path string, entry os.DirEntry) *File {
	return &File{Path: path, Entry: entry}
}

type File struct {
	Path  string
	Entry os.DirEntry
}

func (f *File) Match(pattern *regexp.Regexp) []string {
	return pattern.FindStringSubmatch(f.Entry.Name())
}
