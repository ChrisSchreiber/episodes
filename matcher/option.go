package matcher

type Option func(matcher *Matcher)

func WithPattern(pattern string) Option {
	return func(matcher *Matcher) {
		matcher.AddPattern(pattern)
	}
}

func WithShowPatterns() Option {
	return func(matcher *Matcher) {
		matcher.AddPattern(`^(?P<name>.*)-S(?P<season>\d{2})E(?P<episode>\d{2})(?P<extenstion>\.[a-zA-Z0-9]+)?$`)
	}
}

func WithTrackPatterns() Option {
	return func(matcher *Matcher) {
		matcher.AddPattern(`^(.+_t)(?P<track>\d{2})(?P<extenstion>\.[a-zA-Z0-9]+)?$`)
	}
}
