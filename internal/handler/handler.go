package handler

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type useCase interface {
	AddRoot(input AddRootInput) (output AddRootOutput, err error)
	DeleteRoot(input DeleteRootInput) (err error)
}

type AddRootInput struct {
	Path string
}

type AddRootOutput struct {
	DirID int
	Path  string
}

type DeleteRootInput struct {
	DirID int
}

type Handler struct {
	useCase      useCase
	bundle       i18n.Bundle
	engLocalizer i18n.Localizer
}

func New(useCase useCase,
	bundle i18n.Bundle) *Handler {
	return &Handler{
		useCase:      useCase,
		bundle:       bundle,
		engLocalizer: *i18n.NewLocalizer(&bundle, "en_US"),
	}
}
