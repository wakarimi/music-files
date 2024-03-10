package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"music-files/internal/api/http/controller"
	"music-files/internal/api/http/router"
	"music-files/internal/config"
	"music-files/internal/repository/postgres"
	"music-files/internal/usecase"
	"music-files/pkg/core"
	"music-files/pkg/logging"
	"os"
)

func Run() {
	cfg, err := config.Parse()
	if err != nil {
		panic(fmt.Sprintf("failed to get config: %v", err))
	}

	log, err := createLogger(*cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	log.Debug().Msg("Logger initialized")

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect database")
		return
	}
	log.Debug().Msg("Database connected")
	defer func() {
		err = db.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	useCases := usecase.New(log)

	httpController := controller.New(useCases, log)

	httpRouter := router.New(httpController, log)
	httpRouter.RegisterRoutes()
	err = httpRouter.StartHTTPServer(cfg.HTTP)
	if err != nil {
		log.Error().Msg("Failed to start http server")
		return
	}
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
