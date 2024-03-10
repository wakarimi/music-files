package controller

import (
	"github.com/rs/zerolog"
	"music-files/internal/domain/usecase"
)

type Controller struct {
	useCases usecase.UseCases
	log      *zerolog.Logger
}

func New(useCases usecase.UseCases, log *zerolog.Logger) *Controller {
	return &Controller{
		useCases: useCases,
		log:      log,
	}
}
