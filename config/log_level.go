package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

type LogLevel string

func (ll *LogLevel) Set(value string) error {
	value = strings.ToLower(value)
	switch value {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "":
		fallthrough
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		return fmt.Errorf("Unhandled error level: %s", value)
	}
	*ll = LogLevel(value)
	return nil
}

func (ll *LogLevel) String() string {
	return string(*ll)
}
