package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"music-files/internal/app"
	"music-files/internal/config"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	initializeLogger(cfg.App.LoggingLevel)

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create application")
	}

	err = application.StartHTTPServer(cfg.HTTP)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}

func initializeLogger(level zerolog.Level) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().Caller().Logger().
		With().Str("service", "authentication").Logger().
		Level(level)
	log.Debug().Msg("Logger initialized")
}
