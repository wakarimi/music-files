package app

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"music-files/internal/api/http/controller"
	"music-files/internal/api/http/router"
	"music-files/internal/config"
	"music-files/internal/repository/postgres"
	"music-files/internal/usecase"
	"music-files/pkg/core"
	"music-files/pkg/loclzr"
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

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", json.Unmarshal)
	bundle.MustLoadMessageFile("internal/api/locale/en-US.json")
	bundle.MustLoadMessageFile("internal/api/locale/ru-RU.json")
	localizer := loclzr.New(bundle)

	httpController := controller.New(useCases, localizer, log)

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
