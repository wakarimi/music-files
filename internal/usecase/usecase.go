package usecase

import "github.com/rs/zerolog"

type UseCases struct {
	log *zerolog.Logger
}

func New(log *zerolog.Logger) *UseCases {
	return &UseCases{
		log: log,
	}
}
