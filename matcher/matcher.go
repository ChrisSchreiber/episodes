package matcher

import (
	"episodes/file"
	"regexp"
)

func New(options ...Option) *Matcher {
	matcher := &Matcher{}

	for _, option := range options {
		option(matcher)
	}

	return matcher
}

type Matcher struct {
	patterns []*regexp.Regexp
}

func (m *Matcher) AddPattern(pattern string) {
	m.patterns = append(m.patterns, regexp.MustCompile(pattern))
}

func (m *Matcher) Match(file file.File) (parts *map[string]string) {
	for _, pattern := range m.patterns {
		values := pattern.FindStringSubmatch(file.Entry.Name())
		if values != nil {
			parts = &map[string]string{}
			for i, name := range pattern.SubexpNames() {
				(*parts)[name] = values[i]
			}
			break
		}
	}
	return
}
