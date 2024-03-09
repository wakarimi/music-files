package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"music-files/internal/config"
	"music-files/pkg/core"
	"music-files/pkg/logging"
	"os"
)

func Run() {
	cfg, err := config.Parse()
	if err != nil {
		panic(fmt.Sprintf("failed to get config: %v", err))
	}

	logger, err := createLogger(*cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	logger.Debug().Msg("Logger initialized")
	logger.Info().Msg("Logger initialized")
	logger.Warn().Msg("Logger initialized")
	logger.Error().Msg("Logger initialized")
}

func createLogger(cfg config.Config) (*zerolog.Logger, error) {
	level := logging.ParseZerologLevel(cfg.Logger.Level)

	switch cfg.App.Environment {
	case core.EnvProd:
		logger := zerolog.New(os.Stdout).Level(level).With().
			Str("service", cfg.App.Name).Timestamp().Caller().Logger()
		return &logger, nil
	case core.EnvTest:
		logger := zerolog.New(os.Stdout).Level(level).With().
			Str("service", cfg.App.Name).Timestamp().Caller().Logger()
		return &logger, nil
	case core.EnvDev:
		consoleWriter := zerolog.NewConsoleWriter()
		logger := zerolog.New(consoleWriter).Level(level).With().
			Str("service", cfg.App.Name).Timestamp().Caller().Logger()
		return &logger, nil
	default:
		return nil, fmt.Errorf("failed to initialize logger")
	}
}
