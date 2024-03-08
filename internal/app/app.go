package app

import (
	"log"
	"music-files/internal/config"
	"music-files/pkg/lgr/zerolgr"
)

func Run() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("Failed to get config. %v", err)
	}

	logger, err := zerolgr.NewZerologLogger(&cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger. %v", err)
	}
	logger.Info("Logger initialized")
}
