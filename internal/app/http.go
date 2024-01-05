package app

import (
	"fmt"
	"music-files/internal/config"
)

func (a *App) RegisterRoutes() {
	_ = a.router.Group("/api")
}

func (a *App) StartHTTPServer(cfg config.HTTPConfig) error {
	err := a.router.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}
	return nil
}
