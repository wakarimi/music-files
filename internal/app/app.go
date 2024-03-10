package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"music-files/internal/config"
	"music-files/internal/repository/postgres"
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

	db, err := postgres.New(cfg.DB)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect database")
		return
	}
	logger.Debug().Msg("Database connected")
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}()
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
