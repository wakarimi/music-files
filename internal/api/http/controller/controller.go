package controller

import (
	"github.com/rs/zerolog"
	"music-files/internal/domain/usecase"
	"music-files/pkg/loclzr"
)

type Controller struct {
	useCases  usecase.UseCases
	localizer *loclzr.Localizer
	log       *zerolog.Logger
}

func New(useCases usecase.UseCases, localizer *loclzr.Localizer, log *zerolog.Logger) *Controller {
	return &Controller{
		useCases:  useCases,
		localizer: localizer,
		log:       log,
	}
}
