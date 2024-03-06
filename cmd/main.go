package main

import (
	"log"
	"music-files/internal/config"
	pkglogger "music-files/pkg/logger"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("Failed to get config. %v", err)
	}

	logger, err := pkglogger.NewZerologLogger(&cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger. %v", err)
	}
	logger.Info("Logger initialized")
}
