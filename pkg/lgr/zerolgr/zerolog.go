package zerolgr

import (
	"github.com/rs/zerolog"
	"music-files/internal/config"
	"music-files/pkg/util"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger(cfg *config.LoggerConfig) (*ZerologLogger, error) {
	parsedLevel, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, util.WrapError(err, "failed to parse log level")
	}

	consoleWriter := zerolog.NewConsoleWriter()

	logger := zerolog.New(consoleWriter).
		With().Str("service", "files").Logger().
		Level(parsedLevel)

	return &ZerologLogger{logger: logger}, nil
}

func (l *ZerologLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msgf("%v", args)
}

func (l *ZerologLogger) Info(args ...interface{}) {
	l.logger.Info().Msgf("%v", args)
}

func (l *ZerologLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msgf("%v", args)
}

func (l *ZerologLogger) Error(args ...interface{}) {
	l.logger.Error().Msgf("%v", args)
}

func (l *ZerologLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msgf("%v", args)
}
