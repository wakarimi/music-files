package logging

import "github.com/rs/zerolog"

func ParseZerologLevel(levelStr string) zerolog.Level {
	switch levelStr {
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
