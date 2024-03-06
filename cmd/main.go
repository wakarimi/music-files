package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"music-files/internal/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get config")
	}
	fmt.Println(cfg)
}
